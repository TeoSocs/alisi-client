/*
 * ALISI client
 *
 * This is the client API of ALISI. Each device will expose this API in order to be identified by ALISI compliant control units.
 *
 * API version: 1.0.0
 * Contact: matteo.sovilla@studenti.unipd.it
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package datamodel

type Claim struct {

	// Iroha ID of the issuer
	Iss string `json:"iss,omitempty"`

	// PublicKey to use for signature verification
	Sgk string `json:"sgk,omitempty"`

	// DID of the subject of the DID
	Sub string `json:"sub,omitempty"`

	// Issued AT, unix time
	Iat int32 `json:"iat,omitempty"`

	// JSON content of the claim
	Claim string `json:"claim,omitempty"`
}
