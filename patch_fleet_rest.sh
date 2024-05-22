#!/bin/bash

subscription="${1}"
resourceGroup="${2}"
fleet="${3}"

body=$(cat <<-END
{
    "identity": {
        "type": "SystemAssigned"
    }
}
END
)

echo $body

    # --headers "fleet-configuration-version=" \
az rest --method Patch \
    --url "https://management.azure.com/subscriptions/$subscription/resourceGroups/$resourceGroup/providers/Microsoft.ContainerService/fleets/$fleet?api-version=2023-08-15-preview" \
    --body "$body" --debug --verbose
