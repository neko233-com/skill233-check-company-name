package namecheck

import "testing"

func TestAssessExactDuplicate(t *testing.T) {
	req := Request{
		Name:         "星澜科技有限公司",
		Industry:     "软件开发",
		Jurisdiction: JurisdictionCNMainland,
		ExistingNames: []ExistingName{{
			Name: "星澜科技有限公司", Jurisdiction: JurisdictionCNMainland, Source: "gsxt",
		}},
	}
	a := Assess(req.Name, req)
	if a.Risk != RiskHigh || a.Available {
		t.Fatalf("want HIGH/risk unavailable, got %+v", a)
	}
	if len(a.Reasons) == 0 || a.Reasons[0].Code != ReasonExactDuplicate {
		t.Fatalf("want exact-duplicate reason, got %+v", a.Reasons)
	}
}

func TestAssessCoreCollisionSameIndustry(t *testing.T) {
	req := Request{
		Name:     "星澜智能科技",
		Industry: "软件开发",
		ExistingNames: []ExistingName{{
			Name: "星澜科技有限公司", Industry: "信息技术", Jurisdiction: JurisdictionCNMainland,
		}},
	}
	a := Assess(req.Name, req)
	if a.Available || a.Risk != RiskMedium {
		t.Fatalf("want MEDIUM/unavailable for same-core same-industry, got %+v", a)
	}
	found := false
	for _, r := range a.Reasons {
		if r.Code == ReasonHighSimilarity {
			found = true
		}
	}
	if !found {
		t.Fatalf("want high-similarity reason, got %+v", a.Reasons)
	}
}

func TestAssessCoreCollisionDifferentIndustry(t *testing.T) {
	req := Request{
		Name:     "星澜科技",
		Industry: "软件开发",
		ExistingNames: []ExistingName{{
			Name: "星澜餐饮有限公司", Industry: "餐饮服务", Jurisdiction: JurisdictionCNMainland,
		}},
	}
	a := Assess(req.Name, req)
	if !a.Available || a.Risk != RiskLow {
		t.Fatalf("different industry should not collide, got %+v", a)
	}
}

func TestAssessTrademarkConflict(t *testing.T) {
	req := Request{
		Name:     "星澜科技",
		Industry: "软件开发",
		ExistingTrademarks: []ExistingTrademark{{
			Mark: "星澜", Class: "42", Jurisdiction: JurisdictionCNMainland,
		}},
	}
	a := Assess(req.Name, req)
	found := false
	for _, r := range a.Reasons {
		if r.Code == ReasonTrademarkConflict {
			found = true
		}
	}
	if !found {
		t.Fatalf("want trademark-conflict reason, got %+v", a.Reasons)
	}
}

func TestAssessReservedWord(t *testing.T) {
	req := Request{Name: "中国星澜科技", Jurisdiction: JurisdictionCNMainland}
	a := Assess(req.Name, req)
	found := false
	for _, r := range a.Reasons {
		if r.Code == ReasonReservedWord {
			found = true
		}
	}
	if !found {
		t.Fatalf("want reserved-word reason, got %+v", a.Reasons)
	}
}

func TestAssessClean(t *testing.T) {
	req := Request{
		Name: "星澜科技", Industry: "软件开发", Jurisdiction: JurisdictionCNMainland,
		ExistingNames: []ExistingName{{Name: "星辰科技有限公司", Industry: "餐饮"}},
	}
	a := Assess(req.Name, req)
	if !a.Available || a.Risk != RiskLow || len(a.Reasons) != 0 {
		t.Fatalf("want clean assessment, got %+v", a)
	}
}
