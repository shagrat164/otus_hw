package main

import (
	"os"
	"testing"
)

func TestReadDir(t *testing.T) {
	dir := "forTest/env"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll("forTest")

	_ = os.WriteFile(dir+"/FOO", []byte("123"), 0o644)
	_ = os.WriteFile(dir+"/BAR", []byte("value\x00nextline"), 0o644)
	_ = os.WriteFile(dir+"/EMPTY", []byte{}, 0o644)

	env, err := ReadDir(dir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if env["FOO"].Value != "123" || env["FOO"].NeedRemove {
		t.Errorf("FOO variable mismatch")
	}
	if env["BAR"].Value != "value\nnextline" || env["BAR"].NeedRemove {
		t.Errorf("BAR variable mismatch")
	}
	if !env["EMPTY"].NeedRemove {
		t.Errorf("EMPTY variable mismatch")
	}
}
