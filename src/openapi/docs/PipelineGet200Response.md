# PipelineGet200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Pipeline** | Pointer to [**PipelineGet200ResponsePipeline**](PipelineGet200ResponsePipeline.md) |  | [optional] 
**Processes** | Pointer to [**[]PipelineGet200ResponseProcessesInner**](PipelineGet200ResponseProcessesInner.md) |  | [optional] 

## Methods

### NewPipelineGet200Response

`func NewPipelineGet200Response() *PipelineGet200Response`

NewPipelineGet200Response instantiates a new PipelineGet200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPipelineGet200ResponseWithDefaults

`func NewPipelineGet200ResponseWithDefaults() *PipelineGet200Response`

NewPipelineGet200ResponseWithDefaults instantiates a new PipelineGet200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPipeline

`func (o *PipelineGet200Response) GetPipeline() PipelineGet200ResponsePipeline`

GetPipeline returns the Pipeline field if non-nil, zero value otherwise.

### GetPipelineOk

`func (o *PipelineGet200Response) GetPipelineOk() (*PipelineGet200ResponsePipeline, bool)`

GetPipelineOk returns a tuple with the Pipeline field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPipeline

`func (o *PipelineGet200Response) SetPipeline(v PipelineGet200ResponsePipeline)`

SetPipeline sets Pipeline field to given value.

### HasPipeline

`func (o *PipelineGet200Response) HasPipeline() bool`

HasPipeline returns a boolean if a field has been set.

### GetProcesses

`func (o *PipelineGet200Response) GetProcesses() []PipelineGet200ResponseProcessesInner`

GetProcesses returns the Processes field if non-nil, zero value otherwise.

### GetProcessesOk

`func (o *PipelineGet200Response) GetProcessesOk() (*[]PipelineGet200ResponseProcessesInner, bool)`

GetProcessesOk returns a tuple with the Processes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProcesses

`func (o *PipelineGet200Response) SetProcesses(v []PipelineGet200ResponseProcessesInner)`

SetProcesses sets Processes field to given value.

### HasProcesses

`func (o *PipelineGet200Response) HasProcesses() bool`

HasProcesses returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


