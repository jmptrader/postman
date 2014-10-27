Website::App.controllers :exception do

  get :index, map: '/exception' do
    @exceptions = params['type'] == 'all' ? SendException.all : SendException.all(is_new: true)
    render 'index', layout: 'application'
  end

  post :update, map: '/exception/:id/update' do
    e = SendException.get params['id']
    halt 404, 'no such exception found' unless e
    treatment = params['exception']['treatment']
    halt 405 unless %w(ignore resendNow resendLater).include? treatment
    if e.update treatment: treatment, is_new: false
      return json code: 200
    end
    json errors: e.errors.to_a.flatten
  end
end
