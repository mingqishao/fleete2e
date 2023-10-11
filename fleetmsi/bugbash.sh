export aciName=$fleetname-aci
export aciSubnetName=aci
export aciSubnetId=$(az network vnet subnet create -g $rgname \
   --vnet-name $vNetName \
   --name $aciSubnetName \
   --address-prefixes 172.19.2.0/24 \
   --query id --output tsv)


az container create \
  --resource-group $rgname \
  --name $aciName \
  --image shaomq/fleet-cli \
  --command-line "tail -f /dev/null" \
  --restart-policy never \
  --vnet $vNetName \
  --subnet $aciSubnetName

az container exec \
  --resource-group $rgname \
  --name $aciName \
  --exec-command /bin/bash

export cluterAdmin="Azure Kubernetes Fleet Manager RBAC Cluster Admin"
export fleetId=$( az fleet show -g $rgname -n $fleetname \
   --query id --output tsv
)

export me=minsha@microsoft.com

az role assignment create --role "${cluterAdmin}" --assignee ${me} --scope ${fleetId}