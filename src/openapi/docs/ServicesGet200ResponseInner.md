# ServicesGet200ResponseInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the service | 
**Status** | [**ServiceStatus**](ServiceStatus.md) |  | 
**EnabledVersion** | **string** | The version that is enabled for this service (if any) | 

## Methods

### NewServicesGet200ResponseInner

`func NewServicesGet200ResponseInner(name string, status ServiceStatus, enabledVersion string, ) *ServicesGet200ResponseInner`

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



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


