#!/bin/bash

subscription="${1}"
resourceGroup="${2}"
cluster="${3}"

body=$(cat <<-END
{
    "location": "eastus",
    "properties": {
    }
}
END
)

echo $body

az rest --method PUT \
    --url "https://management.azure.com/subscriptions/$subscription/resourceGroups/$resourceGroup/providers/Microsoft.ContainerService/managedClusters/$cluster?api-version=2022-04-01" \
    --body "$body"
