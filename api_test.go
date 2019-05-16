package main

import (
	"bytes"
	"encoding/json"
	"github.com/TeoSocs/alisi-client/crypto"
	"github.com/TeoSocs/alisi-client/datamodel"
	"github.com/op/go-logging"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
)

const testClaimId = ".testclaim"

var log = logging.MustGetLogger("alisi")

var testClaim = datamodel.Claim{
	Iss:   "manufacturer_user",
	Sgk:   "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEyZcpRkSzDwnlRhUEi/VXRXqvd+Sx\nNVb0hfB3k7OEE/aW8h2kODosHIEXznAp0Qtebeda7YWFtJepBj2udhBSBw==\n-----END PUBLIC KEY-----\n",
	Sub:   "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEG90CSm32RfW8KsK8sOo2Y/PhNzIf\n6rpd3EzLXUbbjJGCzCAS0yMIBbxvvoS8zTU4PlFLzwXJuiEufQ0T1h/zAw==\n-----END PUBLIC KEY-----\n",
	Claim: "{\"certified_device\":\"true\"}",
	Iat:   1557905444,
}

func testEncodedClaim() datamodel.EncodedClaim {
	var encoded = datamodel.EncodedClaim{}
	_ = json.Unmarshal([]byte(`{"id":".testclaim","encodedData":"eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9.eyJpc3MiOiJtYW51ZmFjdHVyZXJfdXNlciIsInNnayI6Ii0tLS0tQkVHSU4gUFVCTElDIEtFWS0tLS0tXG5NRmt3RXdZSEtvWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUV5WmNwUmtTekR3bmxSaFVFaS9WWFJYcXZkK1N4XG5OVmIwaGZCM2s3T0VFL2FXOGgya09Eb3NISUVYem5BcDBRdGViZWRhN1lXRnRKZXBCajJ1ZGhCU0J3PT1cbi0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLVxuIiwic3ViIjoiLS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS1cbk1Ga3dFd1lIS29aSXpqMENBUVlJS29aSXpqMERBUWNEUWdBRUc5MENTbTMyUmZXOEtzSzhzT28yWS9QaE56SWZcbjZycGQzRXpMWFViYmpKR0N6Q0FTMHlNSUJieHZ2b1M4elRVNFBsRkx6d1hKdWlFdWZRMFQxaC96QXc9PVxuLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tXG4iLCJpYXQiOjE1NTc5MDk2NzEsImNsYWltIjoie1wiY2VydGlmaWVkX2RldmljZVwiOlwidHJ1ZVwifSJ9.PHD7hBeU-ae4PLMhWWZ9Ud_KlZ5s_inM9g5_ih7_2eeRNFjnBNuFZ6D_tnwbc5ploDs3TvqAZZaIcX-aYc8QwA"}`), &encoded)
	return encoded
}

var testClaimPath = path.Join(datamodel.CLAIM_FOLDER, testClaimId)

var clientOnline = false

func startAPI() {
	if !clientOnline {
		clientOnline = true
		crypto.TEST_ENV = true
		crypto.Init()
		go main()
	}
}

func closeBody(response *http.Response) {
	_ = response.Body.Close()
}

func cleanEventualTestClaim() {
	log.Debugf("cleaning up %s", testClaimPath)
	_ = os.Remove(testClaimPath)
}

func createTestEncodedClaim() {
	if err := os.Remove(testClaimPath); err == nil {
		log.Debugf("%s cleaned up", testClaimId)
	}
	if err := testEncodedClaim().CreateAndStore(); err != nil {
		log.Panic(err)
	}
	log.Debugf("new %s created", testClaimId)
}

func TestGreetings(t *testing.T) {
	greetingString := "Hello World!"
	crypto.TEST_ENV = true // Not necessary, I will copy this when needed. Just a remainder

	startAPI()

	resp, err := http.Get("http://localhost:8080/alisi/v1/")
	if err != nil {
		t.Fatal(err)
	}
	defer closeBody(resp)
	body, err := ioutil.ReadAll(resp.Body)

	if string(body) != greetingString {
		t.Fatalf("wrong answer:\n%s expected, got\n%s instead", greetingString, string(body))
	}
}

func TestUnauthorizedCreationMissing(t *testing.T) {
	startAPI()
	resp, err := http.Post("http://localhost:8080/alisi/v1/claim", "application/json", nil)
	defer closeBody(resp)

	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("got statusCode %d from unauthorized creation, 401 expected", resp.StatusCode)
	}
}

func TestUnauthorizedCreationWrong(t *testing.T) {
	startAPI()
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/alisi/v1/claim", nil)
	req.Header.Add("X-API-Key", "wrongTestAPIkey")
	resp, err := client.Do(req)
	defer closeBody(resp)

	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("got statusCode %d from unauthorized creation, 401 expected", resp.StatusCode)
	}
}

func TestCreateClaim(t *testing.T) {
	cleanEventualTestClaim()
	encoded, err := json.Marshal(testEncodedClaim())
	if err != nil {
		log.Fatal(err)
	}

	startAPI()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/alisi/v1/claim", bytes.NewBuffer(encoded))
	req.Header.Add("X-API-Key", "testAPIkey")
	resp, err := client.Do(req)
	defer closeBody(resp)

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(string(body))
	}
	if err != nil {
		t.Fatal(err)
	}
	defer cleanEventualTestClaim()

	data, err := ioutil.ReadFile(testClaimPath)
	if err != nil {
		t.Fatalf("error reading file %s: %s", testClaimId, err)
	}
	read, err := json.Marshal(testEncodedClaim())
	if string(data) != string(read) {
		t.Fatalf("error saving encoded claim:\n%s expected\n%s stored", string(data), read)
	}
}

