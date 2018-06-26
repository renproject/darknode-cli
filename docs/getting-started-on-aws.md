# Getting Started on AWS

Before installing and using the Darknode Deployer, make sure that you have created an AWS account.

[Create an AWS account](https://aws.amazon.com)

## Creating a Darknode User Accounts

It is not recommended to use your default account for managing your Darknodes. Instead, we recommend creating a custom user.

![Go to IAM](./getting-started-on-aws-01.png "Go to IAM")

Enter _IAM_ into the search bar and click enter. You should see something like this

![IAM](./getting-started-on-aws-02.png "IAM")

Click _Users_. You should see something like this

![Users](./getting-started-on-aws-03.png "Users")

Click _Add user_. You should see something like this

![Add user](./getting-started-on-aws-04.png "Add user")

Name the user _darknode_ and make sure to check _Programmatic Access_. Then, click _Next: Permissions_. You should see something like this

![Permissions](./getting-started-on-aws-05.png "Permissions")

Click _Create group_. You should see something like this

![Create group](./getting-started-on-aws-06.png "Create group")

Name the group _darknode_ and make sure to check _AmazonEC2FullAccess_. Then, click _Create group_. You should see something like this

![Created group](./getting-started-on-aws-07.png "Created group")

Click _Next: Review_. You should see something like this

![Review](./getting-started-on-aws-08.png "Review")

Click _Create user_. You should see something like this

![Created user](./getting-started-on-aws-09.png "Created user")

Download the CSV file containing your _access key ID_, and your _secret access key_. Make sure you keep these values secure, anyone with access to them will be able to control your Darknode.

## Installing the Darknode Deployer

To install the Darknode Deployer, open a terminal and run:

```sh
curl https://darknode.republicprotocol.com/install.sh -sSf | sh
```

Once this has finished, close the terminal and open a new one.

## Deploying a Darknode

Now, you can deploy a Darknode. Think of a catchy name, and run:

```sh
darknode up --name my-first-darknode --provider aws --access-key YOUR-AWS-ACCESS-KEY --secret-key YOUR-AWS-SECRET-KEY
```

Once this has finished, it will give you a link that you can use to register you Darknode and send it ETH.

Congratulations! You have deployed your first Darknode. You can deploy as many as you like, distinguishing between them by their name. If you forget what you called them, you can list all available Darknodes by running:

```sh
darknode list
```

## Keeping up to date

It is important to keep your Darknode updated to the latest verison. This will ensure that you are compatible with the rest of the network, and that you have the latest security updates and bug fixes.

To update a Darknode, run:

```sh
darknode update --name my-first-darknode
```