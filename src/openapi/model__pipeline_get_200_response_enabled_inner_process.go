/*
roverd REST API

API exposed from each rover to allow process, service, source and file management

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"bytes"
	"fmt"
)

// checks if the PipelineGet200ResponseEnabledInnerProcess type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PipelineGet200ResponseEnabledInnerProcess{}

// PipelineGet200ResponseEnabledInnerProcess The last process that was started for this service (instantiated from the service). This can be undefined if the pipeline was not started before.
type PipelineGet200ResponseEnabledInnerProcess struct {
	// The process ID. Depending on the status, this PID might not exist anymore
	Pid int32 `json:"pid"`
	Status ProcessStatus `json:"status"`
	// The number of milliseconds the process has been running
	Uptime int64 `json:"uptime"`
	// The amount of memory used by the process in megabytes
	Memory int32 `json:"memory"`
	// The percentage of CPU used by the process
	Cpu int32 `json:"cpu"`
}

type _PipelineGet200ResponseEnabledInnerProcess PipelineGet200ResponseEnabledInnerProcess

// NewPipelineGet200ResponseEnabledInnerProcess instantiates a new PipelineGet200ResponseEnabledInnerProcess object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPipelineGet200ResponseEnabledInnerProcess(pid int32, status ProcessStatus, uptime int64, memory int32, cpu int32) *PipelineGet200ResponseEnabledInnerProcess {
	this := PipelineGet200ResponseEnabledInnerProcess{}
	this.Pid = pid
	this.Status = status
	this.Uptime = uptime
	this.Memory = memory
	this.Cpu = cpu
	return &this
}

// NewPipelineGet200ResponseEnabledInnerProcessWithDefaults instantiates a new PipelineGet200ResponseEnabledInnerProcess object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPipelineGet200ResponseEnabledInnerProcessWithDefaults() *PipelineGet200ResponseEnabledInnerProcess {
	this := PipelineGet200ResponseEnabledInnerProcess{}
	return &this
}

// GetPid returns the Pid field value
func (o *PipelineGet200ResponseEnabledInnerProcess) GetPid() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Pid
}

// GetPidOk returns a tuple with the Pid field value
// and a boolean to check if the value has been set.
func (o *PipelineGet200ResponseEnabledInnerProcess) GetPidOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Pid, true
}

// SetPid sets field value
func (o *PipelineGet200ResponseEnabledInnerProcess) SetPid(v int32) {
	o.Pid = v
}

// GetStatus returns the Status field value
func (o *PipelineGet200ResponseEnabledInnerProcess) GetStatus() ProcessStatus {
	if o == nil {
		var ret ProcessStatus
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *PipelineGet200ResponseEnabledInnerProcess) GetStatusOk() (*ProcessStatus, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *PipelineGet200ResponseEnabledInnerProcess) SetStatus(v ProcessStatus) {
	o.Status = v
}

// GetUptime returns the Uptime field value
func (o *PipelineGet200ResponseEnabledInnerProcess) GetUptime() int64 {
	if o == nil {
		var ret int64
		return ret
	}

	return o.Uptime
}

// GetUptimeOk returns a tuple with the Uptime field value
// and a boolean to check if the value has been set.
func (o *PipelineGet200ResponseEnabledInnerProcess) GetUptimeOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Uptime, true
}

// SetUptime sets field value
func (o *PipelineGet200ResponseEnabledInnerProcess) SetUptime(v int64) {
	o.Uptime = v
}

// GetMemory returns the Memory field value
func (o *PipelineGet200ResponseEnabledInnerProcess) GetMemory() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Memory
}

// GetMemoryOk returns a tuple with the Memory field value
// and a boolean to check if the value has been set.
func (o *PipelineGet200ResponseEnabledInnerProcess) GetMemoryOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Memory, true
}

// SetMemory sets field value
func (o *PipelineGet200ResponseEnabledInnerProcess) SetMemory(v int32) {
	o.Memory = v
}

// GetCpu returns the Cpu field value
func (o *PipelineGet200ResponseEnabledInnerProcess) GetCpu() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Cpu
}

// GetCpuOk returns a tuple with the Cpu field value
// and a boolean to check if the value has been set.
func (o *PipelineGet200ResponseEnabledInnerProcess) GetCpuOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Cpu, true
}

// SetCpu sets field value
func (o *PipelineGet200ResponseEnabledInnerProcess) SetCpu(v int32) {
	o.Cpu = v
}

func (o PipelineGet200ResponseEnabledInnerProcess) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PipelineGet200ResponseEnabledInnerProcess) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["pid"] = o.Pid
	toSerialize["status"] = o.Status
	toSerialize["uptime"] = o.Uptime
	toSerialize["memory"] = o.Memory
	toSerialize["cpu"] = o.Cpu
	return toSerialize, nil
}

func (o *PipelineGet200ResponseEnabledInnerProcess) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"pid",
		"status",
		"uptime",
		"memory",
		"cpu",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varPipelineGet200ResponseEnabledInnerProcess := _PipelineGet200ResponseEnabledInnerProcess{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varPipelineGet200ResponseEnabledInnerProcess)

	if err != nil {
		return err
	}

	*o = PipelineGet200ResponseEnabledInnerProcess(varPipelineGet200ResponseEnabledInnerProcess)

	return err
}

type NullablePipelineGet200ResponseEnabledInnerProcess struct {
	value *PipelineGet200ResponseEnabledInnerProcess
	isSet bool
}

func (v NullablePipelineGet200ResponseEnabledInnerProcess) Get() *PipelineGet200ResponseEnabledInnerProcess {
	return v.value
}

func (v *NullablePipelineGet200ResponseEnabledInnerProcess) Set(val *PipelineGet200ResponseEnabledInnerProcess) {
	v.value = val
	v.isSet = true
}

func (v NullablePipelineGet200ResponseEnabledInnerProcess) IsSet() bool {
	return v.isSet
}

func (v *NullablePipelineGet200ResponseEnabledInnerProcess) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePipelineGet200ResponseEnabledInnerProcess(val *PipelineGet200ResponseEnabledInnerProcess) *NullablePipelineGet200ResponseEnabledInnerProcess {
	return &NullablePipelineGet200ResponseEnabledInnerProcess{value: val, isSet: true}
}

func (v NullablePipelineGet200ResponseEnabledInnerProcess) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePipelineGet200ResponseEnabledInnerProcess) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

