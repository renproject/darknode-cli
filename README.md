# Darknode CLI

The Darknode CLI is a command-line interface for managing Darknodes on Republic Protocol.

## Getting Started on AWS

To jump straight into running your first Darknode on AWS, checkout our [Getting Started on AWS](./docs/getting-started-on-aws.md) tutorial.

## Getting Started on Digital Ocean

> Coming soon!

## Installation

To download and install the Darknode CLI, open a terminal and run:

```sh
curl https://darknode.republicprotocol.com/install.sh -sSf | sh
```

This will download the required binaries and templates and install them to the `$HOME/.darknode` directory. Open a new terminal to begin using the Darknode CLI.

## Usage 

### Deploy a Darknode

#### AWS

To deploy a Darknode on AWS, open a terminal and run:

```sh
darknode up --name my-first-darknode --aws --aws-access-key YOUR-AWS-ACCESS-KEY --aws-secret-key YOUR-AWS-SECRET-KEY
``` 

The Darknode CLI will automatically use the credentials available at `$HOME/.aws/credentials` if you do not explicitly set the `--access-key` and `--secret-key` arguments.

You can also specify the region and instance type you want to use for the Darknode:

```sh
darknode up --name my-first-darknode --aws --aws-access-key YOUR-AWS-ACCESS-KEY --aws-secret-key YOUR-AWS-SECRET-KEY --aws-region eu-west-1 --aws-instance t2.small
``` 

You can find all available regions and instance types at [AWS](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html).

You can also associate the darknode to an Elastic IP by specifing the `allocation_id`. 
Make sure you give the same region of the Elastic to the darknode. 

```sh
darknode up --name my-first-darknode --aws --aws-access-key YOUR-AWS-ACCESS-KEY --aws-secret-key YOUR-AWS-SECRET-KEY --aws-region same-region-as-EIP
``` 

#### Digital Ocean

> Coming soon!

### Destroy a Darknode

_WARNING: Before destroying a Darknode make sure you have deregistered it, and withdrawn all fees earned!_

Destroying a Darknode will turn it off and tear down all resources allocated by the cloud provider. To destroy a Darknode, open a terminal and run:

```sh
darknode destroy --name my-first-darknode
```

### List all Darknodes

The Darknode CLI supports deploying multiple Darknodes. To list all available Darknodes, open a terminal and run:

```sh
darknode list
```

### Start Darknode

To turn on your darknode, open a terminal and run: 

```sh
darknode start --name my-first-darknode
``` 

If it's already on, `start` will do nothing.

### Stop Darknode

To turn off your darknode, open a terminal and run: 

```sh
darknode stop --name my-first-darknode
``` 

If it's already off, `stop` will do nothing.

### SSH into Darknode

To access your Darknode using SSH, open a terminal and run:

```sh
darknode ssh --name my-first-darknode
``` 

### Update a Darknode

To update your Darknode to the latest stable version, open a terminal and run:

```sh
darknode update --name my-first-darknode
``` 

To update the config of your darknode, first edit the local version of config, by running:

```sh
nano $HOME/.darknode/darknodes/YOUR-NODE-NAME/config.json
``` 

Then run 

```sh
darknode update --name my-first-darknode --config
``` 
