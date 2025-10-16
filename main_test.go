package main

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
)

// Mock converter for testing
type MockConverter struct{}

func (m *MockConverter) Convert(input, output, order string) error {
	return nil
}

// Helper function to capture stdout
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// Test printUsage output
func TestPrintUsage(t *testing.T) {
	output := captureOutput(func() {
		printUsage()
	})

	tests := []string{
		"Image to PDF Converter",
		"Supported formats: JPG, JPEG, PNG, WEBP, TIFF",
		"Usage:",
		"Examples:",
		"./img2pdf -i images/",
		"./img2pdf -i \"image1.jpg,photo.png,scan.tiff\" -o result.pdf",
		"./img2pdf -i \"images/,photo.jpg,scan.png\" -o result.pdf -order mod",
		"Note: The -i flag accepts both directories and individual files (comma-separated)",
		"Options:",
		"-i string",
		"-o string",
		"-order string",
		"-help",
		"Sorting order for images: seq (sequential), nam (by name), mod (by modification time)",
	}

	for _, test := range tests {
		if !strings.Contains(output, test) {
			t.Errorf("Expected output to contain %q, but it didn't", test)
		}
	}
}

// Test printInputHelp output
func TestPrintInputHelp(t *testing.T) {
	output := captureOutput(func() {
		printInputHelp()
	})

	tests := []string{
		"Usage:",
		"./img2pdf -i",
	}

	for _, test := range tests {
		if !strings.Contains(output, test) {
			t.Errorf("Expected output to contain %q, but it didn't", test)
		}
	}
}

// Test flag parsing
func TestFlagParsing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantHelp bool
	}{
		{
			name:     "help flag set",
			args:     []string{"prog", "-help"},
			wantHelp: true,
		},
		{
			name:     "help with input",
			args:     []string{"prog", "-help", "-i", "test.jpg"},
			wantHelp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag state
			flag.CommandLine = flag.NewFlagSet(tt.args[0], flag.ContinueOnError)

			input := flag.String("i", "", "Input directory or JPG files (space-separated)")
			order := flag.String("order", "seq", "default sequently order")
			output := flag.String("o", "output.pdf", "Output PDF file path")
			help := flag.Bool("help", false, "Show help")

			flag.CommandLine.Parse(tt.args[1:])

			if *help != tt.wantHelp {
				t.Errorf("help flag: got %v, want %v", *help, tt.wantHelp)
			}

			// Verify default values for other flags
			if *order != "seq" {
				t.Errorf("order default: got %s, want seq", *order)
			}
			if *output != "output.pdf" {
				t.Errorf("output default: got %s, want output.pdf", *output)
			}
			if *input != "" && tt.name != "help with input" {
				t.Errorf("input should be empty, got %s", *input)
			}
		})
	}
}

// Test custom flag values
func TestCustomFlagValues(t *testing.T) {
	// Reset flag state
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)

	input := flag.String("i", "", "Input directory or JPG files (space-separated)")
	order := flag.String("order", "seq", "default sequently order")
	output := flag.String("o", "output.pdf", "Output PDF file path")

	args := []string{"-i", "images/", "-o", "result.pdf", "-order", "mod"}
	flag.CommandLine.Parse(args)

	if *input != "images/" {
		t.Errorf("input: got %s, want images/", *input)
	}
	if *output != "result.pdf" {
		t.Errorf("output: got %s, want result.pdf", *output)
	}
	if *order != "mod" {
		t.Errorf("order: got %s, want mod", *order)
	}
}

// Test short vs long flag names
func TestFlagNamesVariations(t *testing.T) {
	// Reset flag state
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)

	input := flag.String("i", "", "Input directory or JPG files (space-separated)")
	output := flag.String("o", "output.pdf", "Output PDF file path")

	args := []string{"-i", "test.jpg", "-o", "out.pdf"}
	flag.CommandLine.Parse(args)

	if *input != "test.jpg" {
		t.Errorf("short flag -i: got %s, want test.jpg", *input)
	}
	if *output != "out.pdf" {
		t.Errorf("short flag -o: got %s, want out.pdf", *output)
	}
}
