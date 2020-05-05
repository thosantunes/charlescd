#!/bin/bash   
modules=( "darwin-application" "darwin-deploy" "darwin-notifications" "darwin-ui-legacy" "darwin-ui-new" "darwin-villager" "keycloak" "postgresql" "nginx")
rm -rf charles.yaml

for i in ${modules[@]}; do
  helm template $i "./charts/$i" --namespace=charles >> charles.yaml
done