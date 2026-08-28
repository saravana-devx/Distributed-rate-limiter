package client

import "errors"

// Sentinel errors for the client domain, checked with errors.Is at the
// handler boundary to decide which httpx status/message to send.
var (
	ErrClientAlreadyExists = errors.New(MsgClientAlreadyExists)
)
