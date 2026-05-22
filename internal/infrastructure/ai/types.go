package ai

type Plan struct {
    Summary          string        `json:"summary"`
    RiskLevel        string        `json:"risk_level"`
    RequiresApproval bool          `json:"requires_approval"`
    Commands         []PlanCommand `json:"commands"`
    UserExplanation  string       `json:"user_explanation"`
}

type PlanCommand struct {
    Command        string `json:"command"`
    Purpose        string `json:"purpose"`
    ExpectedEffect string `json:"expected_effect"`
    RiskLevel      string `json:"risk_level"`
}
