/*
roverd REST API

API exposed from each rover to allow process, service, source and file management

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"fmt"
)

// ServiceStatus The status of any given service is either enabled or disabled
type ServiceStatus string

// List of ServiceStatus
const (
	ENABLED ServiceStatus = "enabled"
	DISABLED ServiceStatus = "disabled"
)

// All allowed values of ServiceStatus enum
var AllowedServiceStatusEnumValues = []ServiceStatus{
	"enabled",
	"disabled",
}

func (v *ServiceStatus) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ServiceStatus(value)
	for _, existing := range AllowedServiceStatusEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid ServiceStatus", value)
}

// NewServiceStatusFromValue returns a pointer to a valid ServiceStatus
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewServiceStatusFromValue(v string) (*ServiceStatus, error) {
	ev := ServiceStatus(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for ServiceStatus: valid values are %v", v, AllowedServiceStatusEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ServiceStatus) IsValid() bool {
	for _, existing := range AllowedServiceStatusEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ServiceStatus value
func (v ServiceStatus) Ptr() *ServiceStatus {
	return &v
}

type NullableServiceStatus struct {
	value *ServiceStatus
	isSet bool
}

func (v NullableServiceStatus) Get() *ServiceStatus {
	return v.value
}

func (v *NullableServiceStatus) Set(val *ServiceStatus) {
	v.value = val
	v.isSet = true
}

func (v NullableServiceStatus) IsSet() bool {
	return v.isSet
}

func (v *NullableServiceStatus) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableServiceStatus(val *ServiceStatus) *NullableServiceStatus {
	return &NullableServiceStatus{value: val, isSet: true}
}

func (v NullableServiceStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableServiceStatus) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
