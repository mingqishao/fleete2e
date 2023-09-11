#!/bin/bash

resourceGroup="${1:-myResourceGroup}"
fleet="${2:-myFleet}"

body=$(cat <<-END
{
    "tags": {
        "newTag1": "newValue1"
    },
    "identity": {
        "type":"UserAssigned",
        "userAssignedIdentities": {
            "/subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourcegroups/minshae2e/providers/Microsoft.ManagedIdentity/userAssignedIdentities/msi": null,
            "/subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourcegroups/minshae2e/providers/Microsoft.ManagedIdentity/userAssignedIdentities/msi-2": {}
        }
    }
}
END
)

curl -v  -X PATCH \
 -H "Content-Type: application/json" \
 -H "x-ms-home-tenant-id: 72f988bf-86f1-41af-91ab-2d7cd011db47" \
 -H "x-ms-home-tenant-id: 72f988bf-86f1-41af-91ab-2d7cd011db47" \
 -H "x-ms-identity-url: http://msi-simulator.msi-simulator.svc.cluster.local" \
 --data "$body" \
 "http://localhost:8080/subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourceGroups/$resourceGroup/providers/Microsoft.ContainerService/fleets/$fleet?api-version=2023-06-15-preview"
