package readiness

import (
	"slices"
	"strings"

	"github.com/OpenUdon/authoring/decision"
	"github.com/OpenUdon/authoring/session"
)

const (
	SeverityBlocking = "blocking"
	SeverityWarning  = "warning"
	SeverityAdvisory = "advisory"
	SeverityInfo     = "info"

	DecisionConflictIssueCode      = "decision.conflict"
	DecisionLowConfidenceIssueCode = "decision.low_confidence"
	DecisionReviewIssueCode        = "decision.review_required"

	DefaultAnswerSource          = "default"
	DefaultSuggestedAnswerSource = "suggested_answer"
)

// Issue is the shared readiness issue record used by sessions and loops.
type Issue = session.ReadinessIssue

// Result is a deterministic summary of readiness findings.
type Result struct {
	Ready    bool    `json:"ready"`
	Issues   []Issue `json:"issues,omitempty"`
	Blocking []Issue `json:"blocking,omitempty"`
	Warnings []Issue `json:"warnings,omitempty"`
	TopIssue *Issue  `json:"top_issue,omitempty"`
	Summary  Summary `json:"summary,omitempty"`
}

// Summary captures stable readiness counters for transcripts and agent
// results.
type Summary struct {
	IssueCount    int    `json:"issue_count,omitempty"`
	BlockingCount int    `json:"blocking_count,omitempty"`
	WarningCount  int    `json:"warning_count,omitempty"`
	AdvisoryCount int    `json:"advisory_count,omitempty"`
	InfoCount     int    `json:"info_count,omitempty"`
	TopCode       string `json:"top_code,omitempty"`
	TopSlot       string `json:"top_slot,omitempty"`
}

// Question is a product-neutral follow-up question plan.
type Question struct {
	ID            string   `json:"id,omitempty"`
	Prompt        string   `json:"prompt,omitempty"`
	Slots         []string `json:"slots,omitempty"`
	Required      bool     `json:"required,omitempty"`
	Forced        bool     `json:"forced,omitempty"`
	AllowDefault  bool     `json:"allow_default,omitempty"`
	DefaultAnswer string   `json:"default_answer,omitempty"`
	DefaultSource string   `json:"default_source,omitempty"`
}

// Plan is a deterministic set of planned questions.
type Plan struct {
	Questions   []Question  `json:"questions,omitempty"`
	TopQuestion *Question   `json:"top_question,omitempty"`
	Summary     PlanSummary `json:"summary,omitempty"`
}

// PlanSummary captures stable question-plan counters.
type PlanSummary struct {
	QuestionCount int `json:"question_count,omitempty"`
	ForcedCount   int `json:"forced_count,omitempty"`
	RequiredCount int `json:"required_count,omitempty"`
	DefaultCount  int `json:"default_count,omitempty"`
}

// Evaluate normalizes issues and returns a readiness result. Issues with empty
// severity are treated as blocking so incomplete findings do not auto-pass.
func Evaluate(issues []Issue) Result {
	issues = NormalizeIssues(issues)
	result := Result{Issues: issues, Ready: Ready(issues)}
	for _, issue := range issues {
		switch {
		case IsBlocking(issue):
			result.Blocking = append(result.Blocking, issue)
			result.Summary.BlockingCount++
		case IsWarning(issue):
			result.Warnings = append(result.Warnings, issue)
			result.Summary.WarningCount++
		case severity(issue) == SeverityAdvisory:
			result.Summary.AdvisoryCount++
		case severity(issue) == SeverityInfo:
			result.Summary.InfoCount++
		}
	}
	result.Summary.IssueCount = len(issues)
	if top := TopIssue(issues); top != nil {
		result.TopIssue = top
		result.Summary.TopCode = top.Code
		result.Summary.TopSlot = top.Slot
	}
	return result
}

// NormalizeResult returns a deterministic copy of result.
func NormalizeResult(result Result) Result {
	return Evaluate(result.Issues)
}

