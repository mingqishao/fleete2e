#!/bin/bash

resourceGroup="${1:-myResourceGroup}"
fleet="${2:-myFleet}"

body=$(cat <<-END
{
    "tags": {
        "newTag1": "newValue1"
    }
}
END
)

curl -v  -X PATCH \
 -H "Content-Type: application/json" \
 -H "x-ms-home-tenant-id: 72f988bf-86f1-41af-91ab-2d7cd011db47" \
 -H "fleet-configuration-version;" \
 --data "$body" \
 "http://localhost:8080/subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourceGroups/$resourceGroup/providers/Microsoft.ContainerService/fleets/$fleet?api-version=2023-06-15-preview"
