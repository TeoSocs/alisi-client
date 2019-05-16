package crypto

import "math/big"

// "r" and "s" parameter of a typical ecdsa signature

type ECDSASignature struct {
	R *big.Int `json:"r,omitempty"`

	S *big.Int `json:"s,omitempty"`
}
