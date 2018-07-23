# 2. Install
This is part of a series explaining Republic Protocol Darknodes. Head back to the [Overview.]()

#### 2.0. [Installing and deploying a Darknode.](#20-installing-and-deploying-a-darknode-1)
#### 2.1. [Receiving Kovan Testnet Tokens ](#21-receiving-kovan-testnet-tokens)
#### 2.2. [Setting up Amazon Web Servers (AWS)](#22-setting-up-amazon-web-servers-aws-1)
#### 2.3. [Downloading and installing the Darknode CLI](#23-downloading-and-installing-the-darknode-cli-1)
#### 2.4. [Deploy a Darknode](#24-deploy-a-darknode-1)
#### 2.5. [Depositing ETH to operate Darknode](#25-depositing-eth-to-operate-darknode-1)
#### 2.6. [Registering a Darknode](#26-registering-a-darknode-1)



## 2.0. Installing and deploying a Darknode.
This section explains how to install and setup your Darknode once you’ve been accepted. To apply to be a Darknode Operator, [visit section 1: Apply](./01-apply.md).



## 2.1. Receiving Kovan Testnet Tokens
When your application is approved, the Republic Protocol Bot will ask for your wallet’s Kovan public key to deposit tokens into.

This is so you can deposit the 100,000REN required to operate a Darknode. Darknode’s also take a small amount of ETH to operate (about 0.1ETH per week). 

Input your Kovan Test Network public key into the chat with the Republic Protocol Bot. 

The tokens will appear in your wallet shortly. Please allow up to 1 hour. 

Note: Make sure you retain at least 1 ETH and 100,000REN to operate the Darknode.


## 2.2. Setting up Amazon Web Servers (AWS)

If you’ve already done this, you can skip to step 2.3. Downloading and installing the Darknode CLI. 


## 2.3. Downloading and installing the Darknode CLI
To download and install the Darknode CLI, open a terminal and run:

```
curl https://darknode.republicprotocol.com/install.sh -sSf | sh
```

This will download the required binaries and templates and install them to the $HOME/.darknode directory. 
 
Close the terminal.


## 2.4. Deploy a Darknode

To deploy your Darknode after install, open a new terminal and run: 

```
darknode up --name my-first-darknode --aws --aws-access-key YOUR-AWS-ACCESS-KEY --aws-secret-key YOUR-AWS-SECRET-KEY
```

The Darknode CLI will automatically use the credentials available at $HOME/.aws/credentials if you do not explicitly set the --access-key and --secret-key arguments.

You can also specify the region and instance type you want to use for the Darknode:

```
darknode up --name my-first-darknode --aws --aws-access-key YOUR-AWS-ACCESS-KEY --aws-secret-key YOUR-AWS-SECRET-KEY --aws-region eu-west-1 --aws-instance t2.small
```

You can find all available regions and instance types at [AWS](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.RegionsAndAvailabilityZones.html).

You can also associate the Darknode to an elastic IP by specifying the allocation_id. Make sure you give the same region of the elastic to the Darknode.

```
darknode up --name my-first-darknode --aws --aws-access-key YOUR-AWS-ACCESS-KEY --aws-secret-key YOUR-AWS-SECRET-KEY --aws-region same-region-as-EIP -aws-elastic-ip XXX.XXX.XXX.XXX
```

The install will take a few minutes. Please be patient. Once it has completed you will see a message “Congratulations! Your Darknode is deployed and running.”

Take note of the URL contained within the success message. It will read ‘https:/darknode.republicprotocol.com/status/xx.xxx.xx.xxx'


## 2.5. Depositing ETH to operate Darknode

Once you have deployed your Darknode, to get it operational you’ll need to deposit the required ETH and REN. Visit the URL presented to you in step 2.4. This is your Darknode Operator Dashboard. 

From here, you can register and deregister your Darknode. For a Darknode to be operational (and earning fees) it needs to be registered on the Republic Protocol network. 

From the Operator Dashboard you can monitor and top up the Darknode’s current ETH Balance. It requires ETH for gas. It will consume about 0.1 ETH per week of operation. 

You can also withdraw fees earned for conducting computations as part of the Republic Protocol network. For more information on earning fees, visit our Darknode FAQ. 

Under the ‘Top-up Balance’ section enter the amount of ETH you would like to deposit. We recommend 0.1 ETH. 

You’ll be asked to sign for the transaction by your wallet. 

Once approved, the balance under the ‘Balance’ section will reflect your deposit. 

Note: Republic Protocol is interested in improving the systems in place for monitoring the current ETH balance of your Darknode. Let us know how you think it could be improved. 


## 2.6. Registering a Darknode
To register your Darknode, click the Register button present on the Darknode Operation Dashboard.

You’ll be asked to approve the transfer of 100,000REN by your wallet. 

Confirm the transaction and wait. Your Darknode will be registered, however you will need to wait up to 24 hours. 



---
Once approved, to operate and maintain your Darknode, head to [Section 3: Operate.]()


