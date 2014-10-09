Website::App.controllers :sender do

  get :index, map: '/' do
    redirect url(:sender, :list)
  end

  get :list do
    @title = 'Senders'
    @senders = Sender.all
    render 'index', layout: 'application'
  end

  post :create do
    sender_params = params['sender'].keep_if do |k, _|
      %w(ip domain).include?(k)
    end
    sender_params.each_value { |v| v.strip! }
    sender = Sender.create(sender_params)
    if sender.errors.size == 0
      return json code: 200
    end
    json errors: sender.errors.to_a.flatten
  end
end
