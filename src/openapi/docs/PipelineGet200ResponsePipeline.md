# PipelineGet200ResponsePipeline

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | Pointer to [**PipelineStatus**](PipelineStatus.md) |  | [optional] 
**LastStart** | Pointer to **int64** | Milliseconds since epoch when the pipeline was manually started | [optional] 
**LastStop** | Pointer to **int64** | Milliseconds since epoch when the pipeline was manually stopped | [optional] 
**LastRestart** | Pointer to **int64** | Milliseconds since epoch when the pipeline was automatically restarted (on process faults) | [optional] 

## Methods

### NewPipelineGet200ResponsePipeline

`func NewPipelineGet200ResponsePipeline() *PipelineGet200ResponsePipeline`

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

### HasStatus

`func (o *PipelineGet200ResponsePipeline) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

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

### HasLastStart

`func (o *PipelineGet200ResponsePipeline) HasLastStart() bool`

HasLastStart returns a boolean if a field has been set.

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

### HasLastStop

`func (o *PipelineGet200ResponsePipeline) HasLastStop() bool`

HasLastStop returns a boolean if a field has been set.

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

### HasLastRestart

`func (o *PipelineGet200ResponsePipeline) HasLastRestart() bool`

HasLastRestart returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


