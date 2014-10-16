# Helper methods defined here can be accessed in any controller or view in the application

module Website
  class App

    module SenderHelper
      module_function

      def get_txt_record(domain_name)
        @dns_obj = Resolv::DNS.new(:nameserver => %w(8.8.8.8 8.8.4.4))
        records = @dns_obj.getresources(domain_name, Resolv::DNS::Resource::IN::TXT)
        records.map! { |rd| rd.data }
        records.join('\n')
      end

      def check_txt_record(domain_name, record)
        @dns_obj = Resolv::DNS.new(:nameserver => %w(8.8.8.8 8.8.4.4))
        records = @dns_obj.getresources(domain_name, Resolv::DNS::Resource::IN::TXT)
        result = false
        records.each do |rd|
          result = true if rd.data.to_s == record
        end
        result
      end

      def include_txt_record?(domain_name, record)
        records = @dns_obj.getresources(domain_name, Resolv::DNS::Resource::IN::TXT)
        result = false
        records.each do |rd|
          result = true if rd.data.to_s.include? record
        end
        result
      end
    end

    helpers SenderHelper
  end
end
