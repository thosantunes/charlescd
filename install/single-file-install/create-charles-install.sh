#!/bin/bash   
modules=( "charles-moove" "charles-deploy" "charles-notifications" "charles-ui-legacy" "charles-ui-new" "charles-villager" "keycloak" "postgresql" "nginx" "redis" "charles-circle-matcher")
rm -rf charles.yaml

for i in ${modules[@]}; do
echo $i
  helm template $i "./charts/$i" --namespace=charles >> charles.yaml
done