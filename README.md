# darknode-deployer
The tool for deploying darknode. 

## Usage 

### Deploy a new darknode on AWS

Download the [darknode.zip](https://republicprotocol.com) and unzip the file.

Open terminal and go to the unzipped folder

```bash
$ ./darknode up --provider=aws --access_key=YOUR-AWS-ACCESS-KEY --secret_key=YOUR-AWS-SECRET-KEY
# Replace `YOUR-AWS-ACCESS-KEY` and `YOUR-AWS-SECRET-KEY` fields with your actual AWS access key and secret key.
``` 
You can also specify the region and instance type you want to use for the darknode.

```bash
$ ./darknode up --provider=aws --access_key=YOUR-AWS-ACCESS-KEY --secret_key=YOUR-AWS-SECRET-KEY --region=eu-west-1 --instance=t2.small
``` 

You can find all available regions and available instance type Sof each region from [AWS](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html)

When finishing deploying, an url will be displayed in the terminal. You can track your node status by opening the url in your browser.


### Deploy a new darknode on digitalOcean

> TBD

### Deploy multiple nodes 

Unzip the compressed file as many time as you like and keep them in separate folders.
The nodes can use different cloud service providers. Just follow the instructions above 
of the provider you want. 

### Destroy a node

Destroy will stop the darknode and tear down the instance from your service provider.
It will keep the node config in case you want use the same config.

```bash
$ ./darknode destroy
``` 

You will be asked to deregister your node and withdrawn your fees before your destroy node.
Enter "Yes" to continue or "No" to cancel.
You can also specify the "--skip" flag to ignore this question and start destroying node straight.

```bash
$ ./darknode destroy --skip
``` 

### SSH into your node

```bash
$ ./darknode ssh  # Make sure you are in the right directory
``` 

### Update your node 

```bash
$ ./darknode update # Make sure you are in the right directory
``` 

### Help 

If you have any question, you can use the help command.

```bash
$ ./darknode help
``` 
