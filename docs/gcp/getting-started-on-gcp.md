# Getting Started on Google Cloud Platform

Before installing and using the Darknode CLI, make sure that you have [created a Google Account](https://accounts.google.com/SignUp) and configured billing for [Google Cloud Platform](https://cloud.google.com/billing/docs/how-to/manage-billing-account), and [created a Project](https://cloud.google.com/resource-manager/docs/creating-managing-projects).

# Enabling the Resource Manager API

Navigate to the [API Library](https://console.developers.google.com/apis/library/cloudresourcemanager.googleapis.com), make sure your project is selected and enable the Cloud Resource Manager API.

![Visual representation of a block](enable-api.png)

# Create a Service Account

Create a [Service Account](https://cloud.google.com/iam/docs/creating-managing-service-accounts) in your project. During creating, grant it the role of Project \> Editor, and download a key in JSON format. You can choose a name freely.

![Creating a service account](create-sa-1.png)
![Creating a service account](create-sa-2.png)
![Creating a service account](create-sa-3.png)


