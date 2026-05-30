package report

import (
	"encoding/json"
	"slices"

	"github.com/OpenUdon/authoring/internal/norm"
	"github.com/OpenUdon/authoring/trust"
)

const (
	ScorecardVersion = "authoring.scorecard.v1"

	OutcomeSkipped = "skipped"

	FailureFamilyUnclassified = "unclassified"
)

// Scorecard is a product-neutral summary of fixture/variant results.
type Scorecard struct {
	Version  string            `json:"version"`
	Name     string            `json:"name,omitempty"`
	Variants []VariantResult   `json:"variants,omitempty"`
	Summary  ScorecardSummary  `json:"summary"`
	Report   *ReportMetadata   `json:"report,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// FixtureResult groups variant outcomes for one fixture.
type FixtureResult struct {
	FixtureID       string            `json:"fixture_id"`
	Group           string            `json:"group,omitempty"`
	ExpectedOutcome string            `json:"expected_outcome,omitempty"`
	Variants        []VariantResult   `json:"variants,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// VariantResult records one fixture/variant run outcome.
type VariantResult struct {
	FixtureID       string                   `json:"fixture_id"`
	VariantID       string                   `json:"variant_id,omitempty"`
	Group           string                   `json:"group,omitempty"`
	ExpectedOutcome string                   `json:"expected_outcome,omitempty"`
	ObservedOutcome string                   `json:"observed_outcome"`
	OutcomeMatched  bool                     `json:"outcome_matched"`
	FailureFamily   string                   `json:"failure_family,omitempty"`
	Message         string                   `json:"message,omitempty"`
	Diagnostics     []trust.DiagnosticRecord `json:"diagnostics,omitempty"`
	Metadata        map[string]string        `json:"metadata,omitempty"`
}

// ScorecardSummary contains stable aggregate counters.
type ScorecardSummary struct {
	Total              int                    `json:"total,omitempty"`
	Complete           int                    `json:"complete,omitempty"`
	NeedsInput         int                    `json:"needs_input,omitempty"`
	Failed             int                    `json:"failed,omitempty"`
	Canceled           int                    `json:"canceled,omitempty"`
	Skipped            int                    `json:"skipped,omitempty"`
	ExpectedMatched    int                    `json:"expected_matched,omitempty"`
	ExpectedMismatched int                    `json:"expected_mismatched,omitempty"`
	Groups             []GroupSummary         `json:"groups,omitempty"`
	FailureFamilies    []FailureFamilySummary `json:"failure_families,omitempty"`
}

// GroupSummary contains counters for one fixture group.
type GroupSummary struct {
	Group              string `json:"group"`
	Total              int    `json:"total,omitempty"`
	Complete           int    `json:"complete,omitempty"`
	NeedsInput         int    `json:"needs_input,omitempty"`
	Failed             int    `json:"failed,omitempty"`
	Canceled           int    `json:"canceled,omitempty"`
	Skipped            int    `json:"skipped,omitempty"`
	ExpectedMatched    int    `json:"expected_matched,omitempty"`
	ExpectedMismatched int    `json:"expected_mismatched,omitempty"`
}

// FailureFamilySummary contains counters for one failure family.
type FailureFamilySummary struct {
	Family string `json:"family"`
	Total  int    `json:"total,omitempty"`
}

// NormalizeScorecard returns a deterministic scorecard with computed summary.
func NormalizeScorecard(scorecard Scorecard) Scorecard {
	scorecard.Version = norm.FirstNonEmpty(scorecard.Version, ScorecardVersion)
	scorecard.Name = norm.Trim(scorecard.Name)
	scorecard.Variants = NormalizeVariantResults(scorecard.Variants)
	scorecard.Summary = SummarizeVariants(scorecard.Variants)
	scorecard.Report = NormalizeReportMetadata(scorecard.Report)
	scorecard.Metadata = norm.Metadata(scorecard.Metadata)
	return scorecard
}

// CanonicalScorecardJSON returns deterministic indented JSON for scorecard.
func CanonicalScorecardJSON(scorecard Scorecard) ([]byte, error) {
	return json.MarshalIndent(NormalizeScorecard(scorecard), "", "  ")
}

// NormalizeFixtureResult returns a deterministic fixture result and propagates
// fixture-level identity and expected outcome to child variants when missing.
func NormalizeFixtureResult(fixture FixtureResult) FixtureResult {
	fixture.FixtureID = norm.Trim(fixture.FixtureID)
	fixture.Group = norm.Token(fixture.Group)
	fixture.ExpectedOutcome = NormalizeOutcome(fixture.ExpectedOutcome)
	for i := range fixture.Variants {
		if fixture.Variants[i].FixtureID == "" {
			fixture.Variants[i].FixtureID = fixture.FixtureID
		}
		if fixture.Variants[i].Group == "" {
			fixture.Variants[i].Group = fixture.Group
		}
		if fixture.Variants[i].ExpectedOutcome == "" {
			fixture.Variants[i].ExpectedOutcome = fixture.ExpectedOutcome
		}
	}
	fixture.Variants = NormalizeVariantResults(fixture.Variants)
	fixture.Metadata = norm.Metadata(fixture.Metadata)
	return fixture
}

// NormalizeVariantResults returns deterministic variant results.
func NormalizeVariantResults(results []VariantResult) []VariantResult {
	out := make([]VariantResult, 0, len(results))
	for _, result := range results {
		result = NormalizeVariantResult(result)
		if result.FixtureID == "" && result.VariantID == "" && result.ObservedOutcome == "" {
			continue
		}
		out = append(out, result)
	}
	slices.SortStableFunc(out, CompareVariantResult)
	return out
}

// NormalizeVariantResult returns a deterministic variant result.
func NormalizeVariantResult(result VariantResult) VariantResult {
	result.FixtureID = norm.Trim(result.FixtureID)
	result.VariantID = norm.Trim(result.VariantID)
	result.Group = norm.Token(result.Group)
	result.ExpectedOutcome = NormalizeOutcome(result.ExpectedOutcome)
	result.ObservedOutcome = NormalizeOutcome(result.ObservedOutcome)
	result.OutcomeMatched = CompareOutcome(result.ExpectedOutcome, result.ObservedOutcome)
	result.FailureFamily = normalizeFailureFamily(result.FailureFamily, result)
	result.Message = norm.Trim(result.Message)
	result.Diagnostics = trust.NormalizeDiagnostics(result.Diagnostics)
	result.Metadata = norm.Metadata(result.Metadata)
	return result
}

// CompareVariantResult orders variant results deterministically.
func CompareVariantResult(a, b VariantResult) int {
	a = NormalizeVariantResult(a)
	b = NormalizeVariantResult(b)
	return norm.CompareStrings(a.Group, b.Group, a.FixtureID, b.FixtureID, a.VariantID, b.VariantID, a.ObservedOutcome, b.ObservedOutcome)
}

// NormalizeOutcome maps common outcome strings to generic scorecard outcomes.
func NormalizeOutcome(outcome string) string {
	switch norm.Token(outcome) {
	case StatusComplete, "pass", "passed", "success", "succeeded":
		return StatusComplete
	case StatusNeedsInput, "need_input", "input_required":
		return StatusNeedsInput
	case StatusFailed, "fail", "failure", "error":
		return StatusFailed
	case StatusCanceled, "cancelled":
		return StatusCanceled
	case OutcomeSkipped, "skip":
		return OutcomeSkipped
	case "":
		return ""
	default:
		return norm.Token(outcome)
	}
}

// CompareOutcome reports whether an expected outcome is satisfied.
func CompareOutcome(expected, observed string) bool {
	expected = NormalizeOutcome(expected)
	observed = NormalizeOutcome(observed)
	return expected != "" && expected == observed
}

// SummarizeVariants returns stable grouped counters for variant results.
func SummarizeVariants(results []VariantResult) ScorecardSummary {
	results = NormalizeVariantResults(results)
	summary := ScorecardSummary{Total: len(results)}
	groups := map[string]*GroupSummary{}
	families := map[string]*FailureFamilySummary{}
	for _, result := range results {
		addOutcome(&summary, result.ObservedOutcome)
		if result.ExpectedOutcome != "" {
			if result.OutcomeMatched {
				summary.ExpectedMatched++
			} else {
				summary.ExpectedMismatched++
			}
		}
		group := norm.FirstNonEmpty(result.Group, "default")
		if groups[group] == nil {
			groups[group] = &GroupSummary{Group: group}
		}
		addGroupOutcome(groups[group], result)
		if result.FailureFamily != "" {
			if families[result.FailureFamily] == nil {
				families[result.FailureFamily] = &FailureFamilySummary{Family: result.FailureFamily}
			}
			families[result.FailureFamily].Total++
		}
	}
	summary.Groups = sortedGroupSummaries(groups)
	summary.FailureFamilies = sortedFailureFamilySummaries(families)
	return summary
}

// ValidateScorecard returns generic shape diagnostics for a scorecard.
func ValidateScorecard(scorecard Scorecard) []trust.DiagnosticRecord {
	scorecard = NormalizeScorecard(scorecard)
	var diagnostics []trust.DiagnosticRecord
	if len(scorecard.Variants) == 0 {
		diagnostics = append(diagnostics, trust.DiagnosticRecord{
			Code:     "scorecard.no_variants",
			Severity: "error",
			Message:  "scorecard has no variant results",
		})
	}
	for _, result := range scorecard.Variants {
		if result.FixtureID == "" {
			diagnostics = append(diagnostics, scorecardDiagnostic("scorecard.missing_fixture_id", result, "variant result is missing fixture id"))
		}
		if result.ObservedOutcome == "" {
			diagnostics = append(diagnostics, scorecardDiagnostic("scorecard.missing_observed_outcome", result, "variant result is missing observed outcome"))
		}
	}
	return trust.NormalizeDiagnostics(diagnostics)
}

func normalizeFailureFamily(family string, result VariantResult) string {
	family = norm.Token(family)
	if family != "" {
		return family
	}
	if result.ObservedOutcome == StatusFailed || (result.ExpectedOutcome != "" && !result.OutcomeMatched) {
		return FailureFamilyUnclassified
	}
	return ""
}

func addOutcome(summary *ScorecardSummary, outcome string) {
	switch outcome {
	case StatusComplete:
		summary.Complete++
	case StatusNeedsInput:
		summary.NeedsInput++
	case StatusCanceled:
		summary.Canceled++
	case OutcomeSkipped:
		summary.Skipped++
	case StatusFailed:
		summary.Failed++
	}
}

func addGroupOutcome(summary *GroupSummary, result VariantResult) {
	summary.Total++
	switch result.ObservedOutcome {
	case StatusComplete:
		summary.Complete++
	case StatusNeedsInput:
		summary.NeedsInput++
	case StatusCanceled:
		summary.Canceled++
	case OutcomeSkipped:
		summary.Skipped++
	case StatusFailed:
		summary.Failed++
	}
	if result.ExpectedOutcome != "" {
		if result.OutcomeMatched {
			summary.ExpectedMatched++
		} else {
			summary.ExpectedMismatched++
		}
	}
}

func sortedGroupSummaries(groups map[string]*GroupSummary) []GroupSummary {
	out := make([]GroupSummary, 0, len(groups))
	for _, group := range groups {
		out = append(out, *group)
	}
	slices.SortStableFunc(out, func(a, b GroupSummary) int {
		return norm.CompareStrings(a.Group, b.Group)
	})
	return out
}

func sortedFailureFamilySummaries(families map[string]*FailureFamilySummary) []FailureFamilySummary {
	out := make([]FailureFamilySummary, 0, len(families))
	for _, family := range families {
		out = append(out, *family)
	}
	slices.SortStableFunc(out, func(a, b FailureFamilySummary) int {
		return norm.CompareStrings(a.Family, b.Family)
	})
	return out
}

func scorecardDiagnostic(code string, result VariantResult, message string) trust.DiagnosticRecord {
	return trust.DiagnosticRecord{
		Code:     code,
		Severity: "error",
		Message:  message,
		Location: trust.DiagnosticLocation{
			Address: norm.FirstNonEmpty(result.FixtureID, result.VariantID),
		},
	}
}
