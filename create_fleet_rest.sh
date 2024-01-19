#!/bin/bash

subscription="${1}"
resourceGroup="${2}"
fleet="${3}"

body=$(cat <<-END
{
    "location": "westus2",
    "properties": {
        "hubProfile": {
            "dnsPrefix": "$fleet"
        }
    }
}
END
)

echo $body

az rest --method PUT \
    --headers "fleet-configuration-version=-2" \
    --url "https://management.azure.com/subscriptions/$subscription/resourceGroups/$resourceGroup/providers/Microsoft.ContainerService/fleets/$fleet?api-version=2023-08-15-preview" \
    --body "$body"
