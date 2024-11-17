# PipelineGet200ResponsePipelineValidationErrors

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**UnmetStreams** | Pointer to [**[]UnmetStreamError**](UnmetStreamError.md) |  | [optional] 
**UnmetServices** | Pointer to [**[]UnmetServiceError**](UnmetServiceError.md) |  | [optional] 
**DuplicateService** | Pointer to **[]string** |  | [optional] 

## Methods

### NewPipelineGet200ResponsePipelineValidationErrors

`func NewPipelineGet200ResponsePipelineValidationErrors() *PipelineGet200ResponsePipelineValidationErrors`

NewPipelineGet200ResponsePipelineValidationErrors instantiates a new PipelineGet200ResponsePipelineValidationErrors object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPipelineGet200ResponsePipelineValidationErrorsWithDefaults

`func NewPipelineGet200ResponsePipelineValidationErrorsWithDefaults() *PipelineGet200ResponsePipelineValidationErrors`

NewPipelineGet200ResponsePipelineValidationErrorsWithDefaults instantiates a new PipelineGet200ResponsePipelineValidationErrors object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUnmetStreams

`func (o *PipelineGet200ResponsePipelineValidationErrors) GetUnmetStreams() []UnmetStreamError`

GetUnmetStreams returns the UnmetStreams field if non-nil, zero value otherwise.

### GetUnmetStreamsOk

`func (o *PipelineGet200ResponsePipelineValidationErrors) GetUnmetStreamsOk() (*[]UnmetStreamError, bool)`

GetUnmetStreamsOk returns a tuple with the UnmetStreams field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnmetStreams

`func (o *PipelineGet200ResponsePipelineValidationErrors) SetUnmetStreams(v []UnmetStreamError)`

SetUnmetStreams sets UnmetStreams field to given value.

### HasUnmetStreams

`func (o *PipelineGet200ResponsePipelineValidationErrors) HasUnmetStreams() bool`

HasUnmetStreams returns a boolean if a field has been set.

### GetUnmetServices

`func (o *PipelineGet200ResponsePipelineValidationErrors) GetUnmetServices() []UnmetServiceError`

GetUnmetServices returns the UnmetServices field if non-nil, zero value otherwise.

### GetUnmetServicesOk

`func (o *PipelineGet200ResponsePipelineValidationErrors) GetUnmetServicesOk() (*[]UnmetServiceError, bool)`

GetUnmetServicesOk returns a tuple with the UnmetServices field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnmetServices

`func (o *PipelineGet200ResponsePipelineValidationErrors) SetUnmetServices(v []UnmetServiceError)`

SetUnmetServices sets UnmetServices field to given value.

### HasUnmetServices

`func (o *PipelineGet200ResponsePipelineValidationErrors) HasUnmetServices() bool`

HasUnmetServices returns a boolean if a field has been set.

### GetDuplicateService

`func (o *PipelineGet200ResponsePipelineValidationErrors) GetDuplicateService() []string`

GetDuplicateService returns the DuplicateService field if non-nil, zero value otherwise.

### GetDuplicateServiceOk

`func (o *PipelineGet200ResponsePipelineValidationErrors) GetDuplicateServiceOk() (*[]string, bool)`

GetDuplicateServiceOk returns a tuple with the DuplicateService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuplicateService

`func (o *PipelineGet200ResponsePipelineValidationErrors) SetDuplicateService(v []string)`

SetDuplicateService sets DuplicateService field to given value.

### HasDuplicateService

`func (o *PipelineGet200ResponsePipelineValidationErrors) HasDuplicateService() bool`

HasDuplicateService returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


