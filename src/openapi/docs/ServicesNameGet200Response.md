# ServicesNameGet200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the service | 
**Status** | [**ServiceStatus**](ServiceStatus.md) |  | 
**Versions** | **[]string** |  | 
**EnabledVersion** | Pointer to **string** | The version that is enabled for this service (if any) | [optional] 

## Methods

### NewServicesNameGet200Response

`func NewServicesNameGet200Response(name string, status ServiceStatus, versions []string, ) *ServicesNameGet200Response`

NewServicesNameGet200Response instantiates a new ServicesNameGet200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServicesNameGet200ResponseWithDefaults

`func NewServicesNameGet200ResponseWithDefaults() *ServicesNameGet200Response`

NewServicesNameGet200ResponseWithDefaults instantiates a new ServicesNameGet200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ServicesNameGet200Response) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ServicesNameGet200Response) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ServicesNameGet200Response) SetName(v string)`

SetName sets Name field to given value.


### GetStatus

`func (o *ServicesNameGet200Response) GetStatus() ServiceStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ServicesNameGet200Response) GetStatusOk() (*ServiceStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ServicesNameGet200Response) SetStatus(v ServiceStatus)`

SetStatus sets Status field to given value.


### GetVersions

`func (o *ServicesNameGet200Response) GetVersions() []string`

GetVersions returns the Versions field if non-nil, zero value otherwise.

### GetVersionsOk

`func (o *ServicesNameGet200Response) GetVersionsOk() (*[]string, bool)`

GetVersionsOk returns a tuple with the Versions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersions

`func (o *ServicesNameGet200Response) SetVersions(v []string)`

SetVersions sets Versions field to given value.


### GetEnabledVersion

`func (o *ServicesNameGet200Response) GetEnabledVersion() string`

GetEnabledVersion returns the EnabledVersion field if non-nil, zero value otherwise.

### GetEnabledVersionOk

`func (o *ServicesNameGet200Response) GetEnabledVersionOk() (*string, bool)`

GetEnabledVersionOk returns a tuple with the EnabledVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabledVersion

`func (o *ServicesNameGet200Response) SetEnabledVersion(v string)`

SetEnabledVersion sets EnabledVersion field to given value.

### HasEnabledVersion

`func (o *ServicesNameGet200Response) HasEnabledVersion() bool`

HasEnabledVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


