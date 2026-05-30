package norm

import "testing"

func TestToken(t *testing.T) {
	if got := Token(" Needs-Input  Now "); got != "needs_input_now" {
		t.Fatalf("Token = %q, want needs_input_now", got)
	}
}

func TestFirstNonEmpty(t *testing.T) {
	if got := FirstNonEmpty(" ", "  value  ", "fallback"); got != "value" {
		t.Fatalf("FirstNonEmpty = %q, want value", got)
	}
}

func TestMetadata(t *testing.T) {
	got := Metadata(map[string]string{" a ": " b ", "empty": " ", " ": "x"})
	if len(got) != 1 || got["a"] != "b" {
		t.Fatalf("Metadata = %#v, want trimmed non-empty entry", got)
	}
	if got := Metadata(map[string]string{" ": " "}); got != nil {
		t.Fatalf("Metadata empty = %#v, want nil", got)
	}
}

func TestSeverityRank(t *testing.T) {
	cases := map[string]int{
		"":          0,
		"error":     0,
		"blocking":  0,
		"unknown":   0,
		"warning":   1,
		"advisory":  2,
		"info":      3,
		"Needs-Fix": 0,
	}
	for in, want := range cases {
		if got := SeverityRank(in); got != want {
			t.Fatalf("SeverityRank(%q) = %d, want %d", in, got, want)
		}
	}
	if CompareSeverity("warning", "info") >= 0 {
		t.Fatalf("CompareSeverity did not sort warning before info")
	}
}
