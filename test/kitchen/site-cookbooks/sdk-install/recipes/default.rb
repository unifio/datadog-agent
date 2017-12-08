#
# Cookbook Name:: sdk-install
# Recipe:: default
#
# Copyright (C) 2013 Datadog
#
# All rights reserved - Do Not Redistribute
#

if node['sdk-install']['add_new_repo']
  case node['platform_family']
  when 'debian'
    include_recipe 'apt'

    apt_repository 'datadog-sdk' do
      keyserver 'keyserver.ubuntu.com'
      key 'c7a7da52'
      uri node['sdk-install']['aptrepo']
      distribution node['sdk-install']['aptrepo_dist']
      components ['main']
      action :add
    end

  when 'rhel'
    include_recipe 'yum'

    yum_repository 'datadog-sdk' do
      name 'datadog-update'
      description 'datadog-update'
      url node['sdk-install']['yumrepo']
      action :add
      make_cache true
      # Older versions of yum embed M2Crypto with SSL that doesn't support TLS1.2
      protocol = node['platform_version'].to_i < 6 ? 'http' : 'https'
      gpgkey "#{protocol}://yum.datadoghq.com/DATADOG_RPM_KEY.public"
    end
  end
end
