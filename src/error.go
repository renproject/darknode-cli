package main

import (
	"fmt"
)

const (
	reset = "\x1b[0m"
	green = "\x1b[32;1m"
	red   = "\x1b[31;1m"
)

// ErrKeyNotFound is returned when no AWS access-key nor secret-key provided.
var ErrKeyNotFound = fmt.Errorf("%splease provide your AWS access key and secret key%s", red, reset)

// ErrNodeExist is returned when user tries to created a new node with name
// already exists.
var ErrNodeExist = fmt.Errorf("%snode with the name already exists%s", red, reset)

// ErrMultipleProviders is returned when user has given more than one provider
var ErrMultipleProviders = fmt.Errorf("%splease give only one provider%s", red, reset)

// ErrUnknownProvider is returned when user wants to deploy darknode to an
// unknown service provider
var ErrUnknownProvider = fmt.Errorf("%sunknown service provider%s", red, reset)

// ErrNilProvider is returned when the provider is nil.
var ErrNilProvider = fmt.Errorf("%sprovider cannot be nil%s", red, reset)

// UnknownRegion is returned when the provided region is not valid on AWS.
var UnknownRegion = fmt.Errorf("%sthere is no such region on AWS%s", red, reset)

// UnSupportedInstanceType is returned when the provided instance is not
// supported in the selected region.
var UnSupportedInstanceType = fmt.Errorf("%sinstance type is not supported in the region%s", red, reset)

// ErrNoNodesFound is returned when no nodes can be found with the given tag.
var ErrNoNodesFound = fmt.Errorf("%sno nodes can be found with the given tag%s", red, reset)
