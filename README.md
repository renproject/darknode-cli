# darknode-deployer
The tool for deploying darknode. 

## Usage 

### Deploy a new darknode on AWS

To download and install the darknode deployer, run the command in shell 

```bash
curl https://darknode.republicprotocol.com/darknode.sh -sSf | sh
```  
If you're using a custom shell other than bash, you have to include the `./darknode/bin` folder in the `PATH` variable.

To deploy a darknode on AWS, run the command below with your access-key and secret-key replaced.
```bash
$ darknode up --provider=aws --access_key=YOUR-AWS-ACCESS-KEY --secret_key=YOUR-AWS-SECRET-KEY
``` 

You can also specify the region and instance type you want to use for the darknode.

```bash
$ darknode up --provider=aws --access_key=YOUR-AWS-ACCESS-KEY --secret_key=YOUR-AWS-SECRET-KEY --region=eu-west-1 --instance=t2.small
``` 

You can find all available regions and available instance type Sof each region from [AWS](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html)

When finishing deploying, an url will be displayed in the terminal. You can track your node status by opening the url in your browser.


### Deploy a new darknode on digitalOcean

> TBD

### Deploy multiple nodes 

> TBD

### Destroy a node

Destroy will stop the darknode and tear down the instance from your service provider.
It will keep the node config in case you want use the same config.

```bash
$ darknode destroy
# Or
$ darknode down
``` 

You will be asked to deregister your node and withdrawn your fees before your destroy node.
Enter "Yes" to continue or "No" to cancel.
You can also specify the "--skip" flag to ignore this question and start destroying node straight.

```bash
$ darknode destroy --skip
# Or
$ darknode down --skip
``` 

### SSH into your node

```bash
$ darknode ssh  # Make sure you are in the right directory
``` 

### Update your node 

```bash
$ darknode update # Make sure you are in the right directory
``` 

### Help 

If you have any question, you can use the help command.

```bash
$ darknode help
``` 