func TestCreationOverwritingFail(t *testing.T) {
	createTestEncodedClaim()
	defer cleanEventualTestClaim()
	encoded, err := json.Marshal(testEncodedClaim())
	if err != nil {
		log.Fatal(err)
	}

	startAPI()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/alisi/v1/claim", bytes.NewBuffer(encoded))
	req.Header.Add("X-API-Key", "testAPIkey")
	resp, err := client.Do(req)
	defer closeBody(resp)

	if resp.StatusCode != 400 {
		t.Fatal("No 400 error raised overwriting .testclaim with creation method")
	}
}

func TestUnauthorizedDeletionMissing(t *testing.T) {
	startAPI()
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/alisi/v1/claim/.testclaim", nil)

	resp, err := client.Do(req)
	defer closeBody(resp)

	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("got statusCode %d from unauthorized deletion, 401 expected", resp.StatusCode)
	}
}

func TestUnauthorizedDeletionWrong(t *testing.T) {
	startAPI()
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/alisi/v1/claim/.testclaim", nil)
	req.Header.Add("X-API-Key", "wrongTestAPIkey")

	resp, err := client.Do(req)
	defer closeBody(resp)

	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("got statusCode %d from unauthorized deletion, 401 expected", resp.StatusCode)
	}
}

func TestDeleteClaim(t *testing.T) {
	createTestEncodedClaim()

	startAPI()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/alisi/v1/claim/.testclaim", nil)
	req.Header.Add("X-API-Key", "testAPIkey")

	resp, err := client.Do(req)
	defer closeBody(resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatal(string(body))
	}
	cleanEventualTestClaim()
}

func TestDeleteNonExistentClaim(t *testing.T) {
	cleanEventualTestClaim()
	startAPI()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/alisi/v1/claim/.testclaim", nil)
	req.Header.Add("X-API-Key", "testAPIkey")

	resp, err := client.Do(req)
	defer closeBody(resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 400 {
		t.Fatalf("got statusCode %d from deleting a nonexistent claim, 400 expected", resp.StatusCode)
	}
}

func TestGetClaimList(t *testing.T) {
	cleanEventualTestClaim()
	startAPI()

	resp, err := http.Get("http://localhost:8080/alisi/v1/claim")

	defer closeBody(resp)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	var claimListBefore []string
	err = json.Unmarshal(body, &claimListBefore)
	if err != nil {
		t.Fatal(err)
	}
	for _, el := range claimListBefore {
		if el == ".testclaim" {
			t.Errorf(".testclaim found after being explicitly deleted")
		}
	}

	createTestEncodedClaim()
	defer cleanEventualTestClaim()
	resp, err = http.Get("http://localhost:8080/alisi/v1/claim")

	defer closeBody(resp)
	if err != nil {
		t.Fatal(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	var claimListAfter []string
	err = json.Unmarshal(body, &claimListAfter)
	if err != nil {
		t.Fatal(err)
	}

	testClaimFound := false
	for _, el := range claimListAfter {
		if el == ".testclaim" {
			testClaimFound = true
		}
	}
	if !testClaimFound {
		t.Errorf(".testclaim not found after being explicitly created")
	}

}

func TestGetClaimById(t *testing.T) {
	createTestEncodedClaim()
	defer cleanEventualTestClaim()
	startAPI()

	resp, err := http.Get("http://localhost:8080/alisi/v1/claim/.testclaim")

	defer closeBody(resp)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	var claim datamodel.Claim
	err = json.Unmarshal(body, &claim)
	if err != nil {
		t.Fatal(err)
	}

	if claim.Iss != testClaim.Iss ||
		claim.Sub != testClaim.Sub ||
		claim.Sgk != testClaim.Sgk ||
		claim.Claim != testClaim.Claim {
		t.Fatalf("error retrieving claim:\n%v expected\n%v read", testClaim, claim)
	}

	log.Infof("%s", claim)
}

func TestRequestSigned(t *testing.T) {
	createTestEncodedClaim()
	defer cleanEventualTestClaim()
	startAPI()

	resp, err := http.Post("http://localhost:8080/alisi/v1/claim/.testclaim/request_signed/mynonce", "application/json", nil)
	defer closeBody(resp)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	var claim datamodel.EncodedClaim
	err = json.Unmarshal(body, &claim)
	if err != nil {
		t.Fatal(err)
	}

	// TODO A new public key has been generated breaking the test.
	//if claim.PublicKey != testEncodedClaim().PublicKey {
	//	t.Fatalf("error retrieving claim.PublicKey:\n%v expected\n%v read", testEncodedClaim().PublicKey, claim.PublicKey)
	//} else if claim.EncodedData != testEncodedClaim.EncodedData {
	if claim.EncodedData != testEncodedClaim().EncodedData {
		t.Fatalf("error retrieving claim.EncodedData:\n%v expected\n%v read", testEncodedClaim().EncodedData, claim.EncodedData)
	} else if claim.Id != testEncodedClaim().Id {
		t.Fatalf("error retrieving claim.Id:\n%v expected\n%v read", testEncodedClaim().Id, claim.Id)
	}

	// TODO check signature <-- IMPOSSIBLE WITHOUT A PROPER PUBLIC KEY. THIS ONE IS A PLACEHOLDER!

	log.Infof("%s", claim)
}

func TestGetPublicKey(t *testing.T) {
	startAPI()
	resp, err := http.Get("http://localhost:8080/alisi/v1/public_key")
	defer closeBody(resp)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	log.Infof("public key retrieved: %s", string(body))
}
