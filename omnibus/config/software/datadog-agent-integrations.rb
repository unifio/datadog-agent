# Unless explicitly stated otherwise all files in this repository are licensed
# under the Apache License Version 2.0.
# This product includes software developed at Datadog (https:#www.datadoghq.com/).
# Copyright 2017 Datadog, Inc.

require './lib/ostools.rb'

name 'datadog-agent-integrations'

dependency 'pip'

relative_path 'integrations-core'
whitelist_file "embedded/lib/python2.7"

source git: 'https://github.com/DataDog/integrations-core.git'

integrations_core_branch = ENV['INTEGRATIONS_CORE_BRANCH']
if integrations_core_branch.nil? || integrations_core_branch.empty?
  integrations_core_branch = 'master'
end
default_version integrations_core_branch

integrations_core_wheels = ENV['INTEGRATIONS_CORE_WHEELS']
if integrations_core_wheels.nil? || integrations_core_wheels.empty?
  integrations_core_wheels = integrations_core_wheels.downcase != 'false'
else
  integrations_core_wheels = false
end

blacklist = [
  'datadog-base',  # namespacing package for wheels (NOT AN INTEGRATION)
  'agent_metrics',
  'docker_daemon',
  'kubernetes',
  'ntp',  # provided as a go check by the core agent
  'vsphere',
]

build do
  # The checks
  checks_dir = "#{install_dir}/agent/checks.d"
  mkdir checks_dir

  # The confs
  conf_dir = "#{install_dir}/etc/datadog-agent/conf.d"
  mkdir conf_dir

  # Copy the checks and generate the global requirements file
  block do
    if integrations_core_wheels
      if windows?
        build_args = "wheel --no-deps ."
        install_args = "install *.whl"
        command "#{windows_safe_path(install_dir)}\\embedded\\scripts\\pip.exe #{build_args}", :cwd => "#{project_dir}\\datadog-base"
        command "#{windows_safe_path(install_dir)}\\embedded\\scripts\\pip.exe #{install_args}", :cwd => "#{project_dir}\\datadog-base"
      else
        build_env = {
          "LD_RUN_PATH" => "#{install_dir}/embedded/lib",
          "PATH" => "#{install_dir}/embedded/bin:#{ENV['PATH']}",
        }
        pip "wheel --no-deps .", :env => build_env, :cwd => "#{project_dir}/datadog-base"
        pip "install *.whl", :env => build_env, :cwd => "#{project_dir}/datadog-base"
      end
    end
    all_reqs_file = File.open("#{project_dir}/check_requirements.txt", 'w+')

    Dir.glob("#{project_dir}/*").each do |check_dir|
      check = check_dir.split('/').last

      next if blacklist.include? check

      if integrations_core_wheels
        if windows?
          command "#{windows_safe_path(install_dir)}\\embedded\\scripts\\pip.exe #{build_args}", :cwd => "#{project_dir}\\#{check_dir}"
          command "#{windows_safe_path(install_dir)}\\embedded\\scripts\\pip.exe #{install_args}", :cwd => "#{project_dir}\\#{check_dir}"
        else
          build_env = {
            "LD_RUN_PATH" => "#{install_dir}/embedded/lib",
            "PATH" => "#{install_dir}/embedded/bin:#{ENV['PATH']}",
          }
          pip "wheel --no-deps .", :env => build_env, :cwd => "#{project_dir}/#{check_dir}"
          pip "install *.whl", :env => build_env, :cwd => "#{project_dir}/#{check_dir}"
        end
      else
        # Check the manifest to be sure that this check is enabled on this system
        # or skip this iteration
        manifest_file_path = "#{check_dir}/manifest.json"

        # If there is no manifest file, then we should assume the folder does not
        # contain a working check and move onto the next
        File.exist?(manifest_file_path) || next

        manifest = JSON.parse(File.read(manifest_file_path))
        manifest['supported_os'].include?(os) || next

        # Copy the checks over
        if File.exist? "#{check_dir}/datadog/#{check_dir}/#{check}.py"
          copy "#{check_dir}/datadog/#{check_dir}/#{check}.py", "#{checks_dir}/#{check}.py"
        end

        check_conf_dir = "#{conf_dir}/#{check}.d"

        # Copy the check config to the conf directories
        if File.exist? "#{check_dir}/conf.yaml.example"
          mkdir check_conf_dir unless File.exists? (check_conf_dir)
          copy "#{check_dir}/conf.yaml.example", "#{check_conf_dir}/"
        end

        # Copy the default config, if it exists
        if File.exist? "#{check_dir}/conf.yaml.default"
          mkdir check_conf_dir unless File.exists? (check_conf_dir)
          copy "#{check_dir}/conf.yaml.default", "#{check_conf_dir}/"
        end

        # Copy the metric file, if it exists
        if File.exist? "#{check_dir}/metrics.yaml"
          mkdir check_conf_dir unless File.exists? (check_conf_dir)
          copy "#{check_dir}/metrics.yaml", "#{check_conf_dir}/"
        end

        # We don't have auto_conf on windows yet
        if os != 'windows'
          if File.exist? "#{check_dir}/auto_conf.yaml"
            mkdir check_conf_dir unless File.exists? (check_conf_dir)
            copy "#{check_dir}/auto_conf.yaml", "#{check_conf_dir}/"
          end
        end

        next unless File.exist?("#{check_dir}/requirements.txt") && !manifest['use_omnibus_reqs']

        reqs = File.open("#{check_dir}/requirements.txt", 'r').read
        reqs.each_line do |line|
          all_reqs_file.puts line if line[0] != '#'
        end
      end
    end

    # Manually add "core" dependencies that are not listed in the checks requirements
    all_reqs_file.puts "requests==2.11.1"
    all_reqs_file.puts "pympler==0.5"

    all_reqs_file.close
  end

  # Install all the requirements
  if windows?
    pip_args = "install  -r #{project_dir}/check_requirements.txt"
    command "#{windows_safe_path(install_dir)}\\embedded\\scripts\\pip.exe #{pip_args}"
  else
    build_env = {
      "LD_RUN_PATH" => "#{install_dir}/embedded/lib",
      "PATH" => "#{install_dir}/embedded/bin:#{ENV['PATH']}",
    }
    pip "install -r #{project_dir}/check_requirements.txt", :env => build_env
  end

  move "#{project_dir}/check_requirements.txt", "#{install_dir}/agent/"
end
