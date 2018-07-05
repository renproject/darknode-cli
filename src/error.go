package main

import (
	"fmt"
)

const (
	RESET = "\x1b[0m"
	GREEN = "\x1b[32;1m"
	RED   = "\x1b[31;1m"
)

// ErrKeyNotFound is returned when no AWS access-key nor secret-key provided.
var ErrKeyNotFound = fmt.Errorf("%splease provide your AWS access key and secret key%s", RED, RESET)

// ErrNodeExist is returned when user tries to created a new node with name
// already exists.
var ErrNodeExist = fmt.Errorf("%snode with same name already exists%s", RED, RESET)

// ErrMultipleProviders is returned when user gives more than one provider.
var ErrMultipleProviders = fmt.Errorf("%splease give only one provider%s", RED, RESET)

// ErrUnknownProvider is returned when user wants to deploy darknode to an
// unknown service provider
var ErrUnknownProvider = fmt.Errorf("%sunknown service provider%s", RED, RESET)

// ErrNilProvider is returned when the provider is nil.
var ErrNilProvider = fmt.Errorf("%sprovider cannot be nil%s", RED, RESET)

// UnknownRegion is returned when the provided region is not valid on AWS.
var UnknownRegion = fmt.Errorf("%sthere is no such region on AWS%s", RED, RESET)

// UnSupportedInstanceType is returned when the provided instance is not
// supported in the selected region.
var UnSupportedInstanceType = fmt.Errorf("%sinstance type is not supported in the region%s", RED, RESET)

// ErrNoNodesFound is returned when no nodes can be found with the given tag.
var ErrNoNodesFound = fmt.Errorf("%sno nodes can be found with the given tag%s", RED, RESET)

// ErrNoDeploymentFound is returned when no node can be found for destroying
var ErrNoDeploymentFound = fmt.Errorf("%scannot find any deployed node%s", RED, RESET)

// ErrEmptyNodeName is returned when user doesn't provide the node name.
var ErrEmptyNodeName = fmt.Errorf("%snode name cannot be empty%s", RED, RESET)
