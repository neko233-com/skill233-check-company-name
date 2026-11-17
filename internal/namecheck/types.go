// Package namecheck implements the deterministic core of company-name checking:
// normalization, duplicate/approximate detection, reason diagnosis and candidate
// generation. It is deliberately free of network I/O so it can be unit-tested and
// embedded in an Agent Skill that supplies registry facts as input.
package namecheck

const (
	// DefaultCandidateCount is the number of alternative names produced by default.
	DefaultCandidateCount = 5

	// SimilarityThreshold is the minimum normalized Levenshtein similarity for a
	// name to be reported as "highly similar" rather than merely mentioned.
	SimilarityThreshold = 0.78
)

// Jurisdiction identifies a registration/trademark authority.
type Jurisdiction string

const (
	JurisdictionCNMainland Jurisdiction = "CN-MAINLAND"
	JurisdictionCNHK       Jurisdiction = "CN-HK"
	JurisdictionCNMO       Jurisdiction = "CN-MO"
	JurisdictionCNTW       Jurisdiction = "CN-TW"
	JurisdictionUS         Jurisdiction = "US"
	JurisdictionEU         Jurisdiction = "EU"
	JurisdictionUK         Jurisdiction = "UK"
	JurisdictionGlobal     Jurisdiction = "GLOBAL"
)

// RiskLevel is the coarse conclusion for a single name.
type RiskLevel string

const (
	RiskLow    RiskLevel = "LOW"    // no exact/near duplicate found in supplied facts
	RiskMedium RiskLevel = "MEDIUM" // similar names exist; needs human review
	RiskHigh   RiskLevel = "HIGH"   // exact duplicate or strong collision
)

// ReasonCode is a stable machine-readable identifier for why a name is unavailable.
type ReasonCode string

const (
	ReasonExactDuplicate    ReasonCode = "EXACT_DUPLICATE"
	ReasonHighSimilarity    ReasonCode = "HIGH_SIMILARITY"
	ReasonTrademarkConflict ReasonCode = "TRADEMARK_CONFLICT"
	ReasonReservedWord      ReasonCode = "RESERVED_WORD"
	ReasonNoReason          ReasonCode = ""
)

// ExistingName is a fact supplied by the caller (typically gathered from a
// registry search) representing an already registered company name.
type ExistingName struct {
	Name         string       `json:"name"`
	Jurisdiction Jurisdiction `json:"jurisdiction"`
	Industry     string       `json:"industry,omitempty"`
	Source       string       `json:"source,omitempty"`
}

// ExistingTrademark is a fact representing an existing/constant trademark mark.
type ExistingTrademark struct {
	Mark         string       `json:"mark"`
	Class        string       `json:"class,omitempty"`
	Jurisdiction Jurisdiction `json:"jurisdiction"`
	Source       string       `json:"source,omitempty"`
}

// Request is the input to Check.
type Request struct {
	// Name is the user-supplied candidate to diagnose.
	Name string `json:"name"`
	// Industry is the intended business, e.g. "软件开发". Used for industry-aware
	// collision rules and for generating relevant alternatives.
	Industry string `json:"industry,omitempty"`
	// Region narrows generated names, e.g. "上海". Optional.
	Region string `json:"region,omitempty"`
	// Jurisdiction is the primary authority to check first. Defaults to CN-MAINLAND.
	Jurisdiction Jurisdiction `json:"jurisdiction,omitempty"`
	// CandidateCount overrides DefaultCandidateCount when > 0.
	CandidateCount int `json:"candidateCount,omitempty"`
	// EnglishHint, when provided, is used as the Latin/transliteration seed for
	// generated names. Optional.
	EnglishHint string `json:"englishHint,omitempty"`

	ExistingNames      []ExistingName      `json:"existingNames,omitempty"`
	ExistingTrademarks []ExistingTrademark `json:"existingTrademarks,omitempty"`
}

// Reason explains, in human-readable form, why a name is problematic.
type Reason struct {
	Code         ReasonCode   `json:"code"`
	Message      string       `json:"message"`
	Jurisdiction Jurisdiction `json:"jurisdiction,omitempty"`
	Source       string       `json:"source,omitempty"`
}

// Assessment is the evaluation of one name against the supplied facts.
type Assessment struct {
	Name       string    `json:"name"`
	Risk       RiskLevel `json:"risk"`
	Available  bool      `json:"available"`
	Reasons    []Reason  `json:"reasons"`
	Similarity float64   `json:"similarity"`
}

// Candidate is one generated alternative together with its assessment.
type Candidate struct {
	Name       string     `json:"name"`
	Strategy   string     `json:"strategy"`
	Assessment Assessment `json:"assessment"`
}

// Report is the full output.
type Report struct {
	Input              Request     `json:"input"`
	Original           Assessment  `json:"original"`
	UnavailableReasons []Reason    `json:"unavailableReasons"`
	Candidates         []Candidate `json:"candidates"`
	Disclaimer         string      `json:"disclaimer"`
}
