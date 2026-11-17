package namecheck

import (
	"fmt"
	"strings"
)

// reservedWords are terms that, in mainland China, generally require special
// approval to use in a company name. Their presence is a strong "reason" signal.
var reservedWords = []string{"中国", "中华", "国家", "人民", "国务院", "中央", "军委", "检察院", "法院"}

// industryRules is ordered (most specific first) so industry key resolution is
// deterministic — iterating a Go map would randomly pick "科技" vs "软件".
var industryRules = []struct {
	key    string
	tokens []string
}{
	{"软件", []string{"软件开发", "信息技术", "互联网", "软件"}},
	{"教育", []string{"教育", "培训", "学堂"}},
	{"医疗", []string{"医疗", "健康", "生物", "医药"}},
	{"金融", []string{"金融", "资本", "投资", "基金", "财务"}},
	{"文化", []string{"文化", "传媒", "影视", "创意", "设计"}},
	{"贸易", []string{"贸易", "商贸", "电商", "实业"}},
	{"科技", []string{"科技", "技术", "信息", "智能", "数字"}},
}

// industryMatches reports whether two free-form industry strings share a token.
func industryMatches(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	keyOf := func(s string) string {
		for _, rule := range industryRules {
			for _, t := range rule.tokens {
				if strings.Contains(s, t) {
					return rule.key
				}
			}
		}
		return ""
	}
	ka, kb := keyOf(a), keyOf(b)
	if ka != "" && kb != "" {
		return ka == kb
	}
	return strings.Contains(a, b) || strings.Contains(b, a)
}

// Assess evaluates name against the registry facts in req and returns a scored
// assessment with human-readable reasons for any problem found.
func Assess(name string, req Request) Assessment {
	normalized := Normalize(name)
	core := Core(normalized)
	a := Assessment{Name: name, Available: true, Risk: RiskLow, Reasons: []Reason{}}

	// Reserved-word rule (China first).
	if req.Jurisdiction == "" || req.Jurisdiction == JurisdictionCNMainland {
		for _, w := range reservedWords {
			if strings.Contains(normalized, w) {
				a.Reasons = append(a.Reasons, Reason{
					Code:         ReasonReservedWord,
					Message:      fmt.Sprintf("名称含需特别核准的字样「%s」，在中国大陆通常不能自主申报使用。", w),
					Jurisdiction: JurisdictionCNMainland,
				})
			}
		}
	}

	best := 0.0
	for _, existing := range req.ExistingNames {
		sim := Similarity(name, existing.Name)
		if sim > best {
			best = sim
		}
		n := Normalize(existing.Name)
		c := Core(n)
		exact := n == normalized
		coreEqual := c != "" && c == core
		sameOrUnknownIndustry := existing.Industry == "" || req.Industry == "" || industryMatches(existing.Industry, req.Industry)

		switch {
		case exact:
			a.Reasons = append(a.Reasons, Reason{
				Code:         ReasonExactDuplicate,
				Message:      fmt.Sprintf("与已登记企业「%s」名称完全相同，构成重名。", existing.Name),
				Jurisdiction: existing.Jurisdiction,
				Source:       existing.Source,
			})
		case coreEqual && sameOrUnknownIndustry:
			a.Reasons = append(a.Reasons, Reason{
				Code:         ReasonHighSimilarity,
				Message:      fmt.Sprintf("字号「%s」与已登记企业「%s」相同且行业相近，核准时大概率被判近似。", core, existing.Name),
				Jurisdiction: existing.Jurisdiction,
				Source:       existing.Source,
			})
		case sim >= SimilarityThreshold:
			a.Reasons = append(a.Reasons, Reason{
				Code:         ReasonHighSimilarity,
				Message:      fmt.Sprintf("与已登记企业「%s」名称高度近似（相似度 %.0f%%）。", existing.Name, sim*100),
				Jurisdiction: existing.Jurisdiction,
				Source:       existing.Source,
			})
		}
	}

	for _, tm := range req.ExistingTrademarks {
		sim := Similarity(name, tm.Mark)
		if sim > best {
			best = sim
		}
		if (Core(Normalize(tm.Mark)) == core && core != "") || sim >= SimilarityThreshold {
			class := tm.Class
			if class == "" {
				class = "相关类别"
			}
			a.Reasons = append(a.Reasons, Reason{
				Code:         ReasonTrademarkConflict,
				Message:      fmt.Sprintf("与第 %s 类商标「%s」相同或近似，存在商标侵权/异议风险。", class, tm.Mark),
				Jurisdiction: tm.Jurisdiction,
				Source:       tm.Source,
			})
		}
	}

	a.Similarity = best

	hasExact := false
	hasStrong := len(a.Reasons) > 0
	for _, r := range a.Reasons {
		if r.Code == ReasonExactDuplicate {
			hasExact = true
		}
	}
	switch {
	case hasExact:
		a.Risk, a.Available = RiskHigh, false
	case hasStrong:
		a.Risk, a.Available = RiskMedium, false
	default:
		a.Risk, a.Available = RiskLow, true
	}
	return a
}
