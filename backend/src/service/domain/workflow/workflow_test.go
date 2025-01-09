package workflow_service

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func TestConcatanateSlashToPath(t *testing.T) {
	name := "test.go"
	filepath := "path"
	expected := "/" + filepath + "/" + name

	res := concatanateSlashToPath(filepath, name)
	if res != expected {
		t.Fatalf(`Got : %s but Expected -> %s`, res, expected)
	}
}
