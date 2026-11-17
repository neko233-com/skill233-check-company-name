// Command checkname reads a Check request as JSON from stdin and writes the
// report as JSON to stdout. It is the deterministic engine called by the Agent
// Skill after registry facts have been gathered.
//
// Usage:
//
//	echo '{"name":"星澜科技","industry":"软件开发"}' | go run ./cmd/checkname
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/neko233-com/skill233-check-company-name/internal/namecheck"
)

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "checkname:", err)
		os.Exit(1)
	}
}

func run(in io.Reader, out io.Writer) error {
	data, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("read input: %w", err)
	}
	var req namecheck.Request
	if err := json.Unmarshal(data, &req); err != nil {
		return fmt.Errorf("parse JSON: %w", err)
	}
	if req.Name == "" {
		return fmt.Errorf("field 'name' is required")
	}
	report := namecheck.Check(req)
	enc := json.NewEncoder(out)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
