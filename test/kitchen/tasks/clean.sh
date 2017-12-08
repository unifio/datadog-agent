#!/usr/bin/env bash -l

vms=($(az vm list --query "[?starts_with(name, 'dd-agent-testing')]|[*].name" --output tsv))

for vm in $vms; do
  echo "az vm delete -n $vm -y"
  az group delete -n $vm -y &
  echo "\n\n"
done

groups=($(az group list -o tsv --query "[?starts_with(name, 'kitchen-dd-agent')]|[*].name"))

for group in $groups; do
  echo "az group delete -n $group -y"
  az group delete -n $group -y &
  echo "\n\n"
done
