package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"math/big"
	"os"
)

var TEST_ENV = false

var keyName = "private.pem"

func newPrivateKey() *ecdsa.PrivateKey {
	log.Println("creating a new private key")
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Println(err)
	}
	log.Println("new private key created")
	return privateKey
}

func encodePrivateKeyToPem(key *ecdsa.PrivateKey) string {
	x509Encoded, _ := x509.MarshalECPrivateKey(key)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	log.Println("private key encoded to PEM")
	return string(pemEncoded)
}

func decodePrivateKeyFromPem(encoded string) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(encoded))
	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		log.Print(err)
	} else {
		log.Println("private key decoded from PEM")
	}
	return privateKey
}

func EncodePublicKeyToPem(key *ecdsa.PublicKey) string {
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(key)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	log.Println("public key encoded to PEM")
	return string(pemEncodedPub)
}

func DecodePublicKeyFromPem(encoded string) (publicKey *ecdsa.PublicKey, err error) {
	blockPub, _ := pem.Decode([]byte(encoded))

	if blockPub == nil {
		log.Printf("Invalid publicKey: %v", encoded)
	}

	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		log.Printf("error parsing x509: %s", err)
		return
	}
	publicKey = genericPublicKey.(*ecdsa.PublicKey)
	log.Println("public key decoded from PEM")
	return
}

func Sign(message string) (r *big.Int, s *big.Int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
			return
		}
	}()
	key, err := getPrivateKey()
	r, s = sign(message, key)
	return
}

func sign(message string, key *ecdsa.PrivateKey) (r *big.Int, s *big.Int) {
	byteMessage := []byte(message)
	r, s, err := ecdsa.Sign(rand.Reader, key, byteMessage)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("message signed")
	return
}

func verify(key *ecdsa.PublicKey, message string, r *big.Int, s *big.Int) bool {
	byteMessage := []byte(message)
	check := ecdsa.Verify(key, byteMessage, r, s)
	if check {
		log.Println("signature verified")
	} else {
		log.Println("signature refused")
	}
	return check
}

func EncodeSignatureDER(r *big.Int, s *big.Int) (der []byte, err error) {
	sig := ECDSASignature{R: r, S: s}
	der, err = asn1.Marshal(sig)
	return
}

func DecodeSignatureDER(der []byte) (r *big.Int, s *big.Int, err error) {
	sig := &ECDSASignature{}
	_, err = asn1.Unmarshal(der, sig)
	if err != nil {
		return
	}
	r = sig.R
	s = sig.S
	return
}

func getPrivateKey() (privateKey *ecdsa.PrivateKey, err error) {
	//var user string
	//if TEST_ENV {
	//	user = "test"
	//} else {
	//	user = "crypto"
	//}
	//// get signingKey
	//secret, err := keyring.Get(service, user)

	if TEST_ENV == true {
		keyName = "test.pem"
	}
	keyPath := "keys/" + keyName
	secret, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("key retrieved from %s", keyPath)
	privateKey = decodePrivateKeyFromPem(string(secret))
	return
}

func GetPublicKey() (publicKey *ecdsa.PublicKey, err error) {
	privKey, err := getPrivateKey()
	if err != nil {
		log.Println(err)
		return
	}
	publicKey = &privKey.PublicKey
	log.Printf("got public key")
	return
}

func storePrivateKey(key *ecdsa.PrivateKey) {
	if TEST_ENV == true {
		keyName = "test.pem"
	}
	keyPath := "keys/" + keyName
	if _, err := os.Stat("keys"); os.IsNotExist(err) {
		log.Print("folder keys doesn't exists. Creating keys")
		err = os.Mkdir("keys", os.FileMode(os.ModePerm))
	}
	data := []byte(encodePrivateKeyToPem(key))
	err := ioutil.WriteFile(keyPath, data, os.FileMode(os.ModePerm))
	//var user string
	//if TEST_ENV {
	//	user = "test"
	//} else {
	//	user = "crypto"
	//}
	//privatePem := encodePrivateKeyToPem(key)
	//err := keyring.Set(service, user, privatePem)
	if err != nil {
		log.Println(err)
	}
	log.Printf("key stored in %s", keyPath)
}

func SignJwt(claims jwt.MapClaims) (encoded string, err error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	privateKey, err := getPrivateKey()

	// Sign and get the complete encoded token as a string using the secret
	encoded, err = token.SignedString(privateKey)

	return
}

func ReadJWT(tokenString string) (claims jwt.MapClaims, err error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, _ := jwt.Parse(tokenString, nil)

	if token == nil {
		log.Printf("error reading jwt")
		err = errors.New("error reading jwt")
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		log.Print("error extracting claims from jwt")
		err = errors.New("error extracting claims from jwt")
		return
	}
	log.Printf("read JWT")
	return
}

func CheckJWTSignature(tokenString string, key *ecdsa.PublicKey) (claims jwt.MapClaims, err error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return key, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		log.Print("JWT validated")
	} else {
		log.Print(err)
	}
	return
}

func Init() {
	key, err := getPrivateKey()
	if err != nil {
		log.Println("key not found, generating a new one")
		privateKey := newPrivateKey()
		storePrivateKey(privateKey)
	} else {
		log.Println("key found")
		log.Println(key.D)
	}
}
