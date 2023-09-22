package main

type FpPayload struct {
	N  string `json:"n"`
	ID string `json:"id"`
}

type Fingerprint struct {
	ID          string `json:"id,omitempty"`
	Fingerprint string `json:"fp"`
}
