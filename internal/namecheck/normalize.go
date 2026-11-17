package namecheck

import (
	"strings"
	"unicode"
)

// removableSuffixes covers corporate forms plus generic industry words. A
// Chinese company name is conventionally [region][字号][industry][form]; only
// the 字号 is distinctive, so all of the others must be stripped for collision
// detection. Order is irrelevant because removal repeats until stable.
var removableSuffixes = []string{
	// corporate forms
	"股份有限公司", "有限责任公司", "有限公司", "合伙企业", "个体工商户",
	"控股", "集团", "limited", "corporation", "inc", "ltd", "llc", "corp", "co",
	// Chinese industry words
	"软件开发", "信息技术", "科技", "技术", "信息", "智能", "数字", "网络", "软件",
	"教育", "培训", "医疗", "健康", "生物", "医药", "金融", "资本", "投资", "基金",
	"财务", "文化", "传媒", "影视", "创意", "设计", "贸易", "商贸", "电商", "实业",
	"供应链", "餐饮", "美食", "服务", "咨询",
	// English industry words
	"technology", "software", "digital", "information", "smart", "tech",
}

// Normalize canonicalizes a name for comparison:
//   - full-width ASCII -> half-width, full-width space -> ASCII space
//   - Unicode lower-casing
//   - drop whitespace and punctuation, but keep ASCII '-' inside Latin names
func Normalize(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 0xFF01 && r <= 0xFF5E: // full-width ASCII
			r = r - 0xFEE0
		case r == 0x3000: // ideographic space
			r = ' '
		}
		if r == '-' || unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

// Core removes corporate forms and generic industry words from a normalized
// name, leaving the distinctive 字号. If nothing remains, the normalized name
// is returned unchanged.
func Core(normalized string) string {
	s := normalized
	changed := true
	for changed {
		changed = false
		for _, suffix := range removableSuffixes {
			if len(s) > len(suffix) && strings.HasSuffix(s, suffix) {
				s = strings.TrimSuffix(s, suffix)
				changed = true
				break
			}
		}
	}
	if s == "" {
		return normalized
	}
	return s
}

// levenshtein returns the edit distance between two rune slices.
func levenshtein(a, b []rune) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	prev := make([]int, len(b)+1)
	cur := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			cur[j] = min3(cur[j-1]+1, prev[j]+1, prev[j-1]+cost)
		}
		prev, cur = cur, prev
	}
	return prev[len(b)]
}

func min3(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}

// Similarity returns a value in [0,1]. Plain Levenshtein over-penalizes a short
// name that is fully contained in a longer one (e.g. "星澜科技" inside
// "星澜科技（上海）"), so a containment bonus is applied for that common case.
func Similarity(a, b string) float64 {
	na, nb := Normalize(a), Normalize(b)
	if na == "" && nb == "" {
		return 1
	}
	if na == "" || nb == "" {
		return 0
	}
	if na == nb {
		return 1
	}
	short, long := na, nb
	if len([]rune(long)) < len([]rune(short)) {
		short, long = long, short
	}
	if len([]rune(short)) >= 2 && strings.Contains(long, short) {
		sr, lr := float64(len([]rune(short))), float64(len([]rune(long)))
		return 0.5 + 0.5*sr/lr
	}
	ar, br := []rune(na), []rune(nb)
	dist := levenshtein(ar, br)
	longer := len(ar)
	if len(br) > longer {
		longer = len(br)
	}
	return 1 - float64(dist)/float64(longer)
}
