
#!/bin/bash
set -e 

memberResourceGroup="${1:-myRG}"
memberCluster="${2:-myCluster}"

memberClusterId="/subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourcegroups/$memberResourceGroup/providers/Microsoft.ContainerService/managedClusters/$memberCluster"

echo "fleetID=$fleetID"
echo "memberClusterId=$memberClusterId"

curl -v -X DELETE \
 -H "x-ms-identity-url: http://msi-simulator.msi-simulator.svc.cluster.local" \
 -H "x-ms-identity-principal-id: 8ff738a5-abcd-4864-a162-6c18f7c9cbd9" \
 -H "x-ms-home-tenant-id: 72f988bf-86f1-41af-91ab-2d7cd011db47" \
 -H "Content-Type: application/json" \
 -H "x-ms-home-tenant-id: 72f988bf-86f1-41af-91ab-2d7cd011db47" \
 --data "$body" \
"http://localhost:8080$memberClusterId/providers/Microsoft.ContainerService/fleetMemberships/default?api-version=2022-06-02-preview"
