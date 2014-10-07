Website::App.controllers :client do
  
  get :index, :map => '/sender/list' do
    @title = 'Senders'
    render 'index', layout: 'application'
  end

end
