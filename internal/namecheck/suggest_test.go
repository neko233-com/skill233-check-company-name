package namecheck

import "testing"

func TestGenerateCandidatesDefaultFive(t *testing.T) {
	req := Request{
		Name: "星澜科技", Industry: "软件开发", Region: "上海", EnglishHint: "Xinglan",
	}
	got := GenerateCandidates(req)
	if len(got) != DefaultCandidateCount {
		t.Fatalf("got %d candidates, want %d: %+v", len(got), DefaultCandidateCount, got)
	}
	seen := map[string]bool{}
	for _, c := range got {
		n := Normalize(c.Name)
		if seen[n] {
			t.Errorf("duplicate candidate %q", c.Name)
		}
		seen[n] = true
		if c.Name == req.Name {
			t.Errorf("candidate equals original name")
		}
		if c.Strategy == "" {
			t.Errorf("candidate %q has empty strategy", c.Name)
		}
	}
}

func TestGenerateCandidatesCustomCount(t *testing.T) {
	req := Request{Name: "星澜科技", Industry: "软件开发", Region: "上海", EnglishHint: "Xinglan", CandidateCount: 3}
	if got := GenerateCandidates(req); len(got) != 3 {
		t.Fatalf("got %d candidates, want 3", len(got))
	}
}

func TestGenerateCandidatesExcludesBlocked(t *testing.T) {
	req := Request{
		Name: "星澜科技", Industry: "软件开发", CandidateCount: 5,
		ExistingNames: []ExistingName{{Name: "星澜控股"}},
	}
	for _, c := range GenerateCandidates(req) {
		if Normalize(c.Name) == Normalize("星澜控股") {
			t.Errorf("blocked name generated: %q", c.Name)
		}
	}
}

func TestGenerateCandidatesNoEnglishHintStillFills(t *testing.T) {
	req := Request{Name: "星澜科技", Industry: "软件开发", Region: "北京"}
	if got := GenerateCandidates(req); len(got) != DefaultCandidateCount {
		t.Fatalf("without EnglishHint got %d candidates, want %d", len(got), DefaultCandidateCount)
	}
}
