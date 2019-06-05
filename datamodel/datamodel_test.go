package datamodel

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var testClaim = Claim{
	Iss:   "manufacturer_user",
	Sgk:   "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEyZcpRkSzDwnlRhUEi/VXRXqvd+Sx\nNVb0hfB3k7OEE/aW8h2kODosHIEXznAp0Qtebeda7YWFtJepBj2udhBSBw==\n-----END PUBLIC KEY-----\n",
	Sub:   "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEG90CSm32RfW8KsK8sOo2Y/PhNzIf\n6rpd3EzLXUbbjJGCzCAS0yMIBbxvvoS8zTU4PlFLzwXJuiEufQ0T1h/zAw==\n-----END PUBLIC KEY-----\n",
	Claim: "{\"certified_device\":\"true\"}",
	Iat:   1557905444,
}

const testClaimId = ".testclaim"

//var testEncodedClaim = EncodedClaim{
//	PublicKey: "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE4QWksDXnawpXJlRz4zadDSB1eJeH\nrTBNWwryZp02b+HL90g3XIcOcWv/7abb55Lj4tpB3dWIq7MdkueDCJpKbA==\n-----END PUBLIC KEY-----",
//	Id:          testClaimId,
//	EncodedData: "eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzUxMiJ9.eyJpc3MiOiJtYW51ZmFjdHVyZXJfdXNlciIsInNnayI6Ii0tLS0tQkVHSU4gUFVCTElDIEtFWS0tLS0tXG5NRmt3RXdZSEtvWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUU0UVdrc0RYbmF3cFhKbFJ6NHphZERTQjFlSmVIXG5yVEJOV3dyeVpwMDJiK0hMOTBnM1hJY09jV3YvN2FiYjU1TGo0dHBCM2RXSXE3TWRrdWVEQ0pwS2JBPT1cbi0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLVxuIiwic3ViIjoiLS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS1cbk1Ga3dFd1lIS29aSXpqMENBUVlJS29aSXpqMERBUWNEUWdBRUc5MENTbTMyUmZXOEtzSzhzT28yWS9QaE56SWZcbjZycGQzRXpMWFViYmpKR0N6Q0FTMHlNSUJieHZ2b1M4elRVNFBsRkx6d1hKdWlFdWZRMFQxaC96QXc9PVxuLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tXG4iLCJpYXQiOjE1NTc5MDU0NDQsImNsYWltIjoie1wiY2VydGlmaWVkX2RldmljZVwiOlwidHJ1ZVwifSJ9.9Ia6PNqkFo1-dUkynOv6vl0CRcz3rjrfi17WumkldAy6Ml-0ikTpiuGYCjxxjKKNzyF5W04cw2MZf5dYH0RHdw",
//	Signature:   "",
//}

var testClaimPath = path.Join(CLAIM_FOLDER, testClaimId)

func testEncodedClaim() EncodedClaim {
	var encoded = EncodedClaim{}
	_ = json.Unmarshal([]byte(`{"id":".testclaim","encodedData":"eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9.eyJpc3MiOiJtYW51ZmFjdHVyZXJfdXNlciIsInNnayI6Ii0tLS0tQkVHSU4gUFVCTElDIEtFWS0tLS0tXG5NRmt3RXdZSEtvWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUV5WmNwUmtTekR3bmxSaFVFaS9WWFJYcXZkK1N4XG5OVmIwaGZCM2s3T0VFL2FXOGgya09Eb3NISUVYem5BcDBRdGViZWRhN1lXRnRKZXBCajJ1ZGhCU0J3PT1cbi0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLVxuIiwic3ViIjoiLS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS1cbk1Ga3dFd1lIS29aSXpqMENBUVlJS29aSXpqMERBUWNEUWdBRUc5MENTbTMyUmZXOEtzSzhzT28yWS9QaE56SWZcbjZycGQzRXpMWFViYmpKR0N6Q0FTMHlNSUJieHZ2b1M4elRVNFBsRkx6d1hKdWlFdWZRMFQxaC96QXc9PVxuLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tXG4iLCJpYXQiOjE1NTc5MDk2NzEsImNsYWltIjoie1wiY2VydGlmaWVkX2RldmljZVwiOlwidHJ1ZVwifSJ9.PHD7hBeU-ae4PLMhWWZ9Ud_KlZ5s_inM9g5_ih7_2eeRNFjnBNuFZ6D_tnwbc5ploDs3TvqAZZaIcX-aYc8QwA"}`), &encoded)
	return encoded
}

func createTestEncodedClaim() {
	if err := os.Remove(testClaimPath); err != nil {
		log.Debugf("there is no %s to clean", testClaimId)
	} else {
		log.Debugf("%s cleaned up", testClaimId)
	}
	if err := testEncodedClaim().CreateAndStore(); err != nil {
		log.Panic(err)
	}
	log.Infof("new %s created", testClaimId)
}

