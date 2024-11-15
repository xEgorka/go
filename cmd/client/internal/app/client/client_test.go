package client

import (
	_ "embed"
	"os"
	"os/exec"
	"testing"
)

func TestStart(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		Start()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestStart")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}
