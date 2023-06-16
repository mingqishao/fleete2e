#!/bin/bash

assignee="$1@microsoft.com"
echo $assignee
objectID=$(az ad user list --filter "mail eq '${assignee}'" --query "[].id" -o tsv)
echo $objectID
az role assignment delete --role "Owner" --assignee $objectID --subscription "3959ec86-5353-4b0c-b5d7-3877122861a0"
az role assignment delete --role "Contributor" --assignee $objectID --subscription "3959ec86-5353-4b0c-b5d7-3877122861a0"
az role assignment delete --role "Azure Kubernetes Fleet Manager RBAC Admin" --assignee $objectID --subscription "3959ec86-5353-4b0c-b5d7-3877122861a0"
az role assignment delete --role "Azure Kubernetes Fleet Manager Contributor Role" --assignee $objectID --subscription "3959ec86-5353-4b0c-b5d7-3877122861a0"
az role assignment delete --role "Azure Kubernetes Service Cluster Admin Role" --assignee $objectID --subscription "3959ec86-5353-4b0c-b5d7-3877122861a0"
az role assignment delete --role "Azure Kubernetes Service Contributor Role" --assignee $objectID --subscription "3959ec86-5353-4b0c-b5d7-3877122861a0"
az role assignment delete --role "Azure Kubernetes Fleet Manager RBAC Cluster Admin" --assignee $objectID --subscription "3959ec86-5353-4b0c-b5d7-3877122861a0"