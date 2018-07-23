# 4. Destroy
This is part of a series explaining Republic Protocol Darknodes. Head back to the [Overview.]()

#### 4.0 [Deregistering a Darknode and providing feedback.](#40-deregistering-a-darknode-and-providing-feedback-1)
#### 4.1. [Deregistering a Darknode](#41-deregistering-a-darknode-1)
#### 4.2. [Destroying a Darknode](#42-destroying-a-darknode-1)
#### 4.3. [Withdrawing REN from a Darknode](#43-withdrawing-ren-from-a-darknode-1)
#### 4.4. [Providing feedback](#44-providing-feedback-1)


## 4.0. Deregistering a Darknode and providing feedback. 
This section contains everything you need to know to deregister, destroy and withdraw your collateral from a Darknode. 

To monitor the status of your Darknode once it is installed and deployed, see section 3: Operation.


## 4.1. Deregistering a Darknode
To deregister your Darknode, navigate to your Darknode Operator Dashboard (accessible at ‘https:/darknode.republicprotocol.com/status/xx.xxx.xx.xxx'. If you’ve forgotten the address, you can access it by following  3.1.).

Remove all funds from the Darknode. You will need to sign for these transactions with your wallet. 

Press ‘Deregister’. You’ll be asked to sign for the transaction to receive the 100,000REN collateral. 

## 4.2. Destroy a Darknode
Destroy a Darknode
/WARNING: Before destroying a Darknode make sure you have deregistered it, and withdrawn all fees earned! See step 4.1/

Destroying a Darknode will turn it off and tear down all resources allocated by the cloud provider. To destroy a Darknode, open a terminal and run:

```sh
darknode destroy --name my-first-darknode
```

To avoid the command-line prompt reminding you to deregister your Darknode, use the `--force` argument:

```sh
darknode destroy --name my-first-darknode --force
```

We do not recommend using the `--force` argument unless you are developing custom tools that manage your Darknodes automatically.

## 4.3. Withdrawing REN from a Darknode


## 4.4. Providing feedback
