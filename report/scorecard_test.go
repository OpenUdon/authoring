package report

import (
	"strings"
	"testing"
)

func TestNormalizeFixtureAndVariants(t *testing.T) {
	fixture := NormalizeFixtureResult(FixtureResult{
		FixtureID:       " fixture-a ",
		Group:           " Core ",
		ExpectedOutcome: "pass",
		Variants: []VariantResult{
			{VariantID: "v1", ObservedOutcome: "success"},
			{VariantID: "v2", ObservedOutcome: "failed", FailureFamily: " Schema "},
		},
	})
	if len(fixture.Variants) != 2 {
		t.Fatalf("fixture = %#v, want two variants", fixture)
	}
	if fixture.Variants[0].FixtureID != "fixture-a" || fixture.Variants[0].Group != "core" || fixture.Variants[0].ExpectedOutcome != StatusComplete {
		t.Fatalf("variant = %#v, want propagated fixture metadata", fixture.Variants[0])
	}
	if !fixture.Variants[0].OutcomeMatched || fixture.Variants[1].OutcomeMatched {
		t.Fatalf("variants = %#v, want comparison outcomes", fixture.Variants)
	}
}

func TestScorecardSummaryGroupsAndFailures(t *testing.T) {
	scorecard := NormalizeScorecard(Scorecard{
		Name: " iCoT ",
		Variants: []VariantResult{
			{FixtureID: "b", VariantID: "base", Group: "extended", ExpectedOutcome: "needs-input", ObservedOutcome: "needs input"},
			{FixtureID: "a", VariantID: "base", Group: "core", ExpectedOutcome: "complete", ObservedOutcome: "complete"},
			{FixtureID: "a", VariantID: "mutant", Group: "core", ExpectedOutcome: "complete", ObservedOutcome: "failed", FailureFamily: "schema"},
			{FixtureID: "c", VariantID: "skip", Group: "extended", ObservedOutcome: "skip"},
		},
	})
	if scorecard.Version != ScorecardVersion || scorecard.Name != "iCoT" {
		t.Fatalf("scorecard = %#v, want normalized scorecard", scorecard)
	}
	if scorecard.Summary.Total != 4 || scorecard.Summary.Complete != 1 || scorecard.Summary.NeedsInput != 1 || scorecard.Summary.Failed != 1 || scorecard.Summary.Skipped != 1 {
		t.Fatalf("summary = %#v, want outcome counters", scorecard.Summary)
	}
	if scorecard.Summary.ExpectedMatched != 2 || scorecard.Summary.ExpectedMismatched != 1 {
		t.Fatalf("summary = %#v, want expected comparison counters", scorecard.Summary)
	}
	if len(scorecard.Summary.Groups) != 2 || scorecard.Summary.Groups[0].Group != "core" || scorecard.Summary.Groups[0].Total != 2 {
		t.Fatalf("groups = %#v, want stable group summaries", scorecard.Summary.Groups)
	}
	if len(scorecard.Summary.FailureFamilies) != 1 || scorecard.Summary.FailureFamilies[0].Family != "schema" {
		t.Fatalf("families = %#v, want schema failure family", scorecard.Summary.FailureFamilies)
	}
}

func TestValidateScorecard(t *testing.T) {
	diagnostics := ValidateScorecard(Scorecard{})
	if len(diagnostics) != 1 || diagnostics[0].Code != "scorecard.no_variants" {
		t.Fatalf("diagnostics = %#v, want no-variants diagnostic", diagnostics)
	}
	diagnostics = ValidateScorecard(Scorecard{Variants: []VariantResult{{VariantID: "v1"}}})
	if len(diagnostics) != 2 {
		t.Fatalf("diagnostics = %#v, want missing fixture and observed outcome diagnostics", diagnostics)
	}
}

func TestCanonicalScorecardJSON(t *testing.T) {
	data, err := CanonicalScorecardJSON(Scorecard{Variants: []VariantResult{{FixtureID: "f", ObservedOutcome: "pass"}}})
	if err != nil {
		t.Fatalf("CanonicalScorecardJSON error = %v", err)
	}
	if !strings.Contains(string(data), `"version": "authoring.scorecard.v1"`) || !strings.Contains(string(data), `"observed_outcome": "complete"`) {
		t.Fatalf("json = %s", data)
	}
	if !strings.Contains(string(data), `"outcome_matched": false`) {
		t.Fatalf("json = %s, want explicit unmatched boolean when no expectation is set", data)
	}
}
