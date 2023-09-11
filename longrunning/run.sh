#!/bin/bash

namespace=hello-world
KUBECONFIG=./.kube/fleet.kubeconfig ./kubectl create namespace $namespace
KUBECONFIG=./.kube/fleet.kubeconfig ./kubectl apply -f ./crp.yaml

deployments=("./deployment-1.yaml" "./deployment-2.yaml")

count=0
while true
do
  echo $count
  n=$(($count % 2))
  deployment=${deployments[$n]}
  echo "apply $deployment on hub"
  KUBECONFIG=./.kube/fleet.kubeconfig ./kubectl apply -f $deployment -n $namespace
  sleep 30
  count=$((count+1))
done