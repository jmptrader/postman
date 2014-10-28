Website::App.controllers :sender do

  before do
    if params['id']
      @sender = Sender.get params['id']
      halt 404, 'no such sender found' unless @sender
    end
  end

  get :index, map: '/' do
    redirect url(:sender, :list)
  end

  get :list do
    @title = 'Senders'
    render 'list', layout: 'application'
  end

  get :list_part, map: '/sender/list.html' do
    if  %w(offline online).include? params['status']
      @senders = Sender.all status: params['status'], order: [:id.desc]
    else
      @senders = Sender.all order: [:id.desc]
    end
    @sender_size = Sender.count
    render '_sender_list', layout: false
  end

  post :create do
    sender_params = params['sender'].keep_if do |k, _|
      %w(ip domain).include?(k)
    end
    sender_params.each_value { |v| v.strip! }
    sender = Sender.create(sender_params)
    if sender.errors.size == 0
      sender.init
      return json code: 200, sender_id: sender.id
    end
    json errors: sender.errors.to_a.flatten
  end

  post :destroy, map: '/sender/:id/delete' do
    if @sender.destroy
      json code: 200
    else
      json error: "Can not remove sender #{@sender.ip}"
    end
  end

  get :dashboard, map: '/sender/:id' do
    @title = 'dashboard'
    @api_addr = "http://#{settings.api_addr}"
    haml :'layouts/dashboard', layout: :application do
      haml :'sender/index', layout: false
    end
  end

  post :update, map: '/sender/:id/update' do
    sender_params = params['sender'].keep_if do |k, _|
      %w(immediate deliver_frequency web_hook).include?(k)
    end
    if @sender.update sender_params
      if sender_params.has_key? 'deliver_frequency'
        send_cmd @sender.id, command: 'frequency', action: 'update',
                 domain: 'default', value: sender_params['deliver_frequency'].to_s
      end
      json code: 200
    else
      json errors: @sender.errors.to_a.flatten
    end
  end

  get :dns_records, map: '/sender/:id/dns-records' do
    return json error: 'dns records has been verified' unless @sender.status == 'unverified'
    json({
             code: 200,
             spf: get_txt_record(@sender.domain),
             dkim: get_txt_record("mx._domainkey.#{@sender.domain}")
         })
  end

  post :check_dns, map: '/sender/:id/check-dns' do
    result = {
        code: 200,
        spf: get_txt_record(@sender.domain).include?(@sender.ip),
        dkim: get_txt_record("mx._domainkey.#{@sender.domain}") == @sender.dkim_record
    }
    if result[:spf] && result[:dkim]
      @sender.status = 'offline'
      @sender.save!
    end
    result[:records] = {
        spf: get_txt_record(@sender.domain),
        dkim: get_txt_record("mx._domainkey.#{@sender.domain}")
    }
    return json(result)
  end

  get :setting, map: '/sender/:id/setting' do
    @title = 'setting'
    haml :'layouts/dashboard', layout: :application do
      haml :'sender/setting', layout: false
    end
  end

  get :frequencies, map: '/sender/:id/frequencies.html' do
    @frequencies = @sender.frequencies.all order: [:id.desc]
    render '_frequencies', layout: false
  end

  post :frequency_create, map: '/sender/:id/frequency/create' do
    frequency_params = params['frequency'].keep_if do |k, _|
      %w(domain deliver_frequency).include?(k)
    end
    frequency = @sender.frequencies.create frequency_params
    if frequency.errors.size == 0
      send_cmd @sender.id, command: 'frequency', action: 'update',
               domain: frequency.domain, value: frequency.deliver_frequency.to_s
      return json code: 200
    end
    json errors: frequency.errors.to_a.flatten
  end

  post :frequency_delete, map: '/sender/:id/frequency/:frequency_id/delete' do
    frequency = @sender.frequencies.first id: params['frequency_id']
    if frequency.destroy
      send_cmd @sender.id, command: 'frequency', action: 'delete',
               domain: frequency.domain
      json code: 200
    else
      json error: "Can not remove role for #{frequency.domain}"
    end
  end

  get :logs, map: '/sender/:id/logs' do
    @title = 'logs'
    @api_addr = "http://#{settings.api_addr}"
    haml :'layouts/dashboard', layout: :application do
      haml :'sender/logs', layout: false
    end
  end

  get :config_download, map: '/sender/:id/config/download' do
    attachment 'config.json'
    {
        authSecret: @sender.secret,
        storeSecret: @sender.storage_key,
        remoteAddr: settings.middleware_addr,
        createAt: @sender.created_at.iso8601(3),
        hostname: @sender.domain
    }.to_json
  end
end
