// main_test.go -- because main wants 100% test coverage too.
//
// (Seriously? Seriously. Design for testing.)

package main

import (
	"bytes"
	"testing"
)

func TestMain(t *testing.T) {

	// Rig up the fake io:
	exited := -1
	var sout bytes.Buffer
	var serr bytes.Buffer
	exit = func(c int) { exited = c }
	stdout = &sout
	stderr = &serr
	args = []string{"program", "--version"}

	// Run it in the simplest form possible:
	main()

	// Check our results, just to be thorough:
	if exited != 0 {
		t.Fatal("nonzero exit for --version")
	}
	if out := serr.String(); out != "" {
		t.Fatalf("wrote to stderr: %s", out)
	}
	if out := sout.String(); out != "binsanity version 0.1.0\n" {
		t.Fatalf("wrote wrong output to stdout: %s", out)
	}

}
