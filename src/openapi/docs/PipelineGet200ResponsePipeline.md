# PipelineGet200ResponsePipeline

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | [**PipelineStatus**](PipelineStatus.md) |  | 
**LastStart** | **int64** | Milliseconds since epoch when the pipeline was manually started | 
**LastStop** | **int64** | Milliseconds since epoch when the pipeline was manually stopped | 
**LastRestart** | **int64** | Milliseconds since epoch when the pipeline was automatically restarted (on process faults) | 
**ValidationErrors** | Pointer to [**PipelineGet200ResponsePipelineValidationErrors**](PipelineGet200ResponsePipelineValidationErrors.md) |  | [optional] 

## Methods

### NewPipelineGet200ResponsePipeline

`func NewPipelineGet200ResponsePipeline(status PipelineStatus, lastStart int64, lastStop int64, lastRestart int64, ) *PipelineGet200ResponsePipeline`

NewPipelineGet200ResponsePipeline instantiates a new PipelineGet200ResponsePipeline object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPipelineGet200ResponsePipelineWithDefaults

`func NewPipelineGet200ResponsePipelineWithDefaults() *PipelineGet200ResponsePipeline`

NewPipelineGet200ResponsePipelineWithDefaults instantiates a new PipelineGet200ResponsePipeline object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *PipelineGet200ResponsePipeline) GetStatus() PipelineStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *PipelineGet200ResponsePipeline) GetStatusOk() (*PipelineStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *PipelineGet200ResponsePipeline) SetStatus(v PipelineStatus)`

SetStatus sets Status field to given value.


### GetLastStart

`func (o *PipelineGet200ResponsePipeline) GetLastStart() int64`

GetLastStart returns the LastStart field if non-nil, zero value otherwise.

### GetLastStartOk

`func (o *PipelineGet200ResponsePipeline) GetLastStartOk() (*int64, bool)`

GetLastStartOk returns a tuple with the LastStart field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastStart

`func (o *PipelineGet200ResponsePipeline) SetLastStart(v int64)`

SetLastStart sets LastStart field to given value.


### GetLastStop

`func (o *PipelineGet200ResponsePipeline) GetLastStop() int64`

GetLastStop returns the LastStop field if non-nil, zero value otherwise.

### GetLastStopOk

`func (o *PipelineGet200ResponsePipeline) GetLastStopOk() (*int64, bool)`

GetLastStopOk returns a tuple with the LastStop field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastStop

`func (o *PipelineGet200ResponsePipeline) SetLastStop(v int64)`

SetLastStop sets LastStop field to given value.


### GetLastRestart

`func (o *PipelineGet200ResponsePipeline) GetLastRestart() int64`

GetLastRestart returns the LastRestart field if non-nil, zero value otherwise.

### GetLastRestartOk

`func (o *PipelineGet200ResponsePipeline) GetLastRestartOk() (*int64, bool)`

GetLastRestartOk returns a tuple with the LastRestart field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLastRestart

`func (o *PipelineGet200ResponsePipeline) SetLastRestart(v int64)`

SetLastRestart sets LastRestart field to given value.


### GetValidationErrors

`func (o *PipelineGet200ResponsePipeline) GetValidationErrors() PipelineGet200ResponsePipelineValidationErrors`

GetValidationErrors returns the ValidationErrors field if non-nil, zero value otherwise.

### GetValidationErrorsOk

`func (o *PipelineGet200ResponsePipeline) GetValidationErrorsOk() (*PipelineGet200ResponsePipelineValidationErrors, bool)`

GetValidationErrorsOk returns a tuple with the ValidationErrors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValidationErrors

`func (o *PipelineGet200ResponsePipeline) SetValidationErrors(v PipelineGet200ResponsePipelineValidationErrors)`

SetValidationErrors sets ValidationErrors field to given value.

### HasValidationErrors

`func (o *PipelineGet200ResponsePipeline) HasValidationErrors() bool`

HasValidationErrors returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


