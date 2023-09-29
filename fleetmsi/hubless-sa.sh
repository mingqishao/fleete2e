#!/bin/bash

RG="${1}"
FLEET="${2}"

parameters=$(cat <<-END
{
  "name": {
    "value": "$FLEET"
  }
}
END
)
echo $parameters
az deployment group create -g $RG --template-file  ./hubless-sa.json  --name $FLEET --parameters "${parameters}"
