package domain

type SetPinRequest struct {
	Pin string `json:"pin"`
}

type VerifyPinRequest struct {
	Pin string `json:"pin"`
}
