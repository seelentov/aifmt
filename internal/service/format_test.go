package service

import (
	"os"
	"strings"
	"testing"
)

func TestFormatCode(t *testing.T) {
	lg := `package main

	func main() {
		return "Hello, World!"
	}`

	token := os.Getenv("API_KEY")
	if token == "" {
		t.Fatal("API_KEY environment variable is not set")
	}

	fmtd, upds, err := FormatCode(lg, "go", "deepseek/deepseek-chat:free", token, false, "", nil)
	if err != nil {
		t.Fatalf("FormatCode failed: %v", err)
	}

	for _, cs := range []string{
		"func main()",
		"package",
		"Hello",
		"World",
		",",
		"!",
	} {
		if !strings.Contains(fmtd, cs) {
			t.Errorf("Expected %s but not found in formatted code", cs)
			t.Log(fmtd)
			t.Log(upds)
			return
		}
	}
}
