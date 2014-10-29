##
# This file mounts each app in the Padrino project to a specified sub-uri.
# You can mount additional applications using any of these commands below:
#
#   Padrino.mount('blog').to('/blog')
#   Padrino.mount('blog', :app_class => 'BlogApp').to('/blog')
#   Padrino.mount('blog', :app_file =>  'path/to/blog/app.rb').to('/blog')
#
# You can also map apps to a specified host:
#
#   Padrino.mount('Admin').host('admin.example.org')
#   Padrino.mount('WebSite').host(/.*\.?example.org/)
#   Padrino.mount('Foo').to('/foo').host('bar.example.org')
#
# Note 1: Mounted apps (by default) should be placed into the project root at '/app_name'.
# Note 2: If you use the host matching remember to respect the order of the rules.
#
# By default, this file mounts the primary app which was generated with this project.
# However, the mounted app can be modified as needed:
#
#   Padrino.mount('AppName', :app_file => 'path/to/file', :app_class => 'BlogApp').to('/')
#

##
# Setup global project settings for your apps. These settings are inherited by every subapp. You can
# override these settings in the subapps as needed.
#
require 'json'
require 'resolv'

Padrino.configure_apps do
  # enable :sessions
  # set :session_secret, '2fb648e3a4f628fdda8620ac11201b2e5f2f5f2a482dc3bcf751da39b5fcd648'
  set :protection, :except => :path_traversal
  set :protect_from_csrf, true

  domain = JSON.parse(IO.read(Padrino.root('../config/domain.json')))
  set :middleware_addr, "#{domain['middleware']['domain']}:#{domain['middleware']['port']}"
  set :api_addr, "#{domain['api']['protocol']}://#{domain['api']['domain']}:#{domain['api']['port']}"
end
# Mounts the core application for this project
Padrino.mount('Website::App', :app_file => Padrino.root('app/app.rb')).to('/')