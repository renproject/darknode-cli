# Getting Started on Google Cloud Platform

Before installing and using the Darknode CLI, make sure that you have [created a Google Account](https://accounts.google.com/SignUp) and configured billing for [Google Cloud Platform](https://cloud.google.com/billing/docs/how-to/manage-billing-account), and [created a Project](https://cloud.google.com/resource-manager/docs/creating-managing-projects).

# Enabling the Resource Manager API

Navigate to the [API Library](https://console.developers.google.com/apis/library/cloudresourcemanager.googleapis.com), make sure your project is selected and enable the Cloud Resource Manager API.

![Visual representation of a block](enable-api.png)

# Create a Service Account

Create a [Service Account](https://cloud.google.com/iam/docs/creating-managing-service-accounts) in your project. During creating, grant it the role of Project \> Editor, and download a key in JSON format. You will have to pass this JSON Key file path to the Darknode-CLI

### Creating a Service account
![Creating a service account](create-sa-1.png)
### Assigning a role
![Creating a service account](create-sa-2.png)
![Creating a service account](create-sa-3.png)

### Generating a JSON Key
![Creating a service account](create-sa-4.png)


## Installing the Darknode CLI

To install the Darknode CLI, open a terminal and run:

```sh
curl https://www.github.com/renproject/darknode-cli/releases/latest/download/install.sh -sSfL | sh
```

Once this has finished, close the terminal and open a new one.

## Deploying a Darknode

Now, you can deploy a Darknode. Think of a catchy name, and run:

```sh
darknode up --name my-first-darknode --gcp --gcp-credentials PATH_TO_YOUR_DOWNLOADED_JSON_FILE
```
Once this has finished, it will give you a link that you can use to register your Darknode and send it ETH.

Congratulations! You have deployed your first Darknode. You can deploy as many as you like, distinguishing between them by their names. If you forget what you called them, you can list all available Darknodes by running:

```sh
darknode list
```

### Choosing a compute zone

You can specify in which [Compute Engine Zone](https://cloud.google.com/compute/docs/regions-zones/) you deploy your node with the --gcp-zone flag. If omitted, a random zone is selected.

```sh
darknode up --name my-first-darknode --gcp --gcp-credentials PATH_TO_YOUR_DOWNLOADED_JSON_FILE --gcp-zone europe-west1-b
```
### Choosing a machine type

You can specify with which [Machine Type](https://cloud.google.com/compute/docs/machine-types) you deploy your node with the --gcp-machine-type flag. If omitted, n1-standard-1 is selected.

```sh
darknode up --name my-first-darknode --gcp --gcp-credentials PATH_TO_YOUR_DOWNLOADED_JSON_FILE --gcp-machine-type f1-micro
```

## Keeping up to date

It is important to keep your Darknode updated to the latest verison. This will ensure that you are compatible with the rest of the network, and that you have the latest security updates and bug fixes.

To update a Darknode, run:

```sh
darknode update --name my-first-darknode
```
