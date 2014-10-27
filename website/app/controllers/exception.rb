Website::App.controllers :exception do

  get :index, map: '/exception' do
    render 'index', layout: 'application'
  end

end
