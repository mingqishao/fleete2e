#!/bin/bash

resourceGroup="${1:-myRG}"
cluster="${2:-myCluster}"

az account set -s 8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8
az group create -n $resourceGroup -l eastus
body=$(cat <<-END
{
   "location":"eastus",
   "identity":{
      "type":"SystemAssigned"
   },
   "properties":{
      "dnsPrefix":"$cluster",
      "agentPoolProfiles":[
         {
            "name":"agentpool1",
            "count":1,
            "vmSize":"Standard_D2s_v3",
            "osType":"Linux",
            "type":"VirtualMachineScaleSets",
            "mode":"System"
         }
      ],
      "sku":{
         "name":"Basic",
         "tier":"Paid"
      }
   }
}
END
)
curl -v -X PUT \
 -H "x-ms-identity-url: http://msi-simulator.msi-simulator.svc.cluster.local" \
 -H "x-ms-identity-principal-id: 8ff738a5-abcd-4864-a162-6c18f7c9cbd9" \
 -H "x-ms-home-tenant-id: 72f988bf-86f1-41af-91ab-2d7cd011db47" \
 -H "Content-Type: application/json" \
 -H "x-ms-home-tenant-id: 72f988bf-86f1-41af-91ab-2d7cd011db47" \
 --data "$body" \
"http://localhost:8081/subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourceGroups/$resourceGroup/providers/Microsoft.ContainerService/managedClusters/$cluster?api-version=2022-04-01"

