package namecheck

const disclaimer = "本报告基于调用方提供的公开检索事实做初步判断，不构成法律意见，也不保证任何名称一定可注册。正式注册前请以登记机关与商标局的审查结果为准。"

// Check runs the full pipeline: assess the original name, collect the reasons
// it is unavailable, and generate the requested number of alternatives.
func Check(req Request) Report {
	if req.Jurisdiction == "" {
		req.Jurisdiction = JurisdictionCNMainland
	}
	original := Assess(req.Name, req)
	reasons := make([]Reason, 0, len(original.Reasons))
	reasons = append(reasons, original.Reasons...)
	return Report{
		Input:              req,
		Original:           original,
		UnavailableReasons: reasons,
		Candidates:         GenerateCandidates(req),
		Disclaimer:         disclaimer,
	}
}
