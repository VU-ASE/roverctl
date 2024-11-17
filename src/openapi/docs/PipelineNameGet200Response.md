# PipelineNameGet200Response

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
**ServiceName** | **string** | The name of the service that this process is running | 
**ServiceVersion** | **string** | The version of the service that this process is running | 
**Logs** | Pointer to **[]string** | The latest &lt;log_lines&gt; log lines of the process | [optional] 

## Methods

### NewPipelineNameGet200Response

`func NewPipelineNameGet200Response(name string, status ProcessStatus, pid int32, uptime int64, memory int32, cpu int32, faults int32, serviceName string, serviceVersion string, ) *PipelineNameGet200Response`

NewPipelineNameGet200Response instantiates a new PipelineNameGet200Response object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPipelineNameGet200ResponseWithDefaults

`func NewPipelineNameGet200ResponseWithDefaults() *PipelineNameGet200Response`

NewPipelineNameGet200ResponseWithDefaults instantiates a new PipelineNameGet200Response object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *PipelineNameGet200Response) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PipelineNameGet200Response) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PipelineNameGet200Response) SetName(v string)`

SetName sets Name field to given value.


### GetStatus

`func (o *PipelineNameGet200Response) GetStatus() ProcessStatus`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *PipelineNameGet200Response) GetStatusOk() (*ProcessStatus, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *PipelineNameGet200Response) SetStatus(v ProcessStatus)`

SetStatus sets Status field to given value.


### GetPid

`func (o *PipelineNameGet200Response) GetPid() int32`

GetPid returns the Pid field if non-nil, zero value otherwise.

### GetPidOk

`func (o *PipelineNameGet200Response) GetPidOk() (*int32, bool)`

GetPidOk returns a tuple with the Pid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPid

`func (o *PipelineNameGet200Response) SetPid(v int32)`

SetPid sets Pid field to given value.


### GetUptime

`func (o *PipelineNameGet200Response) GetUptime() int64`

GetUptime returns the Uptime field if non-nil, zero value otherwise.

### GetUptimeOk

`func (o *PipelineNameGet200Response) GetUptimeOk() (*int64, bool)`

GetUptimeOk returns a tuple with the Uptime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUptime

`func (o *PipelineNameGet200Response) SetUptime(v int64)`

SetUptime sets Uptime field to given value.


### GetMemory

`func (o *PipelineNameGet200Response) GetMemory() int32`

GetMemory returns the Memory field if non-nil, zero value otherwise.

### GetMemoryOk

`func (o *PipelineNameGet200Response) GetMemoryOk() (*int32, bool)`

GetMemoryOk returns a tuple with the Memory field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMemory

`func (o *PipelineNameGet200Response) SetMemory(v int32)`

SetMemory sets Memory field to given value.


### GetCpu

`func (o *PipelineNameGet200Response) GetCpu() int32`

GetCpu returns the Cpu field if non-nil, zero value otherwise.

### GetCpuOk

`func (o *PipelineNameGet200Response) GetCpuOk() (*int32, bool)`

GetCpuOk returns a tuple with the Cpu field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCpu

`func (o *PipelineNameGet200Response) SetCpu(v int32)`

SetCpu sets Cpu field to given value.


### GetFaults

`func (o *PipelineNameGet200Response) GetFaults() int32`

GetFaults returns the Faults field if non-nil, zero value otherwise.

### GetFaultsOk

`func (o *PipelineNameGet200Response) GetFaultsOk() (*int32, bool)`

GetFaultsOk returns a tuple with the Faults field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFaults

`func (o *PipelineNameGet200Response) SetFaults(v int32)`

SetFaults sets Faults field to given value.


### GetServiceName

`func (o *PipelineNameGet200Response) GetServiceName() string`

GetServiceName returns the ServiceName field if non-nil, zero value otherwise.

### GetServiceNameOk

`func (o *PipelineNameGet200Response) GetServiceNameOk() (*string, bool)`

GetServiceNameOk returns a tuple with the ServiceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceName

`func (o *PipelineNameGet200Response) SetServiceName(v string)`

SetServiceName sets ServiceName field to given value.


### GetServiceVersion

`func (o *PipelineNameGet200Response) GetServiceVersion() string`

GetServiceVersion returns the ServiceVersion field if non-nil, zero value otherwise.

### GetServiceVersionOk

`func (o *PipelineNameGet200Response) GetServiceVersionOk() (*string, bool)`

GetServiceVersionOk returns a tuple with the ServiceVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceVersion

`func (o *PipelineNameGet200Response) SetServiceVersion(v string)`

SetServiceVersion sets ServiceVersion field to given value.


### GetLogs

`func (o *PipelineNameGet200Response) GetLogs() []string`

GetLogs returns the Logs field if non-nil, zero value otherwise.

### GetLogsOk

`func (o *PipelineNameGet200Response) GetLogsOk() (*[]string, bool)`

GetLogsOk returns a tuple with the Logs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLogs

`func (o *PipelineNameGet200Response) SetLogs(v []string)`

SetLogs sets Logs field to given value.

### HasLogs

`func (o *PipelineNameGet200Response) HasLogs() bool`

HasLogs returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


