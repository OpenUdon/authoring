package icot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/OpenUdon/authoring/lifecycle"
	"github.com/OpenUdon/authoring/prompt"
	readinesspkg "github.com/OpenUdon/authoring/readiness"
	sessionpkg "github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/transcript"
)

// PromptTurn records one local prompt and answer.
type PromptTurn = sessionpkg.PromptTurn

// Event records a structured interactive-loop event. The payload is owned by
// the downstream product. Kind is retained for compatibility; Type is the
// durable transcript event name.
type Event struct {
	Kind string `json:"kind"`
	Type string `json:"type,omitempty"`
	Data any    `json:"data,omitempty"`
}

// PromptTranscript is a local transcript envelope for replay and review.
type PromptTranscript struct {
	Version string       `json:"version"`
	TimeUTC string       `json:"time_utc"`
	Turns   []PromptTurn `json:"turns"`
	Events  []Event      `json:"events,omitempty"`
	Session any          `json:"session,omitempty"`
}

// PromptSession prompts on a reader/writer pair and records prompt turns.
type PromptSession struct {
	session *prompt.Session
}

// NewPromptSession creates a local prompt session.
func NewPromptSession(in io.Reader, out io.Writer) *PromptSession {
	return &PromptSession{session: prompt.NewSession(in, out)}
}

// SetDefaultMode controls how defaulted prompts are handled.
func (session *PromptSession) SetDefaultMode(mode prompt.DefaultMode) {
	if session == nil {
		return
	}
	session.session.SetDefaultMode(mode)
}

// Ask prompts for a required free-form value.
func (session *PromptSession) Ask(label string) (string, error) {
	return session.session.Ask(label)
}

// AskDefault prompts for a value, returning current when the answer is blank.
func (session *PromptSession) AskDefault(label, current string) (string, error) {
	return session.session.AskDefault(label, current)
}

// AskDefaultForced prints a defaulted prompt and waits for user input.
func (session *PromptSession) AskDefaultForced(label, current string) (string, error) {
	return session.session.AskDefaultForced(label, current)
}

// AskOptionalDefault prompts for an optional value with a default.
func (session *PromptSession) AskOptionalDefault(label, current string) (string, error) {
	return session.session.AskOptionalDefault(label, current)
}

// AskDefaultRequired prompts until a non-empty value is available.
func (session *PromptSession) AskDefaultRequired(label, current string) (string, error) {
	return session.session.AskDefaultRequired(label, current)
}

// AskYesNo prompts for a yes/no answer with a default.
func (session *PromptSession) AskYesNo(label string, defaultYes bool) (bool, error) {
	return session.session.AskYesNo(label, defaultYes)
}

// Turns returns a copy of recorded prompt turns.
func (session *PromptSession) Turns() []PromptTurn {
	if session == nil {
		return nil
	}
	return fromPromptTurns(session.session.Turns())
}

func fromPromptTurns(turns []sessionpkg.PromptTurn) []PromptTurn {
	out := make([]PromptTurn, 0, len(turns))
	for _, turn := range turns {
		out = append(out, turn)
	}
	return out
}

// OneLine normalizes a prompt default for display.
func OneLine(value string) string {
	return prompt.OneLine(value)
}

// AssertPromptLabelsInOrder verifies that prompt labels were emitted in replay
// order.
func AssertPromptLabelsInOrder(output string, turns []PromptTurn) error {
	offset := 0
	for _, turn := range turns {
		index := strings.Index(output[offset:], turn.Label)
		if index < 0 {
			return fmt.Errorf("prompt label %q not found after offset %d", turn.Label, offset)
		}
		offset += index + len(turn.Label)
	}
	return nil
}

