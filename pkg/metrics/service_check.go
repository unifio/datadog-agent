package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gogo/protobuf/proto"

	agentpayload "github.com/DataDog/agent-payload/gogen"
)

// ServiceCheckStatus represents the status associated with a service check
type ServiceCheckStatus int

// Enumeration of the existing service check statuses, and their values
const (
	ServiceCheckOK       ServiceCheckStatus = 0
	ServiceCheckWarning  ServiceCheckStatus = 1
	ServiceCheckCritical ServiceCheckStatus = 2
	ServiceCheckUnknown  ServiceCheckStatus = 3
)

// GetServiceCheckStatus returns the ServiceCheckStatus from and integer value
func GetServiceCheckStatus(val int) (ServiceCheckStatus, error) {
	switch val {
	case int(ServiceCheckOK):
		return ServiceCheckOK, nil
	case int(ServiceCheckWarning):
		return ServiceCheckWarning, nil
	case int(ServiceCheckCritical):
		return ServiceCheckCritical, nil
	case int(ServiceCheckUnknown):
		return ServiceCheckUnknown, nil
	default:
		return ServiceCheckUnknown, fmt.Errorf("invalid value for a ServiceCheckStatus")
	}
}

// String returns a string representation of ServiceCheckStatus
func (s ServiceCheckStatus) String() string {
	switch s {
	case ServiceCheckOK:
		return "OK"
	case ServiceCheckWarning:
		return "WARNING"
	case ServiceCheckCritical:
		return "CRITICAL"
	case ServiceCheckUnknown:
		return "UNKNOWN"
	default:
		return ""
	}
}

// ServiceCheck holds a service check (w/ serialization to DD api format)
type ServiceCheck struct {
	CheckName string             `json:"check"`
	Host      string             `json:"host_name"`
	Ts        int64              `json:"timestamp"`
	Status    ServiceCheckStatus `json:"status"`
	Message   string             `json:"message"`
	Tags      []string           `json:"tags"`
}

// ServiceChecks represents a list of service checks ready to be serialize
type ServiceChecks []*ServiceCheck

// Marshal serialize service checks using agent-payload definition
func (sc ServiceChecks) Marshal() ([]byte, error) {
	payload := &agentpayload.ServiceChecksPayload{
		ServiceChecks: []*agentpayload.ServiceChecksPayload_ServiceCheck{},
		Metadata:      &agentpayload.CommonMetadata{},
	}

	for _, c := range sc {
		payload.ServiceChecks = append(payload.ServiceChecks,
			&agentpayload.ServiceChecksPayload_ServiceCheck{
				Name:    c.CheckName,
				Host:    c.Host,
				Ts:      c.Ts,
				Status:  int32(c.Status),
				Message: c.Message,
				Tags:    c.Tags,
			})
	}

	return proto.Marshal(payload)
}

// MarshalJSON serializes service checks to JSON so it can be sent to V1 endpoints
//FIXME(olivier): to be removed when v2 endpoints are available
func (sc ServiceChecks) MarshalJSON() ([]byte, error) {
	// use an alias to avoid infinit recursion while serializing
	type ServiceChecksAlias ServiceChecks

	reqBody := &bytes.Buffer{}
	err := json.NewEncoder(reqBody).Encode(ServiceChecksAlias(sc))
	return reqBody.Bytes(), err
}