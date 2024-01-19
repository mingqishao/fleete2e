#!/bin/bash


set -x
set -e

RG="${1}"
FLEET="${2}"
az group create -n $RG -l westus2 
fleetId=$(az fleet create -g $RG -n $FLEET --enable-hub \
    --query id --output tsv)

logAnalyticsId=$(az monitor log-analytics workspace create -g $RG -n ${FLEET}-logs \
    --query id --output tsv)

eventhubId=$(az eventhubs namespace create -g $RG -n ${FLEET}-eventhub \
    --query id --output tsv)

# --event-hub /subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourceGroups/minsha-ln/providers/Microsoft.EventHub/namespaces/fleet-1-eventhub \
# /subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourcegroups/minsha-ln/providers/Microsoft.EventHub/namespaces/fleet-1-eventhub/authorizationrules/RootManageSharedAccessKey
az monitor diagnostic-settings create -n ${FLEET}-diag \
    --resource $fleetId \
    --event-hub-rule ${eventhubId}/authorizationrules/RootManageSharedAccessKey  \
    --workspace ${logAnalyticsId} \
    --logs @logs.json
    # --debug

# az rest --method GET -u https://management.azure.com/subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourceGroups/minsha-ln/providers/Microsoft.ContainerService/fleets/fleet-1/providers/Microsoft.Insights/diagnosticSettings/test?api-version=2021-05-01-preview
