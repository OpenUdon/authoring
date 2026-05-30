package authoring_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/OpenUdon/authoring/lifecycle"
	"github.com/OpenUdon/authoring/prompt"
	"github.com/OpenUdon/authoring/promptcontext"
	"github.com/OpenUdon/authoring/report"
	"github.com/OpenUdon/authoring/session"
	"github.com/OpenUdon/authoring/transcript"
)

func TestDurableContractVersions(t *testing.T) {
	versions := map[string]string{
		"lifecycle.DraftVersion":   lifecycle.DraftVersion,
		"prompt.TranscriptVersion": prompt.TranscriptVersion,
		"promptcontext.Version":    promptcontext.Version,
		"report.MetadataVersion":   report.MetadataVersion,
		"report.ScorecardVersion":  report.ScorecardVersion,
		"report.Version":           report.Version,
		"session.Version":          session.Version,
		"transcript.Version":       transcript.Version,
	}
	for name, version := range versions {
		if !strings.HasPrefix(version, "authoring.") || !strings.HasSuffix(version, ".v1") {
			t.Fatalf("%s = %q, want authoring.*.v1", name, version)
		}
	}
}

func TestDurableContractJSONTags(t *testing.T) {
	types := []reflect.Type{
		reflect.TypeOf(lifecycle.Draft[map[string]string]{}),
		reflect.TypeOf(prompt.PromptTranscript{}),
		reflect.TypeOf(prompt.ReplayScript{}),
		reflect.TypeOf(promptcontext.Context{}),
		reflect.TypeOf(report.ReportMetadata{}),
		reflect.TypeOf(report.Result{}),
		reflect.TypeOf(report.Scorecard{}),
		reflect.TypeOf(session.State{}),
		reflect.TypeOf(transcript.Record{}),
	}
	for _, typ := range types {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if field.PkgPath != "" {
				continue
			}
			tag := field.Tag.Get("json")
			if tag == "" || tag == "-" {
				t.Fatalf("%s.%s is exported without a durable json tag", typ.Name(), field.Name)
			}
		}
	}
}
