#!/bin/bash

# Unless explicitly stated otherwise all files in this repository are licensed
# under the Apache License Version 2.0.
# This product includes software developed at Datadog (https://www.datadoghq.com/).
# Copyright 2017 Datadog, Inc.


##### Core config #####

if [[ -z "${DD_API_KEY}" ]]; then
    echo "You must set an DD_API_KEY environment variable to run the Datadog Agent container" >&2
    exit 1
fi

if [[ -z "${DD_DD_URL}" ]]; then
    export DD_DD_URL="https://app.datadoghq.com"
fi

if [[ -z "${DD_DOGSTATSD_SOCKET}" ]]; then
    export DD_DOGSTATSD_NON_LOCAL_TRAFFIC=1
elif [[ -e "${DD_DOGSTATSD_SOCKET}" ]]; then
    if [[ -S "${DD_DOGSTATSD_SOCKET}" ]]; then
        echo "Deleting existing socket at ${DD_DOGSTATSD_SOCKET}"
        rm -v "${DD_DOGSTATSD_SOCKET}" || exit $?
    else
        echo "${DD_DOGSTATSD_SOCKET} exists and is not a socket, please check your volume options" >&2
        ls -l "${DD_DOGSTATSD_SOCKET}" >&2
        exit 1
    fi
fi

if [[ "${KUBERNETES_SERVICE_PORT}" ]]; then
    export KUBERNETES="yes"
fi

# Install default datadog.yaml
if [[ "${KUBERNETES}" ]]; then
    ln -s /etc/datadog-agent/datadog-kubernetes.yaml /etc/datadog-agent/datadog.yaml
else
    ln -s /etc/datadog-agent/datadog-docker.yaml /etc/datadog-agent/datadog.yaml
fi

# Copy custom confs

find /conf.d -name '*.yaml' -exec cp --parents {} /etc/datadog-agent/ \;

find /checks.d -name '*.py' -exec cp --parents {} /etc/datadog-agent/ \;

# Get all of the datadog configuration files.
export TOOLS_PREFIX=/usr/local/bin
export INTEGRATION_DIR=/etc/datadog-agent/conf.d
if [[ ${CONSUL_PREFIX} && ${ENABLE_INTEGRATIONS} ]]; then
  [[ ${CONSUL_DC} ]] && CONSUL_KV_DC="-datacenter=${CONSUL_DC}" || CONSUL_KV_DC=""
  [[ ${CONSUL_ADDR_FROM_AWS_META} && $(curl -s http://169.254.169.254/latest/meta-data/local-ipv4) ]] \
  && export CONSUL_HTTP_ADDR=$(curl -s http://169.254.169.254/latest/meta-data/local-ipv4):8500
  if ${TOOLS_PREFIX}/consul kv get ${CONSUL_KV_DC} -keys "${CONSUL_PREFIX}"/integrations/ &>/dev/null; then
    for integration in $(${TOOLS_PREFIX}/consul kv get -keys "${CONSUL_PREFIX}"/integrations/); do
      THIS_INTEGRATION=$(echo "${integration}" | awk -F '/' '{print $NF}')
      if [[ ! -z "${THIS_INTEGRATION}" ]]; then
        if [[ ${ENABLE_INTEGRATIONS} == *"${THIS_INTEGRATION}"* ]]; then
            ${TOOLS_PREFIX}/consul kv get ${CONSUL_KV_DC} "${integration}" > "${INTEGRATION_DIR}"/"${THIS_INTEGRATION}".yaml
            ${TOOLS_PREFIX}/consul-template -template "${INTEGRATION_DIR}"/"${THIS_INTEGRATION}".yaml:"${INTEGRATION_DIR}"/"${THIS_INTEGRATION}".yaml -once
        fi
      fi
    done
  fi
fi

# Enable logs
if [[ $DD_ENABLE_LOGS ]]; then 
  echo -e "\n# Logs\nlog_enabled: true" | tee -a /etc/datadog-agent/datadog.yaml 
fi 

##### Starting up #####

exec "$@"