// NormalizeIssues returns deterministic readiness issues.
func NormalizeIssues(issues []Issue) []Issue {
	normalized := session.Normalize(session.State{Readiness: issues}).Readiness
	slices.SortStableFunc(normalized, CompareIssue)
	return normalized
}

// Ready reports whether no blocking issue remains.
func Ready(issues []Issue) bool {
	for _, issue := range NormalizeIssues(issues) {
		if IsBlocking(issue) {
			return false
		}
	}
	return true
}

// Blocking returns the deterministic blocking issue subset.
func Blocking(issues []Issue) []Issue {
	var out []Issue
	for _, issue := range NormalizeIssues(issues) {
		if IsBlocking(issue) {
			out = append(out, issue)
		}
	}
	return out
}

// Warnings returns the deterministic warning issue subset.
func Warnings(issues []Issue) []Issue {
	var out []Issue
	for _, issue := range NormalizeIssues(issues) {
		if IsWarning(issue) {
			out = append(out, issue)
		}
	}
	return out
}

// TopIssue returns the most important issue, with blocking findings before
// warnings and lower-severity advisory/info findings.
func TopIssue(issues []Issue) *Issue {
	issues = NormalizeIssues(issues)
	if len(issues) == 0 {
		return nil
	}
	issue := issues[0]
	return &issue
}

// IsBlocking reports whether issue should block default progression.
func IsBlocking(issue Issue) bool {
	switch severity(issue) {
	case SeverityWarning, SeverityAdvisory, SeverityInfo:
		return false
	default:
		return true
	}
}

// IsWarning reports whether issue is a nonblocking warning.
func IsWarning(issue Issue) bool {
	return severity(issue) == SeverityWarning
}

// CompareIssue orders readiness issues deterministically.
func CompareIssue(a, b Issue) int {
	if diff := severityRank(severity(a)) - severityRank(severity(b)); diff != 0 {
		return diff
	}
	return compareStrings(a.Code, b.Code, a.Slot, b.Slot, a.Message, b.Message, a.SuggestedAnswer, b.SuggestedAnswer)
}

// DecisionIssues projects decision evidence that needs confirmation into
// generic readiness issues.
func DecisionIssues(records []decision.Record) []Issue {
	decisions := decision.Merge(records)
	issues := make([]Issue, 0, len(decisions))
	for _, record := range decisions {
		if !decision.RequiresConfirmation(record) {
			continue
		}
		issues = append(issues, decisionIssue(record))
	}
	return NormalizeIssues(issues)
}

// NormalizeQuestion returns a deterministic question plan.
func NormalizeQuestion(question Question) Question {
	question.ID = strings.TrimSpace(question.ID)
	question.Prompt = strings.TrimSpace(question.Prompt)
	question.DefaultAnswer = strings.TrimSpace(question.DefaultAnswer)
	question.DefaultSource = normalizeToken(question.DefaultSource)
	question.Slots = normalizeSlots(question.Slots)
	if question.Forced {
		question.Required = true
	}
	return question
}

// NormalizeQuestions returns deterministically ordered question plans.
func NormalizeQuestions(questions []Question) []Question {
	out := make([]Question, 0, len(questions))
	for _, question := range questions {
		question = NormalizeQuestion(question)
		if emptyQuestion(question) {
			continue
		}
		out = append(out, question)
	}
	slices.SortStableFunc(out, CompareQuestion)
	return out
}

// EvaluatePlan normalizes questions and returns a stable plan summary.
func EvaluatePlan(questions []Question) Plan {
	questions = NormalizeQuestions(questions)
	plan := Plan{Questions: questions}
	for _, question := range questions {
		if question.Forced {
			plan.Summary.ForcedCount++
		}
		if question.Required {
			plan.Summary.RequiredCount++
		}
		if _, _, ok := DefaultAnswer(question); ok {
			plan.Summary.DefaultCount++
		}
	}
	plan.Summary.QuestionCount = len(questions)
	if len(questions) > 0 {
		question := questions[0]
		plan.TopQuestion = &question
	}
	return plan
}

