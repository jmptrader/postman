Website::App.controllers :sender do

  get :index, map: '/' do
    redirect url(:sender, :list)
  end

  get :list do
    @title = 'Senders'
    render 'index', layout: 'application'
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
      return json code: 200, sender_id: sender.id
    end
    json errors: sender.errors.to_a.flatten
  end

  post :destroy, map: '/sender/delete/:id' do
    sender = Sender.get params['id']
    return json error: 'Sender record does not exist!' unless sender
    if sender.destroy
      json code: 200
    else
      json error: "Can not remove sender #{sender.ip}"
    end
  end
end
