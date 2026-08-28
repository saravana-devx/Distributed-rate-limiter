package check

type algorithmType string

type CheckRequest struct {
	Identifier string `json:"identifier"`
}

type CheckResult struct {
	Allowed   bool `json:"allowed"`
	Remaining int  `json:"remaining"`
}
