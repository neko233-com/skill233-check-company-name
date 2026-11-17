package namecheck

import (
	"fmt"
	"strings"
)

// candidateStrategies generate alternative names from a seed. Each returns a
// (name, strategy) pair; empty names are skipped.
type candidateStrategy struct {
	name string
	fn   func(seed Seed) (string, string)
}

// Seed holds the inputs used to generate candidates.
type Seed struct {
	Brand       string
	Industry    string
	Region      string
	EnglishHint string
}

var industryWords = map[string][]string{
	"软件": {"科技", "信息", "智能", "数字", "网络"},
	"科技": {"智能", "数字", "信息", "未来", "云"},
	"教育": {"教育", "学堂", "学苑", "培优", "智教"},
	"医疗": {"健康", "生物", "医药", "康养", "精准"},
	"金融": {"资本", "投资", "财富", "金服", "信通"},
	"文化": {"文化", "传媒", "创意", "影视", "视界"},
	"贸易": {"商贸", "实业", "供应链", "优选", "环球"},
}

var relationWords = []string{"控股", "国际", "数字", "未来", "云", "创投", "智造"}

func industryAlternatives(industry string) []string {
	for key, words := range industryWords {
		for _, t := range strings.Split(industry, " ") {
			if strings.Contains(t, key) || strings.Contains(key, t) {
				return words
			}
		}
	}
	return []string{"科技", "信息", "智能", "数字", "未来"}
}

func brandFromName(name string) string {
	core := Core(Normalize(name))
	if core == "" {
		return Normalize(name)
	}
	return core
}

var strategies = []candidateStrategy{
	{"字号+行业同义词", func(s Seed) (string, string) {
		alts := industryAlternatives(s.Industry)
		if len(alts) == 0 {
			return "", ""
		}
		return s.Brand + alts[0], "字号+行业同义词"
	}},
	{"字号+地域", func(s Seed) (string, string) {
		if s.Region == "" {
			return "", ""
		}
		return fmt.Sprintf("%s（%s）", s.Brand, s.Region), "字号+地域"
	}},
	{"字号+关系词", func(s Seed) (string, string) {
		return s.Brand + relationWords[0], "字号+关系词"
	}},
	{"英文译名", func(s Seed) (string, string) {
		if s.EnglishHint == "" {
			return "", ""
		}
		return s.EnglishHint, "英文译名"
	}},
	{"缩写+业务", func(s Seed) (string, string) {
		alts := industryAlternatives(s.Industry)
		abbr := acronym(s.Brand)
		if abbr == "" || len(alts) < 2 {
			return "", ""
		}
		return abbr + " " + alts[1], "缩写+业务"
	}},
	{"字号+关系词2", func(s Seed) (string, string) {
		return s.Brand + relationWords[1], "字号+关系词"
	}},
	{"字号+行业同义词2", func(s Seed) (string, string) {
		alts := industryAlternatives(s.Industry)
		if len(alts) < 3 {
			return "", ""
		}
		return s.Brand + alts[2], "字号+行业同义词"
	}},
	{"字号+关系词3", func(s Seed) (string, string) {
		return s.Brand + relationWords[2], "字号+关系词"
	}},
	{"字号+未来", func(s Seed) (string, string) {
		return s.Brand + relationWords[3], "字号+关系词"
	}},
	{"字号+云", func(s Seed) (string, string) {
		return s.Brand + relationWords[4], "字号+关系词"
	}},
	{"字号+行业同义词4", func(s Seed) (string, string) {
		alts := industryAlternatives(s.Industry)
		if len(alts) < 5 {
			return "", ""
		}
		return s.Brand + alts[4], "字号+行业同义词"
	}},
}

// acronym builds a short Latin acronym from a seed brand. For pure-CJK brands
// it returns an empty string; callers that want Latin names should supply an
// EnglishHint. For mixed/Latin brands it takes upper-case letters or CJK initials.
func acronym(brand string) string {
	var b strings.Builder
	for _, r := range brand {
		switch {
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= 'a' && r <= 'z':
			b.WriteRune(r - 32)
		}
	}
	if b.Len() < 2 {
		return ""
	}
	return b.String()
}

// GenerateCandidates produces up to count distinct alternatives for req. The
// original name and any name that is an exact duplicate of a supplied registry
// fact are excluded.
func GenerateCandidates(req Request) []Candidate {
	count := req.CandidateCount
	if count <= 0 {
		count = DefaultCandidateCount
	}
	seed := Seed{
		Brand:       brandFromName(req.Name),
		Industry:    req.Industry,
		Region:      req.Region,
		EnglishHint: req.EnglishHint,
	}
	blocked := map[string]bool{Normalize(req.Name): true}
	for _, e := range req.ExistingNames {
		blocked[Normalize(e.Name)] = true
	}

	results := make([]Candidate, 0, count)
	seen := map[string]bool{}
	for _, st := range strategies {
		if len(results) >= count {
			break
		}
		name, strategy := st.fn(seed)
		if name == "" {
			continue
		}
		n := Normalize(name)
		if n == "" || seen[n] || blocked[n] {
			continue
		}
		seen[n] = true
		results = append(results, Candidate{
			Name:       name,
			Strategy:   strategy,
			Assessment: Assess(name, req),
		})
	}
	return results
}
