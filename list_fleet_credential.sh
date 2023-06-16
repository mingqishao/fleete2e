#!/bin/bash

resourceGroup="${1:-myResourceGroup}"
fleet="${2:-myFleet}"

curl -v -X POST \
 -H "x-ms-home-tenant-id: 72f988bf-86f1-41af-91ab-2d7cd011db47" \
http://localhost:8080/subscriptions/18153b17-4e27-4b58-863e-f8105b8892a2/resourceGroups/$resourceGroup/providers/Microsoft.ContainerService/fleets/$fleet/listCredentials?api-version=2022-06-02-preview
