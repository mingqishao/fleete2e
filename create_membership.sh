
#!/bin/bash
set -e 

memberResourceGroup="${1:-myRG}"
memberCluster="${2:-myCluster}"
fleetResourceGroup="${3:-myRG}"
fleet="${4:-myFleet}"

fleetID="/subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourceGroups/$fleetResourceGroup/providers/Microsoft.ContainerService/fleets/$fleet"
memberClusterId="/subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourcegroups/$memberResourceGroup/providers/Microsoft.ContainerService/managedClusters/$memberCluster"

echo "fleetID=$fleetID"
echo "memberClusterId=$memberClusterId"

echo "get the fleet's credentials"

credentials=`curl -v -X POST \
 -H "x-ms-home-tenant-id: 72f988bf-86f1-41af-91ab-2d7cd011db47" \
http://localhost:8080/subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourceGroups/$fleetResourceGroup/providers/Microsoft.ContainerService/fleets/$fleet/listCredentials?api-version=2022-06-02-preview`

# echo "$credentials" | jq '.kubeconfigs[0].value'
kubeconfig=$( jq -r  '.kubeconfigs[0].value' <<< "${credentials}" ) 
echo $?
if [ -z "$kubeconfig" ] || [ "$kubeconfig" == "null" ];
then
      echo "failed to get fleet's credentials"
      exit 1
fi
# echo $kubeconfig | base64 -d
kubeconfig=$( base64 -d <<< $kubeconfig)
fleetHubURL=$( yq '.clusters[0].cluster.server' <<< "$kubeconfig" )
ca=$( yq '.clusters[0].cluster.certificate-authority-data' <<< "$kubeconfig" )


body=$(cat <<-END
{
    "properties":{
        "fleetResourceId":"$fleetID",
        "tenantId":"72f988bf-86f1-41af-91ab-2d7cd011db47",
        "fleetHubUrl":"$fleetHubURL",
        "certificateAuthorityData":"$ca"
    },
    "name":"default"
}
END
)
curl -v -X PUT \
 -H "x-ms-identity-url: http://msi-simulator.msi-simulator.svc.cluster.local" \
 -H "x-ms-identity-principal-id: 8ff738a5-abcd-4864-a162-6c18f7c9cbd9" \
 -H "x-ms-home-tenant-id: 72f988bf-86f1-41af-91ab-2d7cd011db47" \
 -H "Content-Type: application/json" \
 -H "x-ms-home-tenant-id: 72f988bf-86f1-41af-91ab-2d7cd011db47" \
 --data "$body" \
"http://localhost:8080/subscriptions/8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8/resourcegroups/$memberResourceGroup/providers/Microsoft.ContainerService/managedClusters/$memberCluster/providers/Microsoft.ContainerService/fleetMemberships/default?api-version=2022-06-02-preview"
