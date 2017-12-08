default['datadog']['api_key'] = nil
default['dd-agent-upgrade']['version'] = nil # => install the latest available version
default['dd-agent-upgrade']['add_new_repo'] = false # If set to true, be sure to set aptrepo and yumrepo
default['dd-agent-upgrade']['aptrepo'] =  nil
default['dd-agent-upgrade']['yumrepo'] = nil
default['dd-agent-upgrade']['package_name'] = 'datadog-agent'
