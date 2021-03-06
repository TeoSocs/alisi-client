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

type EncodedClaim struct {
	Id string `json:"id,omitempty"`

	// JWT-encoded claim
	EncodedData string `json:"encodedData,omitempty"`

	// der encoding of a typical ecdsa signature
	Signature string `json:"signature,omitempty"`
}
