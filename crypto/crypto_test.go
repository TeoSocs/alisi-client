package crypto

import (
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"testing"
)

func TestNewPrivateKey(t *testing.T) {
	log.Println("creating 2 keys and checking they are different")
	privKey := newPrivateKey()
	privKey1 := newPrivateKey()
	if privKey.D.Cmp(privKey1.D) == 0 {
		t.Error("got the same key twice")
	}
}

func TestPrivPem(t *testing.T) {
	log.Println("creating a key and checking for unwanted mutation during PEM conversion")
	privateKey := newPrivateKey()
	privatePem := encodePrivateKeyToPem(privateKey)
	privKeyFromPem := decodePrivateKeyFromPem(privatePem)
	if privKeyFromPem.D.Cmp(privateKey.D) != 0 {
		t.Errorf("something changed during PEM conversion of the private Key, got \n%d \n from \n%d",
			privKeyFromPem.D, privateKey.D)
	}
}

func TestPubPem(t *testing.T) {
	log.Println("creating a key, extracting the public key and checking for unwanted mutation during PEM conversion")
	privateKey := newPrivateKey()
	publicKey := &privateKey.PublicKey
	publicPem := EncodePublicKeyToPem(publicKey)
	pubKeyFromPem, err := DecodePublicKeyFromPem(publicPem)
	if err != nil {
		t.Fatal(err)
	}
	if pubKeyFromPem.X.Cmp(publicKey.X) != 0 || pubKeyFromPem.Y.Cmp(publicKey.Y) != 0 {
		t.Errorf("something changed during PEM conversion of the public Key, got \n%d, %d \n from \n%d, %d",
			pubKeyFromPem.X, pubKeyFromPem.Y, publicKey.X, publicKey.Y)
	}
}

func TestSignature(t *testing.T) {
	log.Println("creating a key, extracting the public key, signing and verifying a sample message")
	privateKey := newPrivateKey()
	publicKey := &privateKey.PublicKey
	message := "Hello, world!"
	r, s := sign(message, privateKey)
	check := verify(publicKey, message, r, s)
	if !check {
		t.Error("error validating self-signed message")
	}
}

func TestDerEncoding(t *testing.T) {
	log.Println("same of TestSignature, but it checks the der encoding too")
	privateKey := newPrivateKey()
	publicKey := &privateKey.PublicKey
	message := "Hello, world!"
	r, s := sign(message, privateKey)
	encoding, err := EncodeSignatureDER(r, s)
	encodedString := base64.StdEncoding.EncodeToString(encoding)
	if err != nil {
		t.Error(err)
	}
	info := `
message:
%s
key:
%s
signature:
%s
`
	log.Printf(info, message, EncodePublicKeyToPem(publicKey), encodedString)
	rFromDer, sFromDer, err := DecodeSignatureDER(encoding)
	if err != nil {
		t.Error(err)
	}
	if rFromDer.Cmp(r) != 0 || sFromDer.Cmp(s) != 0 {
		t.Errorf("something changed during the der encoding of the signature, got \n%d, %d \n from \n%d, %d",
			rFromDer, sFromDer, r, s)
	}
	check := verify(publicKey, message, r, s)
	if !check {
		t.Error("error validating self-signed message")
	}
}

func TestSecureStorage(t *testing.T) {
	TEST_ENV = true
	privateKey := newPrivateKey()
	storePrivateKey(privateKey)
	readPrivateKey, err := getPrivateKey()
	if err != nil {
		t.Error("error reading the private key just stored")
	}
	if readPrivateKey.D.Cmp(privateKey.D) != 0 {
		t.Errorf("something changed during the archiviation of the private Key, got \n%d \n from \n%d",
			readPrivateKey.D, privateKey.D)
	}

	readPublicKey, err := GetPublicKey()
	if err != nil {
		t.Error("error reading the public key just stored")
	}
	if readPublicKey.X.Cmp(privateKey.PublicKey.X) != 0 || readPublicKey.Y.Cmp(privateKey.PublicKey.Y) != 0 {
		t.Errorf("something changed during the archiviation of the private Key, got \n%d, %d \n from \n%d, %d",
			readPublicKey.X, readPublicKey.Y, privateKey.PublicKey.X, privateKey.PublicKey.Y)
	}
	_ = os.Remove("keys/test.pem")
}

func TestInit(t *testing.T) {
	TEST_ENV = true

	//_ = keyring.Delete(service, "test")
	_ = os.Remove("keys/test.pem")
	_, err := getPrivateKey()
	if err == nil {
		t.Error("key found right after being deleted")
	}
	Init()
	key, err := getPrivateKey()
	if err != nil {
		t.Error("key not found after Init()")
	}
	Init()
	key2, err := getPrivateKey()
	if err != nil {
		t.Error("key not found after second Init()")
	}
	if key.D.Cmp(key2.D) != 0 {
		t.Errorf("different keys retrieved after Init(): \n%d, \n%d", key.D, key2.D)
	}
	_ = os.Remove("keys/test.pem")

}

func TestSignJwt(t *testing.T) {
	TEST_ENV = true
	privateKey := decodePrivateKeyFromPem("-----BEGIN PRIVATE KEY-----\n" +
		"MHcCAQEEIAjxIo9hZ/5NtpEBApv60LTjnzOIe8I3SDVz+vEG8jsloAoGCCqGSM49\n" +
		"AwEHoUQDQgAEG90CSm32RfW8KsK8sOo2Y/PhNzIf6rpd3EzLXUbbjJGCzCAS0yMI\n" +
		"BbxvvoS8zTU4PlFLzwXJuiEufQ0T1h/zAw==\n" +
		"-----END PRIVATE KEY-----")

	storePrivateKey(privateKey)

	claims := jwt.MapClaims{
		"iss":   "c.Iss",
		"sgk":   "c.Sgk",
		"sub":   "c.Sub",
		"iat":   "c.Iat",
		"claim": "c.Claim",
	}

	encoded, err := SignJwt(claims)
	if err != nil {
		t.Error(err)
	}
	//publicKey, err := DecodePublicKeyFromPem("-----BEGIN PUBLIC KEY-----\n" +
	//	"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEUpHfQ4+jRec/C4u5opDJHifushVA" +
	//	"2yN1L0fat6Rk1+V4dfX+R/LJVKYFiSXHjTj5lvVF97cLGXm3hrGozGVeSg==" +
	//	"\n-----END PUBLIC KEY-----")
	//if err != nil {
	//	t.Fatal(err)
	//}
	validatedClaims, err := CheckJWTSignature(encoded, &privateKey.PublicKey)
	if err != nil {
		t.Fatalf("signature invalid: %s", err)
	}
	if validatedClaims["iss"] != "c.Iss" {
		t.Fatalf("wrong claim attribute: iss")
	}
	if validatedClaims["sgk"] != "c.Sgk" {
		t.Fatalf("wrong claim attribute: sgk")
	}
	if validatedClaims["sub"] != "c.Sub" {
		t.Fatalf("wrong claim attribute: sub")
	}
	if validatedClaims["iat"] != "c.Iat" {
		t.Fatalf("wrong claim attribute: iat")
	}
	if validatedClaims["claim"] != "c.Claim" {
		t.Fatalf("wrong claim attribute: claim")
	}
	_ = os.Remove("keys/test.pem")
}
