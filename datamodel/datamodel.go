package datamodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TeoSocs/alisi-client/crypto"
	"github.com/dgrijalva/jwt-go"
	"github.com/op/go-logging"
	"io/ioutil"
	"os"
	"path"
)

var log = logging.MustGetLogger("alisi")

const CLAIM_FOLDER = "claims"

func (c EncodedClaim) CreateAndStore() (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
		}
	}()
	claimId := c.Id
	claimPath := getPathFor(claimId)
	checkNonExistent(claimPath)
	c.writeInFile(claimPath)

	log.Infof("claim %s stored", path.Base(claimPath))
	return
}

func GetClaim(claimId string) (claim Claim, err error) {

	// WARNING: checkExistent can Panic
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			return
		}
	}()

	claimPath := getPathFor(claimId)
	log.Debugf("retrieving claim from %s", claimPath)
	checkExistent(claimPath)
	data, err := ioutil.ReadFile(claimPath)
	if err != nil {
		log.Error("error reading file %s: %s", path.Base(claimPath), err)
		return
	}
	enClaim := EncodedClaim{}
	err = json.Unmarshal(data, &enClaim)
	if err != nil {
		log.Error("error reading encodedClaim: %s", err)
		return
	}

	clearData, err := crypto.ReadJWT(enClaim.EncodedData)
	publicKey, err := crypto.DecodePublicKeyFromPem(clearData["sgk"].(string))
	if err != nil {
		log.Error("error reading publicKey: %s", err)
		return
	}
	mapClaims, err := crypto.CheckJWTSignature(enClaim.EncodedData, publicKey)
	if err != nil {
		log.Error("error validating JWT: %s", err)
		return
	}
	//parsedTime, err := time.Parse(time.RFC3339, mapClaims["iat"].(string))
	//parsedTime := time.Unix(mapClaims["iat"].(int64), 0)

	//parsedTime, err := strconv.Atoi(mapClaims["iat"].(float64))
	//if err != nil {
	//	log.Error("error parsing timestamp: %s", err)
	//	return
	//}

	claim = Claim{
		Iss:   mapClaims["iss"].(string),
		Iat:   int32(mapClaims["iat"].(float64)),
		Sgk:   mapClaims["sgk"].(string),
		Sub:   mapClaims["sub"].(string),
		Claim: mapClaims["claim"].(string),
	}
	log.Infof("claim %s retrieved", path.Base(claimPath))
	return
}

func GetEncoded(claimId string) (claim EncodedClaim, err error) {

	// WARNING: checkExistent can Panic
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
			return
		}
	}()

	claimPath := getPathFor(claimId)
	log.Debugf("retrieving claim from %s", claimPath)
	checkExistent(claimPath)
	data, err := ioutil.ReadFile(claimPath)
	if err != nil {
		log.Error("error reading file %s: %s", path.Base(claimPath), err)
		return
	}
	err = json.Unmarshal(data, &claim)
	if err != nil {
		log.Error("error decoding JWT: %s", err)
		return
	}

	log.Infof("encodedClaim %s retrieved", path.Base(claimPath))
	return
}

func GetClaimList() (claimList []string, err error) {
	if _, err = os.Stat(CLAIM_FOLDER); os.IsNotExist(err) {
		log.Debugf("folder %s doesn't exists. Creating %s", CLAIM_FOLDER, CLAIM_FOLDER)
		err = os.Mkdir(CLAIM_FOLDER, os.FileMode(os.ModePerm))
	}

	if err != nil {
		log.Error(err)
		return
	}

	fileInfoList, err := ioutil.ReadDir(CLAIM_FOLDER)
	if err != nil {
		return
	}
	claimList = []string{}
	for _, fInfo := range fileInfoList {
		claimList = append(claimList, fInfo.Name())
	}

	return
}

func (c EncodedClaim) Overwrite() (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
			return
		}
	}()
	claimId := c.Id
	claimPath := getPathFor(claimId)
	checkExistent(claimPath)
	c.writeInFile(claimPath)

	log.Infof("claim %s stored", path.Base(claimPath))
	return
}

