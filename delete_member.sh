!/bin/bash

memberResourceGroup="${1:-myRG}"
memberCluster="${2:-myCluster}"
fleetResourceGroup="${3:-myRG}"
fleet="${4:-myFleet}"

curl -v  -X DELETE \
 -H "Content-Type: application/json" \
 -H "x-ms-home-tenant-id: 72f988bf-86f1-41af-91ab-2d7cd011db47" \
 --data "$body" \
 "http://localhost:8080/subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourceGroups/$fleetResourceGroup/providers/Microsoft.ContainerService/fleets/$fleet/members/$memberCluster?api-version=2022-06-02-preview"
