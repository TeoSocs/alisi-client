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

	// the signing public key
	PublicKey string `json:"publicKey,omitempty"`

	// the nonce provided with the request, signed by the client
	Signature string `json:"signature,omitempty"`
}
