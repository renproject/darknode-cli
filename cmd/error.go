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

// ErrNodeNotExist is returned when user tries to refund a node which doesn't
// exist
var ErrNodeNotExist = fmt.Errorf("%snode with the name doesn't exist%s", RED, RESET)

// ErrMultipleProviders is returned when user gives more than one provider.
var ErrMultipleProviders = fmt.Errorf("%splease give only one provider%s", RED, RESET)

// ErrUnknownProvider is returned when user wants to deploy darknode to an
// unknown service provider
var ErrUnknownProvider = fmt.Errorf("%sunknown service provider%s", RED, RESET)

// ErrNilProvider is returned when the provider is nil.
var ErrNilProvider = fmt.Errorf("%sprovider cannot be nil%s", RED, RESET)

// ErrUnknownRegion is returned when the provided region is not valid.
var ErrUnknownRegion = fmt.Errorf("%sthere is no such region or the region is not available%s", RED, RESET)

// ErrUnSupportedInstanceType is returned when the provided instance is not
// supported in the selected region.
var ErrUnSupportedInstanceType = fmt.Errorf("%sinstance type is not supported in the region%s", RED, RESET)

// ErrNoNodesFound is returned when no nodes can be found with the given tag.
var ErrNoNodesFound = fmt.Errorf("%sno nodes can be found with the given tag%s", RED, RESET)

// ErrNoDeploymentFound is returned when no node can be found for destroying
var ErrNoDeploymentFound = fmt.Errorf("%scannot find any deployed node%s", RED, RESET)

// ErrEmptyNodeName is returned when user doesn't provide the node name.
var ErrEmptyNodeName = fmt.Errorf("%snode name cannot be empty%s", RED, RESET)

// ErrUnknownNetwork is returned when user wants to deploy darknode to an
// unknown darkpool network
var ErrUnknownNetwork = fmt.Errorf("%sunknown network%s", RED, RESET)

// ErrNameAndTags is returned when both name and tags are given.
var ErrNameAndTags = fmt.Errorf("%stoo many arguments, cannot have both name and tags%s", RED, RESET)

// ErrEmptyNameAndTags is returned when both name and tags are not given.
var ErrEmptyNameAndTags = fmt.Errorf("%splease provide name or tags of the node you want to operate%s", RED, RESET)

// ErrFilePath is returned when user doesn't provide the file path.
var ErrEmptyFilePath = fmt.Errorf("%sfile path cannot be empty%s", RED, RESET)

// ErrEmptyDoToken is returned when the digital ocean token is empty.
var ErrEmptyDoToken = fmt.Errorf("%sdigital ocean token cannot be empty%s", RED, RESET)

// ErrUnknownDropletSize is returned when the user gives an unknown droplet size
var ErrUnknownDropletSize = fmt.Errorf("%sunknown droplet size%s", RED, RESET)

// ErrInvalidEthereumAddress is returned when user gives an invalid Ethereum address.
var ErrInvalidEthereumAddress = fmt.Errorf("%sinvalid Ethereum address%s", RED, RESET)

// ErrFailedTx is returned when the transaction gets reverted on Ethereum.
var ErrFailedTx = fmt.Errorf("%stransaction failed%s", RED, RESET)

// ErrEmptyAddress is returned when the Ethereum address is empty.
var ErrEmptyAddress = fmt.Errorf("%sethereum address cannot be empty%s", RED, RESET)

// ErrRejectedTx is returned when the tx is rejected by Ethereum.
var ErrRejectedTx = fmt.Errorf("%stransaction rejected by Ethereum%s", RED, RESET)

// ErrUnsupportedOS is returned when the operating system is not supported.
var ErrUnsupportedOS = fmt.Errorf("%sunsupported operating system%s", RED, RESET)

// ErrNoAvailableRegion is returned when the provided account has no available region
var ErrNoAvailableRegion = fmt.Errorf("%sno available region to your account%s", RED, RESET)
