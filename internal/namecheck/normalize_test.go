package namecheck

import "testing"

func TestNormalize(t *testing.T) {
	cases := []struct{ in, want string }{
		{"星澜科技", "星澜科技"},
		{"Star-Lane Tech", "star-lanetech"},
		{"  ＡＢＣ　１２３ ", "abc123"},
		{"星澜（上海）科技有限公司", "星澜上海科技有限公司"},
		{"Foo_Bar. Co", "foobarco"},
	}
	for _, c := range cases {
		if got := Normalize(c.in); got != c.want {
			t.Errorf("Normalize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestCore(t *testing.T) {
	cases := []struct{ in, want string }{
		{"星澜科技有限公司", "星澜"},
		{"北京星澜科技有限责任公司", "北京星澜"},
		{"Xinglan Technology Limited", "xinglan"},
		{"星澜", "星澜"},
	}
	for _, c := range cases {
		if got := Core(Normalize(c.in)); got != c.want {
			t.Errorf("Core(Normalize(%q)) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestSimilarity(t *testing.T) {
	if got := Similarity("星澜科技", "星澜科技"); got != 1 {
		t.Errorf("identical similarity = %v, want 1", got)
	}
	if got := Similarity("星澜科技", "星澜科技（上海）"); got < SimilarityThreshold {
		t.Errorf("near-identical similarity = %v, want >= %v", got, SimilarityThreshold)
	}
	if got := Similarity("星澜科技", "完全不同"); got >= SimilarityThreshold {
		t.Errorf("different similarity = %v, want < %v", got, SimilarityThreshold)
	}
}
