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

  post :destroy, map: '/sender/delete/:id' do
    if @sender.destroy
      json code: 200
    else
      json error: "Can not remove sender #{@sender.ip}"
    end
  end

  get :dashboard, map: '/sender/:id' do
    @title = 'dashboard'
    haml :'layouts/dashboard', layout: :application do
      haml :'sender/getting_started', layout: false
    end
  end

  get :setting, map: '/sender/:id/setting' do
    @title = 'setting'
    haml :'layouts/dashboard', layout: :application do
      haml :'sender/setting', layout: false
    end
  end

  get :logs, map: '/sender/:id/logs' do
    @title = 'logs'
    haml :'layouts/dashboard', layout: :application do
      haml :'sender/logs', layout: false
    end
  end

  get :dns_records, map: '/sender/:id/dns-records' do
    @sender = Sender.get params['id']
    halt 404, 'no sender found' unless @sender
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
    return json(result)
  end

  get :config_download, map: '/sender/:id/config/download' do
    attachment 'config.json'
    {
        authSecret: @sender.secret,
        storeSecret: @sender.storage_key,
        remoteAddr: settings.middleware_addr,
        createAt: @sender.created_at.iso8601(3)
    }.to_json
  end
end
