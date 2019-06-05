# ALISI client


<a name="overview"></a>
## Overview
This is the client API of ALISI. Each device will expose this API in order to be identified by ALISI compliant control units.


### Version information
*Version* : 1.0.0


### Contact information
*Contact Email* : matteo.sovilla@studenti.unipd.it


### URI scheme
*BasePath* : /alisi/v1  
*Schemes* : HTTP


### Tags

* Claims : CRUD operations on the stored claims


### External Docs
*Description* : ALISI system documentation  
*URL* : http://www.github.com/TeoSocs/alisi-network




<a name="paths"></a>
## Paths

<a name="createclaim"></a>
### Add a claim
```
POST /claim
```


#### Description
Create a claim providing the claimID and the JWT-encoded content


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Body**|**body**  <br>*required*|Claim to create|[EncodedClaim](#encodedclaim)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**201**|created|No Content|
|**401**|API key is missing or invalid  <br>**Headers** :   <br>`WWW_Authenticate` (string)|No Content|


#### Tags

* Claims


#### Security

|Type|Name|
|---|---|
|**apiKey**|**[APIKeyHeader](#apikeyheader)**|


<a name="getclaimlist"></a>
### Return the claim list
```
GET /claim
```


#### Description
Returns the list of claimID stored by the client


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|successful operation|< string > array|


#### Produces

* `application/json`


#### Tags

* Claims


<a name="getclaimbyid"></a>
### Return claim by ID
```
GET /claim/{claimID}
```


#### Description
Given a claimID, it looks for the corresponding claim and returns it if exists


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Path**|**claimID**  <br>*required*|ID of the claim to fetch|string|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|successful operation|[Claim](#claim)|
|**404**|claim ID not found|No Content|


#### Produces

* `application/json`


#### Tags

* Claims


<a name="deleteclaim"></a>
### Delete a claim
```
DELETE /claim/{claimID}
```


#### Description
Deletes the claim with this specific claimID, if it exists


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Path**|**claimID**  <br>*required*|ID of the claim to fetch|string|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|successful operation|No Content|
|**401**|API key is missing or invalid  <br>**Headers** :   <br>`WWW_Authenticate` (string)|No Content|
|**404**|claim ID not found|No Content|


#### Tags

* Claims


#### Security

|Type|Name|
|---|---|
|**apiKey**|**[APIKeyHeader](#apikeyheader)**|


<a name="requestsigned"></a>
### Request a claim signed by the client
```
POST /claim/{claimID}/request_signed/{nonce}
```


#### Description
Create a claim request providing the claimID and the nonce that the client has to sign.


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Path**|**claimID**  <br>*required*|ID of the claim to fetch|string|
|**Path**|**nonce**  <br>*required*|Nonce the client needs to sign. Prevents replay attacks|integer|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||No Content|
|**404**|claim ID not found|No Content|


#### Tags

* Claims


<a name="getpublickey"></a>
### Returns the public key
```
GET /public_key
```


#### Description
Returns the public key of the device. Invoked by the manufacturer endpoint in order to create the corresponding claim


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|key retrieved|string|
|**500**|Internal error on crypto material|No Content|


#### Tags

* Crypto




<a name="definitions"></a>
## Definitions

<a name="claim"></a>
### Claim

|Name|Description|Schema|
|---|---|---|
|**claim**  <br>*required*|JSON content of the claim|string|
|**iat**  <br>*required*|Issued AT, unix time|integer|
|**iss**  <br>*required*|Iroha ID of the issuer|string|
|**sgk**  <br>*required*|PublicKey to use for signature verification|string|
|**sub**  <br>*required*|DID of the subject of the DID|string|


<a name="encodedclaim"></a>
### EncodedClaim

|Name|Description|Schema|
|---|---|---|
|**encodedData**  <br>*required*|JWT-encoded claim|string|
|**id**  <br>*required*||string|
|**signature**  <br>*required*|der encoding of a typical ecdsa signature|string|




<a name="securityscheme"></a>
## Security

<a name="apikeyheader"></a>
### APIKeyHeader
*Type* : apiKey  
*Name* : X-API-Key  
*In* : HEADER

<a name="knownissues"></a>
### Known issues
Actually, the private key is stored in the folder `keys`.

This is a security issues, the key should be more protected.
For example, it could be integrated with the gnome keyring,
or a dedicated hardware.

This is not already done because this system is still running
in a simulated environment, so a proper integration with such
hardware is not possible iet. Software emulation is unfeasible
too: including in a docker container the entire gnome keyring
would skyrocket the size of the image, preventing any meaningful
evaluation of the total size.