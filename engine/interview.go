package engine

import (
	"fmt"
	"strings"

	"github.com/OpenUdon/authoring/interview"
)

// InterviewBinding binds product state to Authoring's interview transaction.
// State, Clone, and Resolve are required. Question may customize the default
// node-to-question projection; Normalize and Validate apply product-owned state
// rules around the shared atomic settlement.
type InterviewBinding[S, D any] struct {
	State     func(*S) *interview.State
	Clone     func(S) (S, error)
	Prepare   func(*S, []D) error
	Question  func(S, []D, interview.Node) Question
	Resolve   func(*S, []D, interview.Node, RoundAnswer) (interview.Resolution, error)
	Normalize func(*S)
	Validate  func(S) error
}

// Plan returns the complete dependency-ready question frontier for state.
func (binding InterviewBinding[S, D]) Plan(state *S, docs []D) ([]Question, error) {
	if err := binding.validate(); err != nil {
		return nil, err
	}
	if state == nil {
		return nil, fmt.Errorf("icot interview state is required")
	}
	clone, err := binding.Clone(*state)
	if err != nil {
		return nil, err
	}
	if binding.Normalize != nil {
		binding.Normalize(&clone)
	}
	if binding.Prepare != nil {
		if err := binding.Prepare(&clone, append([]D(nil), docs...)); err != nil {
			return nil, err
		}
	}
	if binding.Normalize != nil {
		binding.Normalize(&clone)
	}
	interviewState := binding.State(&clone)
	if interviewState == nil {
		return nil, fmt.Errorf("icot interview binding returned nil state")
	}
	nodes, err := interview.Frontier(*interviewState)
	if err != nil {
		return nil, err
	}
	questions := make([]Question, 0, len(nodes))
	for _, node := range nodes {
		question := Question{
			ID: node.ID, Prompt: firstNonEmpty(node.Prompt, node.Title), Slots: []string{node.ID},
			Required: node.Required, Recommendation: node.Recommendation, Priority: node.Priority,
			Rationale: node.Rationale, EvidenceRefs: append([]string(nil), node.EvidenceRefs...),
		}
		if binding.Question != nil {
			question = binding.Question(clone, append([]D(nil), docs...), node)
		}
		// Node identity remains the transaction key even when the adapter
		// customizes every user-facing question field.
		question.ID = node.ID
		questions = append(questions, question)
	}
	*state = clone
	return PlanFrontier(questions).Questions, nil
}

// Apply clones product state, applies product mutations through Resolve,
// settles the exact displayed frontier through interview.ApplyRound, validates
// the result, and swaps state only after the whole transaction succeeds.
func (binding InterviewBinding[S, D]) Apply(state *S, answers []RoundAnswer, docs []D) error {
	if err := binding.validate(); err != nil {
		return err
	}
	if state == nil {
		return fmt.Errorf("icot interview state is required")
	}
	clone, err := binding.Clone(*state)
	if err != nil {
		return err
	}
	if binding.Normalize != nil {
		binding.Normalize(&clone)
	}
	interviewState := binding.State(&clone)
	if interviewState == nil {
		return fmt.Errorf("icot interview binding returned nil state")
	}
	frontier, err := interview.Frontier(*interviewState)
	if err != nil {
		return err
	}
	if len(answers) != len(frontier) {
		return fmt.Errorf("frontier round must answer all %d displayed questions; got %d", len(frontier), len(answers))
	}
	byID := make(map[string]interview.Node, len(frontier))
	for _, node := range frontier {
		byID[node.ID] = node
	}
	seen := map[string]bool{}
	resolutions := make([]interview.Resolution, 0, len(answers))
	for _, answer := range answers {
		node, ok := byID[strings.TrimSpace(answer.QuestionID)]
		if !ok {
			return fmt.Errorf("answer references unknown frontier question %q", answer.QuestionID)
		}
		if seen[node.ID] {
			return fmt.Errorf("frontier round contains multiple answers for question %q", node.ID)
		}
		seen[node.ID] = true
		resolution, err := binding.Resolve(&clone, append([]D(nil), docs...), node, answer)
		if err != nil {
			return err
		}
		if strings.TrimSpace(resolution.NodeID) == "" {
			resolution.NodeID = node.ID
		}
		resolutions = append(resolutions, resolution)
	}
	settled, err := interview.ApplyRound(*interviewState, resolutions)
	if err != nil {
		return err
	}
	*interviewState = settled
	if binding.Normalize != nil {
		binding.Normalize(&clone)
	}
	if current := binding.State(&clone); current == nil {
		return fmt.Errorf("icot interview binding returned nil state after normalization")
	} else if err := interview.Validate(*current); err != nil {
		return err
	}
	if binding.Validate != nil {
		if err := binding.Validate(clone); err != nil {
			return err
		}
	}
	*state = clone
	return nil
}

func (binding InterviewBinding[S, D]) validate() error {
	switch {
	case binding.State == nil:
		return fmt.Errorf("icot interview state callback is required")
	case binding.Clone == nil:
		return fmt.Errorf("icot interview clone callback is required")
	case binding.Resolve == nil:
		return fmt.Errorf("icot interview resolution callback is required")
	default:
		return nil
	}
}

func (binding InterviewBinding[S, D]) noProgressRounds(state *S) int {
	if binding.State == nil || state == nil {
		return 0
	}
	current := binding.State(state)
	if current == nil || current.NoProgressRounds < 0 {
		return 0
	}
	return current.NoProgressRounds
}

func (binding InterviewBinding[S, D]) setNoProgressRounds(state *S, rounds int) {
	if binding.State == nil || state == nil {
		return
	}
	current := binding.State(state)
	if current == nil {
		return
	}
	if rounds < 0 {
		rounds = 0
	}
	current.NoProgressRounds = rounds
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}
