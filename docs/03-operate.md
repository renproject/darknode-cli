# 3. Operate
This is part of a series explaining Republic Protocol Darknodes. Head back to the [Overview.]()

#### 3.0. [Monitoring and maintaining a Darknode](#30-monitoring-and-maintaining-a-darknode-1)
#### 3.1. [List all Darknodes on the network](#31-list-all-darknodes-on-the-network-1)
#### 3.2. [Starting a Darknode](#32-starting-a-darknode-1)
#### 3.3. [Stopping a Darknode](#33-stopping-a-darknode-1)
#### 3.4. [Updating a Darknode](#34-updating-a-darknode-1)
#### 3.5. [SSH into a Darknode](#35-ssh-into-a-darknode-1)
#### 3.6. [Withdrawing funds from a Darknode](#36-withdrawing-funds-from-a-darknode-1)
#### 3.7. [Deploying multiple Darknodes](#37-deploying-multiple-darknodes-1)


## 3.0. Monitoring and maintaining a Darknode
This section explains how to monitor the status of your Darknode once it is installed and deployed. We’ll also cover various commands to maintain your Darknode. 

If you haven’t yet installed and deployed your Darknode, visit [Section 2: Install.]()


## 3.1. List all Darknodes on the network
To list all available Darknodes, open a terminal and run:
```
darknode list
```


## 3.2. Starting a Darknode
To turn on your darknode, open a terminal and run:
```
darknode start --name my-first-darknode
```

Note: Darknodes will deploy already started. If it's already on, start will do nothing.

## 3.3. Stopping a Darknode
To turn off your darknode, open a terminal and run:
`darknode stop --name my-first-darknode`

Note: If it's already off, stop will do nothing.

## 3.4. Updating a Darknode
To update your Darknode to the latest stable version, open a terminal and run:
`darknode update --name my-first-darknode`

To update the config of your darknode, first edit the local version of config, by running:
`nano $HOME/.darknode/darknodes/YOUR-NODE-NAME/config.json`

Then run:
```
darknode update --name my-first-darknode --config
```



## 3.5. SSH into a Darknode
To access your Darknode using SSH, open a terminal and run:
```
darknode ssh --name my-first-darknode
```

## 3.6. Withdrawing funds from a Darknode

## 3.7. Deploying multiple Darknodes



---
To destroy your Darknode, head to [Section 4: Destroy]()

