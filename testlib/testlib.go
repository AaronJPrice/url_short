package testlib

import (
	"path/filepath"
	"runtime"
	"testing"
)

func logPointOfFailure(t *testing.T) {
	_, filePath, line, _ := runtime.Caller(2)
	fileName := filepath.Base(filePath)
	t.Logf("Failure on %s:%v", fileName, line)
}

func Assert(t *testing.T, b bool) {
	if !b {
		logPointOfFailure(t)
		t.Fatalf("expression is not true")
	}
}

func AssertEqual(t *testing.T, have interface{}, want interface{}) {
	if have != want {
		logPointOfFailure(t)
		t.Fatalf("\nhave: %v \nwant: %v", have, want)
	}
}

func ErrCheck(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("%v", err)
	}
}
