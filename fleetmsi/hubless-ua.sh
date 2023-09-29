#!/bin/bash

RG="${1}"
FLEET="${2}"
MSI=$FLEET-msi
MSI_ID=$(az identity create -g $RG -n $MSI \
  --query id --output tsv)

echo "MSI ID: $MSI_ID"
parameters=$(cat <<-END
{
  "name": {
    "value": "$FLEET"
  },
  "msi": {
    "value": "$MSI_ID"
  }
}
END
)
az deployment group create -g $RG --template-file  ./hubless-ua.json  --name $FLEET --parameters "${parameters}"
