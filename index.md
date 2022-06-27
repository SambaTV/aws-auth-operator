 # Samba TV AWS Auth Operator
 
This operator manages the `aws-auth` ConfigMap used by [AWS EKS](https://aws.amazon.com/eks/)
clusters to [map AWS IAM Users and Roles to Kubernetes Users and Groups](https://docs.aws.amazon.com/eks/latest/userguide/add-user-role.html). 

This operator provides two custom resource definitions used to manage the
`data.mapRoles` and `data.mapUsers` arrays with the mappings.

* MapRoles.aws-auth.samba.tv provides a mapping between an AWS IAM Role
  and a Kubernetes user and groups
* MapUsers.aws-auth.samba.tv provides a mapping between an AWS IAM User
  and a Kubernetes user and groups

## Use cases

You should consider using this operator to avoid potential corruption when
manually editing the `aws-auth` ConfigMap.

Don't use this operator in non-AWS EKS clusters as they don't use the
`aws-auth` ConfigMap to configure Kubernetes users and groups.

## Installation

Before you can install the chart you will need to add its repository to `helm`.

```shell
helm repo add aws-auth-operator-charts https://sambatv.github.io/aws-auth-operator
```

After adding the chart repository, you can install the chart from it.

```shell
helm install aws-auth-operator aws-auth-operator-charts/aws-auth-operator --namespace aws-auth-operator --create-namespace
```