// SuggestedQuestion builds a question that may safely use an issue's suggested
// answer when the issue is not forced.
func SuggestedQuestion(issue Issue, prompt string) Question {
	issues := NormalizeIssues([]Issue{issue})
	if len(issues) > 0 {
		issue = issues[0]
	}
	return NormalizeQuestion(Question{
		ID:            issueID(issue),
		Prompt:        prompt,
		Slots:         issueSlots(issue),
		Required:      IsBlocking(issue),
		AllowDefault:  issue.SuggestedAnswer != "",
		DefaultAnswer: issue.SuggestedAnswer,
		DefaultSource: DefaultSuggestedAnswerSource,
	})
}

// ForcedQuestion builds a required operator question. Suggested answers may be
// shown as defaults, but they are not auto-applied.
func ForcedQuestion(issue Issue, prompt string) Question {
	question := SuggestedQuestion(issue, prompt)
	question.Forced = true
	question.Required = true
	question.AllowDefault = false
	return NormalizeQuestion(question)
}

// DefaultAnswer returns an auto-usable default answer for question.
func DefaultAnswer(question Question) (string, string, bool) {
	question = NormalizeQuestion(question)
	if question.Forced || !question.AllowDefault || question.DefaultAnswer == "" {
		return "", "", false
	}
	return question.DefaultAnswer, firstNonEmpty(question.DefaultSource, DefaultAnswerSource), true
}

// CompareQuestion orders question plans deterministically.
func CompareQuestion(a, b Question) int {
	a = NormalizeQuestion(a)
	b = NormalizeQuestion(b)
	if diff := boolRank(b.Forced) - boolRank(a.Forced); diff != 0 {
		return diff
	}
	if diff := boolRank(b.Required) - boolRank(a.Required); diff != 0 {
		return diff
	}
	if diff := boolRank(b.AllowDefault) - boolRank(a.AllowDefault); diff != 0 {
		return diff
	}
	return compareStrings(a.ID, b.ID, strings.Join(a.Slots, "\x00"), strings.Join(b.Slots, "\x00"), a.Prompt, b.Prompt)
}

func decisionIssue(record decision.Record) Issue {
	record = decision.Normalize(record)
	code := DecisionReviewIssueCode
	message := "decision requires confirmation"
	switch decision.Behavior(record) {
	case decision.BehaviorConflict:
		code = DecisionConflictIssueCode
		message = "decision has conflicting evidence"
	case decision.BehaviorLowConfidence:
		code = DecisionLowConfidenceIssueCode
		message = "decision has low confidence"
	}
	return Issue{
		Code:            code,
		Severity:        SeverityBlocking,
		Slot:            record.Slot,
		Message:         message,
		SuggestedAnswer: record.Value,
	}
}

func severity(issue Issue) string {
	return normalizeToken(issue.Severity)
}

func severityRank(severity string) int {
	switch normalizeToken(severity) {
	case "", "error", SeverityBlocking:
		return 0
	case SeverityWarning:
		return 1
	case SeverityAdvisory:
		return 2
	case SeverityInfo:
		return 3
	default:
		return 0
	}
}

func normalizeSlots(slots []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(slots))
	for _, slot := range slots {
		slot = strings.TrimSpace(slot)
		if slot == "" || seen[slot] {
			continue
		}
		seen[slot] = true
		out = append(out, slot)
	}
	return out
}

func issueID(issue Issue) string {
	return firstNonEmpty(issue.Code, issue.Slot)
}

func issueSlots(issue Issue) []string {
	if strings.TrimSpace(issue.Slot) == "" {
		return nil
	}
	return []string{issue.Slot}
}

func emptyQuestion(question Question) bool {
	return question.ID == "" && question.Prompt == "" && len(question.Slots) == 0
}

func boolRank(value bool) int {
	if value {
		return 1
	}
	return 0
}

func compareStrings(values ...string) int {
	for i := 0; i+1 < len(values); i += 2 {
		if values[i] < values[i+1] {
			return -1
		}
		if values[i] > values[i+1] {
			return 1
		}
	}
	return 0
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func normalizeToken(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(value)), "_"))
}
