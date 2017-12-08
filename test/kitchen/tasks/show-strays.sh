#!/usr/bin/env bash -l

printf "VMs:\n"

az vm list --query "[?starts_with(name, 'dd-agent-testing')]|[*].{name:name,location:location,state:provisioningState}" -o table

printf "\n"

printf "Groups:\n"
az group list --query "[?starts_with(name, 'kitchen-dd-agent')]|[*].{name:name,location:location,state:properties.provisioningState}" -o table
