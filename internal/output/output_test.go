package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func TestResultMarshalJSON(t *testing.T) {
	// Test with no error
	result := Result{
		Site:   "TestSite",
		URL:    "https://example.com/user",
		Exists: true,
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal result: %v", err)
	}

	var unmarshaled map[string]any
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if unmarshaled["site"] != "TestSite" {
		t.Errorf("Expected site to be TestSite, got %v", unmarshaled["site"])
	}
	if unmarshaled["url"] != "https://example.com/user" {
		t.Errorf("Expected url to be https://example.com/user, got %v", unmarshaled["url"])
	}
	if unmarshaled["exists"] != true {
		t.Errorf("Expected exists to be true, got %v", unmarshaled["exists"])
	}
	if _, ok := unmarshaled["error"]; ok {
		t.Errorf("Expected error to be omitted, but it was included")
	}

	// Test with error
	result = Result{
		Site:   "TestSite",
		URL:    "https://example.com/user",
		Exists: false,
		Error:  errors.New("test error"),
	}

	data, err = json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal result: %v", err)
	}

	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	if unmarshaled["error"] != "test error" {
		t.Errorf("Expected error to be 'test error', got %v", unmarshaled["error"])
	}
}

func TestFormatterPrintResult(t *testing.T) {
	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a formatter and result
	formatter := NewFormatter(true).WithColor(false)
	result := Result{
		Site:   "TestSite",
		URL:    "https://example.com/user",
		Exists: true,
	}

	// Print the result
	formatter.PrintResult(result)

	// Restore stdout and get the output
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Check the output
	output := buf.String()
	if !strings.Contains(output, "TestSite") {
		t.Errorf("Expected output to contain TestSite, got %s", output)
	}
	if !strings.Contains(output, "https://example.com/user") {
		t.Errorf("Expected output to contain https://example.com/user, got %s", output)
	}
}

func TestSaveToFile(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "accio-test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Create a formatter and results
	formatter := NewFormatter(false)
	results := []Result{
		{
			Site:   "TestSite1",
			URL:    "https://example.com/user1",
			Exists: true,
		},
		{
			Site:   "TestSite2",
			URL:    "https://example.com/user2",
			Exists: false,
		},
	}

	// Save to file
	err = formatter.SaveToFile(results, tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to save to file: %v", err)
	}

	// Read the file
	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Check the content
	content := string(data)
	if !strings.Contains(content, "TestSite1") {
		t.Errorf("Expected file to contain TestSite1, got %s", content)
	}
	if !strings.Contains(content, "https://example.com/user1") {
		t.Errorf("Expected file to contain https://example.com/user1, got %s", content)
	}
}
