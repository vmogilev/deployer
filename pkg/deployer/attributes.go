package deployer

import "fmt"

const (
	// GlobalAttributeNamePrefix - global Attributes such as IP, HOSTNAME
	// are parsed at the start of run-list execution and injected into All
	// directive data Attributes, which can then reference them in templates
	GlobalAttributeNamePrefix = "GLOBAL_"
)

// Globals - global Attributes shared by All directives
func (x *SDK) Globals() (*Attributes, error) {
	return &Attributes{
		All: map[string]interface{}{
			GlobalAttributeNamePrefix + "IP": GetOutboundIP().String(),
		},
	}, nil
}

type Attributes struct {
	All    map[string]interface{}
	Errors []error
}

func (a *Attributes) StringValue(key string) string {
	v, ok := a.All[key]
	if !ok {
		a.Errors = append(a.Errors, fmt.Errorf("%q key not found", key))
		return ""
	}
	s, ok := v.(string)
	if !ok {
		a.Errors = append(a.Errors, fmt.Errorf("%q key value is not a string: %v", key, v))
		return ""
	}
	return s
}
