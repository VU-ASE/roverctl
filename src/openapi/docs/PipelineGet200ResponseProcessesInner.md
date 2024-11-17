# PipelineGet200ResponseProcessesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the service running as a process | 
**Status** | [**ProcessStatus**](ProcessStatus.md) |  | 
**Pid** | **int32** | The process ID | 
**Uptime** | **int64** | The number of milliseconds the process has been running | 
**Memory** | **int32** | The amount of memory used by the process in megabytes | 
**Cpu** | **int32** | The percentage of CPU used by the process | 
**Faults** | **int32** | The number of faults that have occurred (causing the pipeline to restart) since last_start | 

## Methods

### NewPipelineGet200ResponseProcessesInner

`func NewPipelineGet200ResponseProcessesInner(name string, status ProcessStatus, pid int32, uptime int64, memory int32, cpu int32, faults int32, ) *PipelineGet200ResponseProcessesInner`

NewPipelineGet200ResponseProcessesInner instantiates a new PipelineGet200ResponseProcessesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPipelineGet200ResponseProcessesInnerWithDefaults

`func NewPipelineGet200ResponseProcessesInnerWithDefaults() *PipelineGet200ResponseProcessesInner`

NewPipelineGet200ResponseProcessesInnerWithDefaults instantiates a new PipelineGet200ResponseProcessesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *PipelineGet200ResponseProcessesInner) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PipelineGet200ResponseProcessesInner) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PipelineGet200ResponseProcessesInner) SetName(v string)`

SetName sets Name field to given value.


### GetStatus

`func (o *PipelineGet200ResponseProcessesInner) GetStatus() ProcessStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *PipelineGet200ResponseProcessesInner) GetStatusOk() (*ProcessStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *PipelineGet200ResponseProcessesInner) SetStatus(v ProcessStatus)`

SetStatus sets Status field to given value.


### GetPid

`func (o *PipelineGet200ResponseProcessesInner) GetPid() int32`

GetPid returns the Pid field if non-nil, zero value otherwise.

### GetPidOk

`func (o *PipelineGet200ResponseProcessesInner) GetPidOk() (*int32, bool)`

GetPidOk returns a tuple with the Pid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPid

`func (o *PipelineGet200ResponseProcessesInner) SetPid(v int32)`

SetPid sets Pid field to given value.


### GetUptime

`func (o *PipelineGet200ResponseProcessesInner) GetUptime() int64`

GetUptime returns the Uptime field if non-nil, zero value otherwise.

### GetUptimeOk

`func (o *PipelineGet200ResponseProcessesInner) GetUptimeOk() (*int64, bool)`

GetUptimeOk returns a tuple with the Uptime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUptime

`func (o *PipelineGet200ResponseProcessesInner) SetUptime(v int64)`

SetUptime sets Uptime field to given value.


### GetMemory

`func (o *PipelineGet200ResponseProcessesInner) GetMemory() int32`

GetMemory returns the Memory field if non-nil, zero value otherwise.

### GetMemoryOk

`func (o *PipelineGet200ResponseProcessesInner) GetMemoryOk() (*int32, bool)`

GetMemoryOk returns a tuple with the Memory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemory

`func (o *PipelineGet200ResponseProcessesInner) SetMemory(v int32)`

SetMemory sets Memory field to given value.


### GetCpu

`func (o *PipelineGet200ResponseProcessesInner) GetCpu() int32`

GetCpu returns the Cpu field if non-nil, zero value otherwise.

### GetCpuOk

`func (o *PipelineGet200ResponseProcessesInner) GetCpuOk() (*int32, bool)`

GetCpuOk returns a tuple with the Cpu field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpu

`func (o *PipelineGet200ResponseProcessesInner) SetCpu(v int32)`

SetCpu sets Cpu field to given value.


### GetFaults

`func (o *PipelineGet200ResponseProcessesInner) GetFaults() int32`

GetFaults returns the Faults field if non-nil, zero value otherwise.

### GetFaultsOk

`func (o *PipelineGet200ResponseProcessesInner) GetFaultsOk() (*int32, bool)`

GetFaultsOk returns a tuple with the Faults field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFaults

`func (o *PipelineGet200ResponseProcessesInner) SetFaults(v int32)`

SetFaults sets Faults field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


