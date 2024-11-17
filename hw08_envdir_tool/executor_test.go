package main

import "testing"

func TestRunCmd(t *testing.T) {
	env := Environment{
		"TEST_VAR": {Value: "123", NeedRemove: false},
	}

	cmd := []string{"bash", "-c", "echo $TEST_VAR"}
	code := RunCmd(cmd, env)
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}
