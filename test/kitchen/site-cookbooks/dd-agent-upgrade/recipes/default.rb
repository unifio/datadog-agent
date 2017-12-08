#
# Cookbook Name:: dd-agent-upgrade
# Recipe:: default
#
# Copyright (C) 2013 Datadog
#
# All rights reserved - Do Not Redistribute
#

if node['dd-agent-upgrade']['add_new_repo']
  case node['platform_family']
  when 'debian'
    include_recipe 'apt'

    apt_repository 'datadog-update' do
      keyserver 'keyserver.ubuntu.com'
      key 'c7a7da52'
      uri node['dd-agent-upgrade']['aptrepo']
      distribution node['dd-agent-upgrade']['aptrepo_dist']
      components ['main']
      action :add
    end

  when 'rhel'
    include_recipe 'yum'

    yum_repository 'datadog-update' do
      name 'datadog-update'
      description 'datadog-update'
      url node['dd-agent-upgrade']['yumrepo']
      action :add
      make_cache true
      # Older versions of yum embed M2Crypto with SSL that doesn't support TLS1.2
      protocol = node['platform_version'].to_i < 6 ? 'http' : 'https'
      gpgkey "#{protocol}://yum.datadoghq.com/DATADOG_RPM_KEY.public"
    end
  end
end

package node['dd-agent-upgrade']['package_name'] do
  action :upgrade
  version node['dd-agent-upgrade']['version']
end
