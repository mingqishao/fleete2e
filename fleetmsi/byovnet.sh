#!/bin/bash

set -x

RG="${1}"
FLEET="${2}"

LOCATION=eastus
VNET=$FLEET-vnet
API_SUBNET=api
NODE_SUBNET=node
MSI=msi
MEMBER=$FLEET-member

az network vnet create -n ${VNET} \
-g $RG \
--address-prefixes 172.19.0.0/16 \
-l $LOCATION

API_SUBNET_ID=$(az network vnet subnet create -g $RG \
--vnet-name $VNET \
--name $API_SUBNET \
--delegations Microsoft.ContainerService/managedClusters \
--address-prefixes 172.19.0.0/28 \
--query id --output tsv)

NODE_SUBNET_ID=$(az network vnet subnet create -g $RG \
--vnet-name $VNET \
--name $NODE_SUBNET \
--address-prefixes 172.19.1.0/24 \
--query id --output tsv)

MSI_ID=$(az identity create -g $RG -n $MSI \
--query id --output tsv)

MSI_SP_ID=$(az identity show \
  --resource-group $RG \
  --name $MSI \
  --query principalId --output tsv)

sleep 10 

az role assignment create --scope $API_SUBNET_ID \
--role "Network Contributor" \
--assignee $MSI_SP_ID


az role assignment create --scope $NODE_SUBNET_ID \
--role "Network Contributor" \
--assignee $MSI_SP_ID

echo API_SUBNET_ID: $API_SUBNET_ID
echo NODE_SUBNET_ID: $NODE_SUBNET_ID
echo MSI_ID: $MSI_ID

parameters=$(cat <<-END
{
  "name": {
    "value": "$FLEET"
  },
  "apisubnet": {
    "value": "$API_SUBNET_ID"
  },
  "nodesubnet": {
    "value": "$NODE_SUBNET_ID"
  },
  "msi": {
    "value": "$MSI_ID"
  },
}
END
)


echo "creating fleet..."
az deployment group create -g $RG --template-file  ./byovnet.json  --name $FLEET --parameters "${parameters}"


# echo "creating member..."
# az aks create -g ${RG} -n ${MEMBER} --generate-ssh-keys --network-plugin azure --node-count 1

# MEMBER_ID=$(az aks show -g ${RG} -n ${MEMBER}  \
#     --query id --output tsv)

# echo "join member..."
# az fleet member create -g $RG --fleet-name $FLEET --member-cluster-id=${MEMBER_ID} -n ${MEMBER}
