package authoring_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestImportBoundaryExcludesProductModules(t *testing.T) {
	cmd := exec.Command("go", "list", "-deps", "./...")
	cmd.Dir = "."
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list -deps ./... failed: %v\n%s", err, stderr.String())
	}
	forbidden := []string{
		"github.com/OpenUdon/openudon",
		"github.com/OpenUdon/ramen",
		"github.com/OpenUdon/uws",
		"github.com/OpenUdon/apitools",
	}
	for _, dep := range strings.Fields(string(out)) {
		for _, prefix := range forbidden {
			if dep == prefix || strings.HasPrefix(dep, prefix+"/") {
				t.Fatalf("Authoring import boundary includes forbidden dependency %q", dep)
			}
		}
	}
}
