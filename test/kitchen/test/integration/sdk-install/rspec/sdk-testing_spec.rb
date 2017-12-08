require 'json'
require 'spec_helper'

describe 'the installed integration' do
  let(:expected_installed_packages) {
    JSON.parse(IO.read("/tmp/kitchen/dna.json"))
      .fetch('integration_rspec')
      .fetch('packages')
  }

  it 'has the name of the check in the info command' do
    expected_installed_packages.each do |check|
      expect(info).to include "#{check}"
    end
  end

  it 'has the integration packages' do
    expected_installed_packages.each do |check|
      file_name = "/opt/datadog-agent/integrations/#{check}/check.py"
      expect(File).to exist(file_name)
    end
  end

  it 'has the integration config' do
    expected_installed_packages.each do |check|
      example_conf_file_name = "/etc/dd-agent/conf.d/examples/#{check}.yaml.example"
      conf_file_name = "/etc/dd-agent/conf.d/#{check}.yaml"
      expect(File).to exist(example_conf_file_name)
      expect(File).to exist(conf_file_name)
    end
  end

  it 'should remove the integration' do
    if system('which apt-get &> /dev/null')
      expect(system("sudo apt-get -q -y remove dd-check-#{check} > /dev/null")).to be_truthy
    elsif system('which yum &> /dev/null')
      expect(system("sudo yum -y remove dd-check-#{check} > /dev/null")).to be_truthy
    else
      raise 'Unknown package manager'
    end
  end

end
