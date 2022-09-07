#!/bin/bash

assignee="$1@microsoft.com"
echo $assignee
objectID=$(az ad user list --filter "mail eq 'wenjungao@microsoft.com'" --query "[].id" -o tsv)
az role assignment create --role "Contributor" --assignee-object-id $objectID --subscription "3959ec86-5353-4b0c-b5d7-3877122861a0"
az role assignment create --role "Azure Kubernetes Fleet Manager RBAC Admin" --assignee-object-id $objectID --subscription "3959ec86-5353-4b0c-b5d7-3877122861a0"
az role assignment create --role "Azure Kubernetes Fleet Manager Contributor Role" --assignee-object-id $objectID --subscription "3959ec86-5353-4b0c-b5d7-3877122861a0"
az role assignment create --role "Azure Kubernetes Service Cluster Admin Role" --assignee-object-id $objectID --subscription "3959ec86-5353-4b0c-b5d7-3877122861a0"
az role assignment create --role "Azure Kubernetes Service Contributor Role" --assignee-object-id $objectID --subscription "3959ec86-5353-4b0c-b5d7-3877122861a0"