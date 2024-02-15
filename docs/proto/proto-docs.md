<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [cosmwasm/will/params.proto](#cosmwasm/will/params.proto)
    - [Params](#cosmwasm.will.Params)
  
- [cosmwasm/will/genesis.proto](#cosmwasm/will/genesis.proto)
    - [GenesisState](#cosmwasm.will.GenesisState)
  
- [cosmwasm/will/types.proto](#cosmwasm/will/types.proto)
    - [ClaimComponent](#cosmwasm.will.ClaimComponent)
    - [ExecutionComponent](#cosmwasm.will.ExecutionComponent)
    - [GnarkZkSnark](#cosmwasm.will.GnarkZkSnark)
    - [PedersenCommitment](#cosmwasm.will.PedersenCommitment)
    - [SchnorrSignature](#cosmwasm.will.SchnorrSignature)
    - [TransferComponent](#cosmwasm.will.TransferComponent)
    - [Will](#cosmwasm.will.Will)
    - [WillIds](#cosmwasm.will.WillIds)
    - [Wills](#cosmwasm.will.Wills)
  
- [cosmwasm/will/query.proto](#cosmwasm/will/query.proto)
    - [QueryGetWillRequest](#cosmwasm.will.QueryGetWillRequest)
    - [QueryGetWillResponse](#cosmwasm.will.QueryGetWillResponse)
    - [QueryListWillsRequest](#cosmwasm.will.QueryListWillsRequest)
    - [QueryListWillsResponse](#cosmwasm.will.QueryListWillsResponse)
  
    - [Query](#cosmwasm.will.Query)
  
- [cosmwasm/will/tx.proto](#cosmwasm/will/tx.proto)
    - [MsgCheckInRequest](#cosmwasm.will.MsgCheckInRequest)
    - [MsgCheckInResponse](#cosmwasm.will.MsgCheckInResponse)
    - [MsgCreateWillRequest](#cosmwasm.will.MsgCreateWillRequest)
    - [MsgCreateWillResponse](#cosmwasm.will.MsgCreateWillResponse)
    - [MsgUpdateParams](#cosmwasm.will.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#cosmwasm.will.MsgUpdateParamsResponse)
  
    - [Msg](#cosmwasm.will.Msg)
  
- [Scalar Value Types](#scalar-value-types)



<a name="cosmwasm/will/params.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmwasm/will/params.proto



<a name="cosmwasm.will.Params"></a>

### Params
Params defines the parameters for the module.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmwasm/will/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmwasm/will/genesis.proto



<a name="cosmwasm.will.GenesisState"></a>

### GenesisState
GenesisState defines the will module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#cosmwasm.will.Params) |  | params defines all the parameters of the module. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmwasm/will/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmwasm/will/types.proto



<a name="cosmwasm.will.ClaimComponent"></a>

### ClaimComponent
ClaimComponent is designed for actions requiring a claim with proof.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `evidence` | [string](#string) |  | Evidence required to support the claim. |
| `pedersen` | [PedersenCommitment](#cosmwasm.will.PedersenCommitment) |  | Represents a Pedersen commitment scheme. |
| `schnorr` | [SchnorrSignature](#cosmwasm.will.SchnorrSignature) |  | Represents a Schnorr signature scheme. |
| `gnark` | [GnarkZkSnark](#cosmwasm.will.GnarkZkSnark) |  | Represents a zk-SNARK scheme using Gnark. |






<a name="cosmwasm.will.ExecutionComponent"></a>

### ExecutionComponent
ExecutionComponent defines a single actionable component within a will.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `transfer` | [TransferComponent](#cosmwasm.will.TransferComponent) |  | Represents an asset transfer action. |
| `claim` | [ClaimComponent](#cosmwasm.will.ClaimComponent) |  | Represents a claim action that a beneficiary must perform. |






<a name="cosmwasm.will.GnarkZkSnark"></a>

### GnarkZkSnark
GnarkZkSnark is for claims using zero-knowledge succinct non-interactive
arguments of knowledge.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `verification_key` | [bytes](#bytes) |  | The public key for verifying the zk-SNARK proof. |
| `public_inputs` | [bytes](#bytes) |  | Public inputs required for the proof verification. |
| `proof` | [bytes](#bytes) |  | The zk-SNARK proof demonstrating knowledge of a secret |






<a name="cosmwasm.will.PedersenCommitment"></a>

### PedersenCommitment
PedersenCommitment enables the use of a Pedersen commitment for claims.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commitment` | [bytes](#bytes) |  | The commitment hash, representing the hidden value. |
| `random_factor` | [bytes](#bytes) |  | The random factor used to generate the commitment, |
| `value` | [bytes](#bytes) |  | enhancing privacy.

The actual value being committed, revealed during claim. |
| `blinding_factor` | [bytes](#bytes) |  | The blinding factor used alongside the value for claim verification. |






<a name="cosmwasm.will.SchnorrSignature"></a>

### SchnorrSignature
SchnorrSignature is used for claims that require a Schnorr signature.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `public_key` | [bytes](#bytes) |  | The public key associated with the signature. |
| `unique_session_id` | [bytes](#bytes) |  | A unique identifier for the session, enhancing security. |
| `signature` | [bytes](#bytes) |  | The digital signature for claim verification. |
| `message` | [string](#string) |  | An optional message that may accompany the signature. |






<a name="cosmwasm.will.TransferComponent"></a>

### TransferComponent
TransferComponent is used for direct asset transfers.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to` | [string](#string) |  | Destination address for the asset transfer. |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | The amount and type of asset to transfer. |






<a name="cosmwasm.will.Will"></a>

### Will
Will represents the entire structure of a will.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | Unique identifier for the will. |
| `name` | [string](#string) |  | User-generated name for the will. |
| `beneficiary` | [string](#string) |  | The designated beneficiary or receiver of the will's assets. |
| `height` | [int64](#int64) |  | The designated block to trigger the will |
| `components` | [ExecutionComponent](#cosmwasm.will.ExecutionComponent) | repeated | The list of execution components that make up the will. |






<a name="cosmwasm.will.WillIds"></a>

### WillIds
WillIds represents a list of will IDs.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `ids` | [string](#string) | repeated |  |






<a name="cosmwasm.will.Wills"></a>

### Wills
type to hold wills


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `wills` | [Will](#cosmwasm.will.Will) | repeated | the set of wills to return |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmwasm/will/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmwasm/will/query.proto



<a name="cosmwasm.will.QueryGetWillRequest"></a>

### QueryGetWillRequest
QueryGetWillRequest is the request type for retrieving a will by its ID.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `will_id` | [string](#string) |  |  |






<a name="cosmwasm.will.QueryGetWillResponse"></a>

### QueryGetWillResponse
QueryGetWillResponse is the response type returned after retrieving a will.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `will` | [Will](#cosmwasm.will.Will) |  | will is the will body |






<a name="cosmwasm.will.QueryListWillsRequest"></a>

### QueryListWillsRequest
QueryListWillsRequest request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the contract to query |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<a name="cosmwasm.will.QueryListWillsResponse"></a>

### QueryListWillsResponse
QueryListWillsRequest response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `wills` | [Will](#cosmwasm.will.Will) | repeated | the will struct for the entries of the response |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmwasm.will.Query"></a>

### Query
Query defines the gRPC querier service for the will module.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `GetWill` | [QueryGetWillRequest](#cosmwasm.will.QueryGetWillRequest) | [QueryGetWillResponse](#cosmwasm.will.QueryGetWillResponse) | GetWill retrieves a will by its ID. | GET|/cosmwasm/wasmd/will/{will_id}|
| `ListWills` | [QueryListWillsRequest](#cosmwasm.will.QueryListWillsRequest) | [QueryListWillsResponse](#cosmwasm.will.QueryListWillsResponse) | GetWill retrieves all wills by an account address | GET|/cosmwasm/wasmd/will/list/{address}|

 <!-- end services -->



<a name="cosmwasm/will/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmwasm/will/tx.proto



<a name="cosmwasm.will.MsgCheckInRequest"></a>

### MsgCheckInRequest
message for checking in


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `creator` | [string](#string) |  |  |
| `id` | [string](#string) |  |  |
| `height` | [int64](#int64) |  |  |






<a name="cosmwasm.will.MsgCheckInResponse"></a>

### MsgCheckInResponse
response for checkin


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `status` | [bool](#bool) |  |  |
| `height` | [int64](#int64) |  |  |






<a name="cosmwasm.will.MsgCreateWillRequest"></a>

### MsgCreateWillRequest
message for creating a will


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `creator` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `beneficiary` | [string](#string) |  |  |
| `height` | [int64](#int64) |  |  |
| `components` | [ExecutionComponent](#cosmwasm.will.ExecutionComponent) | repeated |  |






<a name="cosmwasm.will.MsgCreateWillResponse"></a>

### MsgCreateWillResponse
to get the will response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `creator` | [string](#string) |  |  |
| `name` | [string](#string) |  |  |
| `beneficiary` | [string](#string) |  |  |
| `height` | [int64](#int64) |  |  |






<a name="cosmwasm.will.MsgUpdateParams"></a>

### MsgUpdateParams
MsgUpdateParams is the Msg/UpdateParams request type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `authority` | [string](#string) |  | authority is the address that controls the module (defaults to x/gov unless overwritten). |
| `params` | [Params](#cosmwasm.will.Params) |  | params defines the module parameters to update.

NOTE: All parameters must be supplied. |






<a name="cosmwasm.will.MsgUpdateParamsResponse"></a>

### MsgUpdateParamsResponse
MsgUpdateParamsResponse defines the response structure for executing a
MsgUpdateParams message.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="cosmwasm.will.Msg"></a>

### Msg
Msg defines the Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `UpdateParams` | [MsgUpdateParams](#cosmwasm.will.MsgUpdateParams) | [MsgUpdateParamsResponse](#cosmwasm.will.MsgUpdateParamsResponse) | UpdateParams defines a (governance) operation for updating the module parameters. The authority defaults to the x/gov module account. | |
| `CreateWill` | [MsgCreateWillRequest](#cosmwasm.will.MsgCreateWillRequest) | [MsgCreateWillResponse](#cosmwasm.will.MsgCreateWillResponse) | create a new will | |
| `CheckIn` | [MsgCheckInRequest](#cosmwasm.will.MsgCheckInRequest) | [MsgCheckInResponse](#cosmwasm.will.MsgCheckInResponse) | checkin into a will | |

 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

