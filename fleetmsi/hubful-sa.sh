#!/bin/bash

RG="${1}"
FLEET="${2}"
PRIVATE="${3}"


MEMBER=$FLEET-member

parameters=$(cat <<-END
{
  "name": {
    "value": "$FLEET"
  },
  "private": {
    "value": $PRIVATE
  }
}
END
)
az deployment group create -g $RG --template-file  ./hubful-sa.json  --name $FLEET --parameters "${parameters}"

echo "creating member..."
az aks create -g ${RG} -n ${MEMBER} --generate-ssh-keys --network-plugin azure --node-count 1

MEMBER_ID=$(az aks show -g ${RG} -n ${MEMBER}  \
    --query id --output tsv)

echo "join member..."
az fleet member create -g $RG --fleet-name $FLEET --member-cluster-id=${MEMBER_ID} -n ${MEMBER}
