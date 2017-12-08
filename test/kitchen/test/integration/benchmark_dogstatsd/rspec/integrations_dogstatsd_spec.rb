require 'spec_helper'

# Maximum memory usage allowed in kilobytes
MAX_MEMORY_FOOTPRINT = 20000

describe 'integrations-dogstatsd' do
  include_examples 'a running Agent with no errors'

  context "After dogstatd has been running five minutes" do
    require 'statsd'
    before(:each) do
      restart
      statsd = Statsd.new('localhost', 8125)
      puts 'Starting to feed dogstatsd for 300s'
      300.times do
        (1..100).each do |i|
          statsd.histogram('dd-agent.testing.bench_'<< i.to_s, Random.rand(10))
        end
        sleep 1
      end
    end

    it "doesn't consume too much memory" do
      expect(dogstatsd_memory).to be < MAX_MEMORY_FOOTPRINT
      puts "Memory used: "<< dogstatsd_memory.to_s << "kb"
    end
  end

end

def dogstatsd_memory
  `ps -u dd-agent -o rss,cmd|grep dogstatsd`.split(pattern=' ').first.to_i
end

