# ServicesGet200ResponseInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | The name of the service | [optional] 
**Status** | Pointer to [**ServiceStatus**](ServiceStatus.md) |  | [optional] 
**EnabledVersion** | Pointer to **string** | The version that is enabled for this service (if any) | [optional] 

## Methods

### NewServicesGet200ResponseInner

`func NewServicesGet200ResponseInner() *ServicesGet200ResponseInner`

NewServicesGet200ResponseInner instantiates a new ServicesGet200ResponseInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServicesGet200ResponseInnerWithDefaults

`func NewServicesGet200ResponseInnerWithDefaults() *ServicesGet200ResponseInner`

NewServicesGet200ResponseInnerWithDefaults instantiates a new ServicesGet200ResponseInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ServicesGet200ResponseInner) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ServicesGet200ResponseInner) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ServicesGet200ResponseInner) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ServicesGet200ResponseInner) HasName() bool`

HasName returns a boolean if a field has been set.

### GetStatus

`func (o *ServicesGet200ResponseInner) GetStatus() ServiceStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ServicesGet200ResponseInner) GetStatusOk() (*ServiceStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ServicesGet200ResponseInner) SetStatus(v ServiceStatus)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ServicesGet200ResponseInner) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetEnabledVersion

`func (o *ServicesGet200ResponseInner) GetEnabledVersion() string`

GetEnabledVersion returns the EnabledVersion field if non-nil, zero value otherwise.

### GetEnabledVersionOk

`func (o *ServicesGet200ResponseInner) GetEnabledVersionOk() (*string, bool)`

GetEnabledVersionOk returns a tuple with the EnabledVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabledVersion

`func (o *ServicesGet200ResponseInner) SetEnabledVersion(v string)`

SetEnabledVersion sets EnabledVersion field to given value.

### HasEnabledVersion

`func (o *ServicesGet200ResponseInner) HasEnabledVersion() bool`

HasEnabledVersion returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