func DeleteClaim(claimId string) (err error) {

	claimPath := getPathFor(claimId)
	log.Debugf("checking if %s exists", claimId)

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(string))
			return
		}
	}()

	checkExistent(claimPath)

	err = os.Remove(claimPath)
	if err != nil {
		log.Infof("deleted claim %s", claimId)
	}
	return
}

func (c Claim) isEqual(other Claim) bool {
	return c.Iat == other.Iat &&
		c.Iss == other.Iss &&
		c.Sub == other.Sub &&
		c.Sgk == other.Sgk &&
		c.Claim == other.Claim
}

func (c Claim) isEqualExceptTime(other Claim) bool {
	return c.Iss == other.Iss &&
		c.Sub == other.Sub &&
		c.Sgk == other.Sgk &&
		c.Claim == other.Claim
}

func (c EncodedClaim) isEqual(other EncodedClaim) bool {
	return c.EncodedData == other.EncodedData &&
		c.Id == other.Id &&
		c.Signature == other.Signature
}

func (c Claim) encode() (encoded string) {

	mapClaims := jwt.MapClaims{
		"iss":   c.Iss,
		"sgk":   c.Sgk,
		"sub":   c.Sub,
		"iat":   c.Iat,
		"claim": c.Claim,
	}

	encoded, err := crypto.SignJwt(mapClaims)

	if err != nil {
		log.Panicf("error encoding the claim: %s", err)
	}

	log.Info("claim encoded")
	return
}

func getPathFor(claimId string) (claimPath string) {

	log.Debug("checking the path")
	check, err := path.Match("*", claimId)
	if err != nil {
		log.Panicf(err.Error())
	}
	if !check {
		log.Panicf("invalid claimId. Use a literal name without '/' characters instead")
	}
	claimPath = path.Join(CLAIM_FOLDER, claimId)
	return
}

func checkNonExistent(claimPath string) {

	if _, err := os.Stat(claimPath); err == nil {
		log.Panicf("the claim %s already exists", path.Base(claimPath))
	}
}

func checkExistent(claimPath string) {

	if _, err := os.Stat(claimPath); os.IsNotExist(err) {
		log.Panicf("the claim %s doesn't exists", path.Base(claimPath))
	}

}

func (c Claim) writeInFile(claimPath string) {
	// TODO check if better using ioutils
	if _, err := os.Stat(CLAIM_FOLDER); os.IsNotExist(err) {
		log.Debugf("folder %s doesn't exists. Creating %s", CLAIM_FOLDER, CLAIM_FOLDER)
		if err = os.Mkdir(CLAIM_FOLDER, os.FileMode(os.ModePerm)); err != nil {
			log.Panicf(err.Error())
		}
	}

	log.Info("creating the new file")
	file, err := os.Create(claimPath)
	if err != nil {
		log.Panicf("error creating file %s: %s", path.Base(claimPath), err)
	}
	defer file.Close()

	encoded := c.encode()

	_, err = fmt.Fprintf(file, encoded)
	if err != nil {
		log.Panicf("error writing data into file %s: %s", path.Base(claimPath), err)
	}

	return
}

func (c EncodedClaim) writeInFile(claimPath string) {
	// TODO check if better using ioutils
	if _, err := os.Stat(CLAIM_FOLDER); os.IsNotExist(err) {
		log.Debugf("folder %s doesn't exists. Creating %s", CLAIM_FOLDER, CLAIM_FOLDER)
		if err = os.Mkdir(CLAIM_FOLDER, os.FileMode(os.ModePerm)); err != nil {
			log.Panicf(err.Error())
		}
	}

	log.Info("creating the new file")
	file, err := os.Create(claimPath)
	if err != nil {
		log.Panicf("error creating file %s: %s", path.Base(claimPath), err)
	}
	defer file.Close()

	encoded, err := json.Marshal(c)

	_, err = fmt.Fprintf(file, string(encoded))
	if err != nil {
		log.Panicf("error writing data into file %s: %s", path.Base(claimPath), err)
	}

	return
}
