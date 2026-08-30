package check

type CheckRequest struct {
	Identifier string `json:"identifier" binding:"required"`
}

type CheckResult struct {
	Allowed   bool  `json:"allowed"`
	Remaining int   `json:"remaining"`
	ResetAt   int64 `json:"reset_at"`
}
