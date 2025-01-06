package workflow_service

import (
	"testing"
)

func TestconcatanateSlashToPath(t *testing.T) {
	name := "test.go"
	filepath := "path"
	want := "/" + filepath + "/" + name

	answer := concatanateSlashToPath(filepath, name)
	if answer != want {
		t.Fatalf(`concatanateSlashToPath("path", "test.go") = %s, want match for %s, nil`, answer, want)
	}
}
