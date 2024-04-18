<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [cosmwasm/will/params.proto](#cosmwasm/will/params.proto)
    - [Params](#cosmwasm.will.Params)
  
- [cosmwasm/will/genesis.proto](#cosmwasm/will/genesis.proto)
    - [GenesisState](#cosmwasm.will.GenesisState)
  
- [cosmwasm/will/types.proto](#cosmwasm/will/types.proto)
    - [ClaimAccessControl](#cosmwasm.will.ClaimAccessControl)
    - [ClaimAccessPrivate](#cosmwasm.will.ClaimAccessPrivate)
    - [ClaimAccessPublic](#cosmwasm.will.ClaimAccessPublic)
    - [ClaimComponent](#cosmwasm.will.ClaimComponent)
    - [ContractComponent](#cosmwasm.will.ContractComponent)
    - [ExecutionComponent](#cosmwasm.will.ExecutionComponent)
    - [GnarkZkSnark](#cosmwasm.will.GnarkZkSnark)
    - [IBCContractCall](#cosmwasm.will.IBCContractCall)
    - [IBCMsgComponent](#cosmwasm.will.IBCMsgComponent)
    - [IBCSend](#cosmwasm.will.IBCSend)
    - [OutputContractCall](#cosmwasm.will.OutputContractCall)
    - [OutputTransfer](#cosmwasm.will.OutputTransfer)
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
    - [GnarkClaim](#cosmwasm.will.GnarkClaim)
    - [MsgCheckInRequest](#cosmwasm.will.MsgCheckInRequest)
    - [MsgCheckInResponse](#cosmwasm.will.MsgCheckInResponse)
    - [MsgClaimRequest](#cosmwasm.will.MsgClaimRequest)
    - [MsgClaimResponse](#cosmwasm.will.MsgClaimResponse)
    - [MsgCreateWillRequest](#cosmwasm.will.MsgCreateWillRequest)
    - [MsgCreateWillResponse](#cosmwasm.will.MsgCreateWillResponse)
    - [MsgUpdateParams](#cosmwasm.will.MsgUpdateParams)
    - [MsgUpdateParamsResponse](#cosmwasm.will.MsgUpdateParamsResponse)
    - [PedersenClaim](#cosmwasm.will.PedersenClaim)
    - [SchnorrClaim](#cosmwasm.will.SchnorrClaim)
  
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
| `port_id` | [string](#string) |  | holds the ibc port for the module |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmwasm/will/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmwasm/will/types.proto



<a name="cosmwasm.will.ClaimAccessControl"></a>

### ClaimAccessControl
claim access control


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `public` | [ClaimAccessPublic](#cosmwasm.will.ClaimAccessPublic) |  | public access |
| `private` | [ClaimAccessPrivate](#cosmwasm.will.ClaimAccessPrivate) |  | private access |






<a name="cosmwasm.will.ClaimAccessPrivate"></a>

### ClaimAccessPrivate
claim access private type


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `addresses` | [string](#string) | repeated | set of address to allow claim |






<a name="cosmwasm.will.ClaimAccessPublic"></a>

### ClaimAccessPublic
claim access public type






<a name="cosmwasm.will.ClaimComponent"></a>

### ClaimComponent
ClaimComponent is designed for actions requiring a claim with proof.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `access` | [ClaimAccessControl](#cosmwasm.will.ClaimAccessControl) |  |  |
| `pedersen` | [PedersenCommitment](#cosmwasm.will.PedersenCommitment) |  | Represents a Pedersen commitment scheme. |
| `schnorr` | [SchnorrSignature](#cosmwasm.will.SchnorrSignature) |  | Represents a Schnorr signature scheme. |
| `gnark` | [GnarkZkSnark](#cosmwasm.will.GnarkZkSnark) |  | Represents a zk-SNARK scheme using Gnark. |
| `transfer` | [OutputTransfer](#cosmwasm.will.OutputTransfer) |  | output for simple native transfer |
| `contract_call` | [OutputContractCall](#cosmwasm.will.OutputContractCall) |  | output for native-to-contract call |
| `ibc_contract_call` | [IBCContractCall](#cosmwasm.will.IBCContractCall) |  | output for ibc contract call |
| `ibc_send` | [IBCSend](#cosmwasm.will.IBCSend) |  | output for ibc send |






<a name="cosmwasm.will.ContractComponent"></a>

### ContractComponent
contract component


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | contract address |
| `data` | [bytes](#bytes) |  | data to be passed to the contract |






<a name="cosmwasm.will.ExecutionComponent"></a>

### ExecutionComponent
ExecutionComponent defines a single actionable component within a will.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | component_type enables the inclusion of different types of execution component name |
| `id` | [string](#string) |  | component id |
| `status` | [string](#string) |  | components within a will. |
| `transfer` | [TransferComponent](#cosmwasm.will.TransferComponent) |  | Represents an asset transfer action. |
| `claim` | [ClaimComponent](#cosmwasm.will.ClaimComponent) |  | Represents a claim action that a beneficiary must perform. |
| `contract` | [ContractComponent](#cosmwasm.will.ContractComponent) |  | Future use: Represents an interaction with a smart contract. |
| `ibc_msg` | [IBCMsgComponent](#cosmwasm.will.IBCMsgComponent) |  | future use: for ibc message |






<a name="cosmwasm.will.GnarkZkSnark"></a>

### GnarkZkSnark
GnarkZkSnark is for claims using zero-knowledge succinct non-interactive
arguments of knowledge.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `verification_key` | [bytes](#bytes) |  | The public key for verifying the zk-SNARK proof. |
| `public_inputs` | [bytes](#bytes) |  | Public inputs required for the proof verification. |
| `proof` | [bytes](#bytes) |  | The zk-SNARK proof demonstrating knowledge of a secret |






<a name="cosmwasm.will.IBCContractCall"></a>

### IBCContractCall
for ibc output message, we could make this be contract, or IBC send...


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | contract address |
| `chain` | [bytes](#bytes) |  | destination chain |
| `data` | [bytes](#bytes) |  | data to be passed in the packet |
| `channel` | [bytes](#bytes) |  | channel to be passed in the packet |






<a name="cosmwasm.will.IBCMsgComponent"></a>

### IBCMsgComponent
ibc msg component


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `ibc_contract_call` | [IBCContractCall](#cosmwasm.will.IBCContractCall) |  | output for ibc contract call |
| `ibc_send` | [IBCSend](#cosmwasm.will.IBCSend) |  | output for ibc send |






<a name="cosmwasm.will.IBCSend"></a>

### IBCSend
output for ibc send


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | contract address |
| `chain` | [bytes](#bytes) |  | destination chain |
| `denom` | [bytes](#bytes) |  | data to be passed in the packet |
| `channel` | [bytes](#bytes) |  | channel to be passed in the packet |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | amount to send over IBC |






<a name="cosmwasm.will.OutputContractCall"></a>

### OutputContractCall
output for contract call


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | contract address |
| `data` | [bytes](#bytes) |  | data to be passed to the contract |






<a name="cosmwasm.will.OutputTransfer"></a>

### OutputTransfer
types of outputs for components


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | recipient |
| `denom` | [string](#string) |  | denom to send |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | amout to send |






<a name="cosmwasm.will.PedersenCommitment"></a>

### PedersenCommitment
PedersenCommitment enables the use of a Pedersen commitment for claims.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commitment` | [bytes](#bytes) |  | The commitment hash, representing the hidden value. |
| `target_commitment` | [bytes](#bytes) |  | The target comitment hash |






<a name="cosmwasm.will.SchnorrSignature"></a>

### SchnorrSignature
CLAIM TYPES
SchnorrSignature is used for claims that require a Schnorr signature.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `public_key` | [bytes](#bytes) |  | The public key associated with the signature. |
| `signature` | [bytes](#bytes) |  | The digital signature for claim verification. |
| `message` | [string](#string) |  | An optional message that may accompany the signature. |






<a name="cosmwasm.will.TransferComponent"></a>

### TransferComponent
TransferComponent is used for direct asset transfers.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `to` | [string](#string) |  | Destination address for the asset transfer. |
| `denom` | [string](#string) |  | denom to send |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) |  | The amount and type of asset to transfer. |






<a name="cosmwasm.will.Will"></a>

### Will
Will represents the entire structure of a will.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  | Unique identifier for the will. |
| `creator` | [string](#string) |  | will creator |
| `name` | [string](#string) |  | User-generated name for the will. |
| `beneficiary` | [string](#string) |  | The designated beneficiary or receiver of the will's assets. |
| `height` | [int64](#int64) |  | The designated block to trigger the will |
| `status` | [string](#string) |  | The designated block to trigger the will |
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



<a name="cosmwasm.will.GnarkClaim"></a>

### GnarkClaim
gnark


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `proof` | [bytes](#bytes) |  | Specific fields for Gnark zk-SNARK claims, e.g., proof, public inputs, etc. |
| `public_inputs` | [bytes](#bytes) |  | Public inputs required for the proof verification |






<a name="cosmwasm.will.MsgCheckInRequest"></a>

### MsgCheckInRequest
checkins
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






<a name="cosmwasm.will.MsgClaimRequest"></a>

### MsgClaimRequest
claims


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `will_id` | [string](#string) |  | ID of the will being claimed |
| `claimer` | [string](#string) |  | Claimer's address |
| `component_id` | [string](#string) |  | component id

Add this line |
| `schnorr_claim` | [SchnorrClaim](#cosmwasm.will.SchnorrClaim) |  |  |
| `pedersen_claim` | [PedersenClaim](#cosmwasm.will.PedersenClaim) |  |  |
| `gnark_claim` | [GnarkClaim](#cosmwasm.will.GnarkClaim) |  |  |






<a name="cosmwasm.will.MsgClaimResponse"></a>

### MsgClaimResponse
MsgClaimResponse


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `success` | [bool](#bool) |  | Indicates whether the claim was successful or not |
| `message` | [string](#string) |  | Optional message providing more details on the claim result |






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






<a name="cosmwasm.will.PedersenClaim"></a>

### PedersenClaim
pedersen


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `commitment` | [bytes](#bytes) |  | Specific fields for Pedersen claims, e.g., commitment, blinding factor, etc. |
| `blinding_factor` | [bytes](#bytes) |  |  |
| `value` | [bytes](#bytes) |  | The actual value being revealed in the claim |






<a name="cosmwasm.will.SchnorrClaim"></a>

### SchnorrClaim
SchnorrClaim is specifically structured for claims requiring a Schnorr
signature.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `public_key` | [bytes](#bytes) |  | The public key associated with the Schnorr signature, typically 33 bytes. |
| `signature` | [bytes](#bytes) |  | The actual Schnorr signature, could be 64 bytes (r || s). |
| `message` | [string](#string) |  | The original message that was signed, needed for verification. |





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
| `Claim` | [MsgClaimRequest](#cosmwasm.will.MsgClaimRequest) | [MsgClaimResponse](#cosmwasm.will.MsgClaimResponse) | make a claim | |

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

