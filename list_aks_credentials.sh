#!/bin/bash

cluster="${1:-myCluster}"
resourceGroup="$cluster-rg"

curl -v -X POST \
 -H "x-ms-home-tenant-id: 72f988bf-86f1-41af-91ab-2d7cd011db47" \
 http://localhost:8081/subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourceGroups/$resourceGroup/providers/Microsoft.ContainerService/managedClusters/$cluster/listClusterAdminCredential?api-version=2022-04-01