// SavePromptTranscript writes a prompt transcript with private-file
// permissions. Empty paths are ignored.
func SavePromptTranscript(path, version string, turns []PromptTurn, events []Event, session any) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	if strings.TrimSpace(version) == "" {
		version = "authoring.icot-transcript.v1"
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	transcript := PromptTranscript{
		Version: version,
		TimeUTC: time.Now().UTC().Format(time.RFC3339),
		Turns:   append([]PromptTurn(nil), turns...),
		Events:  NormalizeEvents(events),
		Session: session,
	}
	data, err := json.MarshalIndent(transcript, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return lifecycle.AtomicWrite(path, data, 0o600)
}

// NormalizeEvents returns a compatibility-preserving copy of product payload
// events with both Kind and Type populated when either is present.
func NormalizeEvents(events []Event) []Event {
	out := make([]Event, 0, len(events))
	for _, event := range events {
		event.Kind = strings.TrimSpace(event.Kind)
		event.Type = strings.TrimSpace(event.Type)
		if event.Type == "" {
			event.Type = event.Kind
		}
		if event.Kind == "" {
			event.Kind = event.Type
		}
		if event.Kind == "" && event.Type == "" && event.Data == nil {
			continue
		}
		out = append(out, event)
	}
	return out
}

// TranscriptEvent projects a product-payload event into the durable transcript
// event contract without interpreting downstream payload fields.
func TranscriptEvent(event Event) transcript.Event {
	events := NormalizeEvents([]Event{event})
	if len(events) == 0 {
		return transcript.Event{}
	}
	event = events[0]
	out := transcript.Event{
		Type:   event.Type,
		Fields: map[string]string{},
	}
	if event.Kind != "" && event.Kind != event.Type {
		out.Fields["kind"] = event.Kind
	}
	if event.Data != nil {
		if data, err := json.Marshal(event.Data); err == nil && len(data) > 0 {
			out.Fields["data_json"] = string(data)
		}
	}
	if len(out.Fields) == 0 {
		out.Fields = nil
	}
	return out
}

// TranscriptEvents projects product-payload events into durable transcript
// events.
func TranscriptEvents(events []Event) []transcript.Event {
	normalized := NormalizeEvents(events)
	out := make([]transcript.Event, 0, len(normalized))
	for _, event := range normalized {
		projected := TranscriptEvent(event)
		if projected.Type != "" {
			out = append(out, projected)
		}
	}
	return out
}

// ReadinessIssue explains why an interactive draft needs more input before it
// can be saved by the downstream adapter.
type ReadinessIssue = sessionpkg.ReadinessIssue

// InteractiveQuestion is one next-question decision in an interactive loop.
type InteractiveQuestion = readinesspkg.Question

// DraftRequest is the model-facing input for an interactive draft.
type DraftRequest[S, D any] struct {
	Opening           string           `json:"opening"`
	Brief             string           `json:"brief,omitempty"`
	Session           S                `json:"session"`
	Docs              []D              `json:"docs"`
	TranscriptTurns   []PromptTurn     `json:"transcript_turns,omitempty"`
	ReadinessFeedback []ReadinessIssue `json:"readiness_feedback,omitempty"`
}

// Extractor provides optional AI assistance for an interactive loop.
type Extractor[S, D any] interface {
	Kickoff(context.Context, string) (S, error)
	Draft(context.Context, DraftRequest[S, D]) (S, error)
	Refine(context.Context, S) (S, error)
	Disambiguate(context.Context, string, []D) ([]string, error)
}

// NoopExtractor disables AI assistance.
type NoopExtractor[S, D any] struct{}

func (NoopExtractor[S, D]) noopInteractiveExtractor() {}

func (NoopExtractor[S, D]) Kickoff(context.Context, string) (S, error) {
	var zero S
	return zero, nil
}

func (NoopExtractor[S, D]) Draft(context.Context, DraftRequest[S, D]) (S, error) {
	var zero S
	return zero, nil
}

func (NoopExtractor[S, D]) Refine(_ context.Context, session S) (S, error) {
	return session, nil
}

func (NoopExtractor[S, D]) Disambiguate(context.Context, string, []D) ([]string, error) {
	return nil, nil
}

// InteractiveHooks supplies product-specific behavior for the generic loop.
type InteractiveHooks[S, D, A any] struct {
	Session       S
	Documents     []D
	Opening       string
	Brief         string
	NoLLM         bool
	DefaultMode   prompt.DefaultMode
	MaxAttempts   int
	OpeningPrompt string

	Extractor Extractor[S, D]

	Normalize            func(*S)
	ApplyOpeningAnswer   func(*S, string, []D) error
	OpeningEvents        func(S) []Event
	Autosave             func(S) error
	RankDocuments        func([]D, []string) []D
	DeterministicPrefill func(*S, []D) bool
	LooksLikeSession     func(S) bool
	MergeDraft           func(S, S, []D) S
	AfterDraft           func(S) error
	DraftResultSummary   func(S) any
	DraftEvents          func(S) []Event
	OnDraftError         func(error)
	RefreshDocuments     func(S, []D) ([]D, error)
	ShouldDraft          func(S, []D, []ReadinessIssue) bool
	ShouldDraftQuestion  func(S, []D, []ReadinessIssue, InteractiveQuestion) bool
	DraftQuestion        func(context.Context, *S, []D, []ReadinessIssue, InteractiveQuestion) (bool, error)
	CheckReadiness       func(S, []D) []ReadinessIssue
	Ready                func(S, []ReadinessIssue) bool
	PlanFrontier         func(S, []D, []ReadinessIssue) []InteractiveQuestion
	ApplyRound           func(*S, []RoundAnswer, []D) error
	// Deprecated: adapted to a one-question frontier for source compatibility.
	PlanQuestion       func(S, []D, []ReadinessIssue) InteractiveQuestion
	ApplyAnswer        func(*S, InteractiveQuestion, string, []D) error
	FinalConfirm       func(*PromptSession, *S, []D, *[]Event) (A, error)
	FinalResultSummary func(A) any
	SaveTranscript     func([]PromptTurn, []Event, A) error
}

// InteractiveLifecycleOptions adds draft/transcript lifecycle behavior around
// RunInteractive. Draft persistence functions are supplied by the downstream
// product so Authoring does not own artifact formats.
type InteractiveLifecycleOptions[S, D, A any] struct {
	DraftPath            string
	TranscriptPath       string
	TranscriptVersion    string
	DeleteDraftOnSuccess bool
	Normalize            func(*S)
	LooksLikeSession     func(S) bool
	Opening              func(S) string
	TranscriptSession    func(A) any
	LoadDraft            func(string) (S, bool, error)
	SaveDraft            func(string, S) error
	DeleteDraft          func(string) error
}

// RunInteractiveWithLifecycle binds interactive hooks to caller-owned draft,
// autosave, transcript, and cleanup lifecycle behavior.
func RunInteractiveWithLifecycle[S, D, A any](ctx context.Context, in io.Reader, out io.Writer, hooks InteractiveHooks[S, D, A], opts InteractiveLifecycleOptions[S, D, A]) (A, error) {
	draftPath := strings.TrimSpace(opts.DraftPath)
	session := hooks.Session
	if draftPath != "" && opts.LoadDraft != nil {
		loaded, ok, err := opts.LoadDraft(draftPath)
		if err != nil {
			var zero A
			return zero, err
		}
		if ok && (opts.LooksLikeSession == nil || opts.LooksLikeSession(loaded)) {
			session = loaded
		}
	}
	if opts.Normalize != nil {
		opts.Normalize(&session)
	}
	hooks.Session = session
	if strings.TrimSpace(hooks.Opening) == "" && opts.Opening != nil {
		hooks.Opening = opts.Opening(session)
	}

	baseAutosave := hooks.Autosave
	hooks.Autosave = func(session S) error {
		if baseAutosave != nil {
			if err := baseAutosave(session); err != nil {
				return err
			}
		}
		if draftPath == "" || opts.SaveDraft == nil {
			return nil
		}
		if opts.LooksLikeSession != nil && !opts.LooksLikeSession(session) {
			return nil
		}
		if opts.Normalize != nil {
			opts.Normalize(&session)
		}
		return opts.SaveDraft(draftPath, session)
	}

	baseTranscript := hooks.SaveTranscript
	hooks.SaveTranscript = func(turns []PromptTurn, events []Event, artifacts A) error {
		if baseTranscript != nil {
			if err := baseTranscript(turns, events, artifacts); err != nil {
				return err
			}
		}
		if strings.TrimSpace(opts.TranscriptPath) == "" {
			return nil
		}
		var transcriptSession any = artifacts
		if opts.TranscriptSession != nil {
			transcriptSession = opts.TranscriptSession(artifacts)
		}
		return SavePromptTranscript(opts.TranscriptPath, opts.TranscriptVersion, turns, events, transcriptSession)
	}

	artifacts, err := RunInteractive(ctx, in, out, hooks)
	if err == nil && opts.DeleteDraftOnSuccess && opts.DeleteDraft != nil {
		err = opts.DeleteDraft(draftPath)
	}
	return artifacts, err
}
