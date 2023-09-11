#!/bin/bash

resourceGroup="${1:-myResourceGroup}"
fleet="${2:-myFleet}"
subscription="${3:-8ecadfc9-d1a3-4ea4-b844-0d9f87e4d7c8}"

set -x 
curl -X GET \
http://localhost:8080/subscriptions/$subscription/resourceGroups/$resourceGroup/providers/Microsoft.ContainerService/fleets/$fleet?api-version=2023-06-15-preview
