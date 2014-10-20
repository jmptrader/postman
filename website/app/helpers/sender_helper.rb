# Helper methods defined here can be accessed in any controller or view in the application

module Website
  class App
    module SenderHelper

      attr_accessor :dns_obj, :redis
      @@dns_obj = Resolv::DNS.new(:nameserver => %w(8.8.8.8 8.8.4.4))
      @@redis = Redis.new(JSON.parse(IO.read(Padrino.root('../config/database.json')))['redis'])

      module_function

      def get_txt_record(domain_name)
        records = @@dns_obj.getresources(domain_name, Resolv::DNS::Resource::IN::TXT)
        records.map! { |rd| rd.data }
        records.join('\n')
      end

      def check_txt_record(domain_name, record)
        records = @@dns_obj.getresources(domain_name, Resolv::DNS::Resource::IN::TXT)
        result = false
        records.each do |rd|
          result = true if rd.data.to_s == record
        end
        result
      end

      def send_cmd(sender_id, command)
        @@redis.lpush('jianxin:command', "#{sender_id}:#{command.to_json}")
      end

    end

    helpers SenderHelper
  end
end
