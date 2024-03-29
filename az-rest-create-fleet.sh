#!/bin/bash

resourceGroup="${1:-myResourceGroup}"
fleet="${2:-myFleet}"

body=$(cat <<-END
{
    "tags": {
        "tag1": "value1",
        "tag2": "value2"
    },
    "location": "centraluseuap",
    "properties": {
        "hubProfile": {
            "dnsPrefix": "$fleet",
        }
    }
}
END
)

az rest \
--method PUT  \
--uri "/subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourceGroups/$resourceGroup/providers/Microsoft.ContainerService/fleets/$fleet?api-version=2022-06-02-preview" \
--header "fleet-configuration-version=-1" \
--body "$body"
