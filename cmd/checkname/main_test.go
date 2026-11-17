package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	in := strings.NewReader(`{"name":"星澜科技","industry":"软件开发","region":"上海","existingNames":[{"name":"星澜科技有限公司","jurisdiction":"CN-MAINLAND"}]}`)
	var out bytes.Buffer
	if err := run(in, &out); err != nil {
		t.Fatalf("run: %v", err)
	}
	var payload struct {
		Original struct {
			Name string `json:"name"`
		} `json:"original"`
		Candidates []struct {
			Name string `json:"name"`
		} `json:"candidates"`
		UnavailableReasons []struct {
			Code string `json:"code"`
		} `json:"unavailableReasons"`
	}
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("parse output: %v\n%s", err, out.String())
	}
	if payload.Original.Name != "星澜科技" {
		t.Fatalf("unexpected original: %+v", payload.Original)
	}
	if len(payload.Candidates) != 5 {
		t.Fatalf("want 5 candidates, got %d", len(payload.Candidates))
	}
	if len(payload.UnavailableReasons) == 0 {
		t.Fatal("expected unavailable reasons")
	}
}

func TestRunMissingName(t *testing.T) {
	var out bytes.Buffer
	if err := run(strings.NewReader(`{}`), &out); err == nil {
		t.Fatal("expected error for missing name")
	}
}
