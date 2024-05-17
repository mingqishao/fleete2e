#!/bin/bash

subscription="${1}"
resourceGroup="${2}"
fleet="${3}"

body=$(cat <<-END
{
    "location": "westus",
    "properties": {
        "hubProfile": {
            "dnsPrefix": "$fleet",
            "agentProfile": {
                "vmSize": "badsize"
            }
        }
    }
}
END
)

echo $body
# --headers "fleet-configuration-version=-2" \

az rest --method PUT \
    --url "https://management.azure.com/subscriptions/$subscription/resourceGroups/$resourceGroup/providers/Microsoft.ContainerService/fleets/$fleet?api-version=2023-08-15-preview" \
    --body "$body"
