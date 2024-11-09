# ServicesNameVersionGet200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | The name of the service | [optional] 
**Version** | Pointer to **string** | The version of the service | [optional] 
**Status** | Pointer to [**ServiceStatus**](ServiceStatus.md) |  | [optional] 
**BuiltAt** | Pointer to **int64** | The time this version was last built as milliseconds since epoch | [optional] 
**Author** | Pointer to **string** | The author of the service | [optional] 
**Inputs** | Pointer to [**[]ServicesNameVersionGet200ResponseInputsInner**](ServicesNameVersionGet200ResponseInputsInner.md) | The dependencies/inputs of this service version | [optional] 
**Outputs** | Pointer to **[]string** | The output streams of this service version | [optional] 
**Errors** | Pointer to **[]string** | The validation errors of this service version (one error per line) | [optional] 

## Methods

### NewServicesNameVersionGet200Response

`func NewServicesNameVersionGet200Response() *ServicesNameVersionGet200Response`

NewServicesNameVersionGet200Response instantiates a new ServicesNameVersionGet200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServicesNameVersionGet200ResponseWithDefaults

`func NewServicesNameVersionGet200ResponseWithDefaults() *ServicesNameVersionGet200Response`

NewServicesNameVersionGet200ResponseWithDefaults instantiates a new ServicesNameVersionGet200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ServicesNameVersionGet200Response) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ServicesNameVersionGet200Response) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ServicesNameVersionGet200Response) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ServicesNameVersionGet200Response) HasName() bool`

HasName returns a boolean if a field has been set.

### GetVersion

`func (o *ServicesNameVersionGet200Response) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *ServicesNameVersionGet200Response) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *ServicesNameVersionGet200Response) SetVersion(v string)`

SetVersion sets Version field to given value.

### HasVersion

`func (o *ServicesNameVersionGet200Response) HasVersion() bool`

HasVersion returns a boolean if a field has been set.

### GetStatus

`func (o *ServicesNameVersionGet200Response) GetStatus() ServiceStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ServicesNameVersionGet200Response) GetStatusOk() (*ServiceStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ServicesNameVersionGet200Response) SetStatus(v ServiceStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ServicesNameVersionGet200Response) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetBuiltAt

`func (o *ServicesNameVersionGet200Response) GetBuiltAt() int64`

GetBuiltAt returns the BuiltAt field if non-nil, zero value otherwise.

### GetBuiltAtOk

`func (o *ServicesNameVersionGet200Response) GetBuiltAtOk() (*int64, bool)`

GetBuiltAtOk returns a tuple with the BuiltAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBuiltAt

`func (o *ServicesNameVersionGet200Response) SetBuiltAt(v int64)`

SetBuiltAt sets BuiltAt field to given value.

### HasBuiltAt

`func (o *ServicesNameVersionGet200Response) HasBuiltAt() bool`

HasBuiltAt returns a boolean if a field has been set.

### GetAuthor

`func (o *ServicesNameVersionGet200Response) GetAuthor() string`

GetAuthor returns the Author field if non-nil, zero value otherwise.

### GetAuthorOk

`func (o *ServicesNameVersionGet200Response) GetAuthorOk() (*string, bool)`

GetAuthorOk returns a tuple with the Author field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthor

`func (o *ServicesNameVersionGet200Response) SetAuthor(v string)`

SetAuthor sets Author field to given value.

### HasAuthor

`func (o *ServicesNameVersionGet200Response) HasAuthor() bool`

HasAuthor returns a boolean if a field has been set.

### GetInputs

`func (o *ServicesNameVersionGet200Response) GetInputs() []ServicesNameVersionGet200ResponseInputsInner`

GetInputs returns the Inputs field if non-nil, zero value otherwise.

### GetInputsOk

`func (o *ServicesNameVersionGet200Response) GetInputsOk() (*[]ServicesNameVersionGet200ResponseInputsInner, bool)`

GetInputsOk returns a tuple with the Inputs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInputs

`func (o *ServicesNameVersionGet200Response) SetInputs(v []ServicesNameVersionGet200ResponseInputsInner)`

SetInputs sets Inputs field to given value.

### HasInputs

`func (o *ServicesNameVersionGet200Response) HasInputs() bool`

HasInputs returns a boolean if a field has been set.

### GetOutputs

`func (o *ServicesNameVersionGet200Response) GetOutputs() []string`

GetOutputs returns the Outputs field if non-nil, zero value otherwise.

### GetOutputsOk

`func (o *ServicesNameVersionGet200Response) GetOutputsOk() (*[]string, bool)`

GetOutputsOk returns a tuple with the Outputs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutputs

`func (o *ServicesNameVersionGet200Response) SetOutputs(v []string)`

SetOutputs sets Outputs field to given value.

### HasOutputs

`func (o *ServicesNameVersionGet200Response) HasOutputs() bool`

HasOutputs returns a boolean if a field has been set.

### GetErrors

`func (o *ServicesNameVersionGet200Response) GetErrors() []string`

GetErrors returns the Errors field if non-nil, zero value otherwise.

### GetErrorsOk

`func (o *ServicesNameVersionGet200Response) GetErrorsOk() (*[]string, bool)`

GetErrorsOk returns a tuple with the Errors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErrors

`func (o *ServicesNameVersionGet200Response) SetErrors(v []string)`

SetErrors sets Errors field to given value.

### HasErrors

`func (o *ServicesNameVersionGet200Response) HasErrors() bool`

HasErrors returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