func cleanTestClaim() {
	log.Debug("cleaning up .testclaim")
	if err := os.Remove(testClaimPath); err != nil {
		log.Panicf("error during cleanup: %s", err)
	}
}

func TestEncodedClaimCreation(t *testing.T) {
	createTestEncodedClaim()
	defer cleanTestClaim()

	data, err := ioutil.ReadFile(testClaimPath)
	if err != nil {
		t.Fatalf("error reading file %s: %s", testClaimId, err)
	}

	read, err := json.Marshal(testEncodedClaim())
	if string(data) != string(read) {
		t.Fatalf("error saving encoded claim:\n%s expected\n%s stored", string(data), read)
	}
}

func TestClaimRead(t *testing.T) {
	createTestEncodedClaim()
	defer cleanTestClaim()
	claimStored, err := GetClaim(testClaimId)
	if err != nil {
		t.Fatal(err)
	}
	if !claimStored.isEqualExceptTime(testClaim) {
		t.Fatalf("error deserializing claim:\n%v expected\n%v read", testClaim, claimStored)
	}
}

func TestGetEncoded(t *testing.T) {
	createTestEncodedClaim()
	defer cleanTestClaim()
	claimStored, err := GetEncoded(testClaimId)
	if err != nil {
		t.Fatal(err)
	}
	if !claimStored.isEqual(testEncodedClaim()) {
		t.Fatalf("error deserializing claim:\n%v expected\n%v read", testClaim, claimStored)
	}
}

func TestEncodedClaimOverwriteFail(t *testing.T) {
	createTestEncodedClaim()
	defer cleanTestClaim()

	if err := testEncodedClaim().CreateAndStore(); err == nil {
		t.Fatal("no error raised by claim overwriting")
	}
}

func TestEncodedClaimUpdate(t *testing.T) {
	createTestEncodedClaim()
	defer cleanTestClaim()

	testClaim2 := EncodedClaim{
		Id:          testClaimId,
		Signature:   "secondNonceSigned",
		EncodedData: "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGFpbSI6InRlc3QtY2xhaW0iLCJpYXQiOiIyMDE5LTAzLTAxVDE1OjQwOjQ1LjE0MDA2ODg3MyswMTowMCIsImlzcyI6InRlc3QtaXNzdWVyIiwic2drIjoidGVzdC1zZ2siLCJzdWIiOiJ0ZXN0LXN1YmplY3QifQ.BZHH5qnI8N6ecvhLMh10nLEJZdQFwzxM6A5naI0i8mqtvlrAVcjT7uo8-LUWyNTDPlyljVRwNBM1IEBia9vP9w",
	}
	if err := testClaim2.Overwrite(); err != nil {
		t.Fatal(err)
		return
	}
	data, err := ioutil.ReadFile(testClaimPath)
	if err != nil {
		t.Fatalf("error reading file %s: %s", testClaimId, err)
		return
	}

	read, err := json.Marshal(testClaim2)
	if string(data) != string(read) {
		t.Fatalf("error overwriting encoded claim:\n%s expected\n%s stored", data, read)
	}
}

func TestUpdateNonexistentClaim(t *testing.T) {
	_ = os.Remove(testClaimPath)
	if err := testEncodedClaim().Overwrite(); err == nil {
		t.Fatalf("no error raised by overwriting a nonexistent claim")
	}
	_ = os.Remove(testClaimPath)

}

func TestUpdateNonexistentEncodedClaim(t *testing.T) {
	_ = os.Remove(testClaimPath)
	if err := testEncodedClaim().Overwrite(); err == nil {
		t.Fatalf("no error raised by overwriting a nonexistent claim")
	}
	_ = os.Remove(testClaimPath)
}

func TestClaimDelete(t *testing.T) {
	createTestEncodedClaim()
	err := DeleteClaim(testClaimId)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(testClaimPath); err == nil {
		t.Fatal(".testclaim still exists")
	}
}

func TestDeleteNonexistentClaim(t *testing.T) {
	_ = os.Remove(testClaimPath)
	if err := testEncodedClaim().Overwrite(); err == nil {
		t.Fatalf("no error raised by deleting a nonexistent claim")
	}
}

func TestClaimListRead(t *testing.T) {

	_ = os.Remove(testClaimPath)
	claimList, err := GetClaimList()
	if err != nil {
		t.Fatal(err)
	}
	for _, claim := range claimList {
		if claim == testClaimId {
			t.Fatalf("%s listed after being explicitly removed", testClaimId)
		}
	}
	createTestEncodedClaim()
	defer cleanTestClaim()

	claimList, err = GetClaimList()
	if err != nil {
		t.Fatal(err)
	}
	ok := false
	for _, claim := range claimList {
		if claim == testClaimId {
			ok = true
		}
	}

	if !ok {
		t.Fatalf("%s not listed after being explicitly created", testClaimId)
	}
}
