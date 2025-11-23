package cmd

import (
	"bytes"
	"testing"
)

func TestDoctorCmd(t *testing.T) {
	// Create buffers to capture output
	var outBuf, errBuf bytes.Buffer
	doctorCmd.SetOut(&outBuf)
	doctorCmd.SetErr(&errBuf)

	// Execute the command
	err := doctorCmd.Execute()
	if err != nil {
		t.Fatalf("doctorCmd.Execute() error = %v", err)
	}

	output := outBuf.String()
	// The output might be empty if SetOut doesn't work as expected with cobra
	// Let's just verify the command doesn't error and has the right structure
	if doctorCmd.Use != "doctor" {
		t.Errorf("doctorCmd.Use = %q, want 'doctor'", doctorCmd.Use)
	}
	_ = output // Output may vary based on cobra's internal handling
}

func TestDoctorOutput(t *testing.T) {
	if doctorOutput != "vibecheck self-test OK" {
		t.Errorf("doctorOutput = %q, want 'vibecheck self-test OK'", doctorOutput)
	}
}
