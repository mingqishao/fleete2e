export SUB=3959ec86-5353-4b0c-b5d7-3877122861a0
export VNET1_RG=minsha-network-cross-region
export VNET1=vnet-swedencentral
export VNET2_RG=bugbash-aks-group3
export VNET2=vnet-southcentralus
az account set -s $SUB

az network vnet peering create \
    -g ${VNET1_RG} \
    --vnet-name ${VNET1} \
    -n "${VNET1}_to_${VNET2}" \
    --remote-vnet /subscriptions/${SUB}/resourceGroups/${VNET2_RG}/providers/Microsoft.Network/virtualNetworks/${VNET2} --allow-vnet-access

az network vnet peering create \
    -g ${VNET2_RG} \
    --vnet-name ${VNET2} \
    -n "${VNET2}_to_${VNET1}" \
    --remote-vnet /subscriptions/${SUB}/resourceGroups/${VNET1_RG}/providers/Microsoft.Network/virtualNetworks/${VNET1} --allow-vnet-access