package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

func helperCommand(command string, s ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, s...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestRunEchoCommand(t *testing.T) {
	execCommand = helperCommand
	defer func() { execCommand = exec.Command }()

	if err := RunCommand([]string{"echo", "hi"}, 1*time.Second); err != nil {
		t.Errorf("Expected nil error, got %#v", err)
	}
}

func TestRunSleep100Command(t *testing.T) {
	execCommand = helperCommand
	defer func() { execCommand = exec.Command }()

	if err := RunCommand([]string{"sleep", "100"}, 10*time.Second); err != nil {
		t.Errorf("Expected nil error, got %#v", err)
	}
}

func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "echo":
		iargs := []interface{}{}
		for _, s := range args {
			iargs = append(iargs, s)
		}
		fmt.Println(iargs...)
	case "sleep":
		time.Sleep(100 * time.Second)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %q\n", cmd)
		os.Exit(3)
	}
}
