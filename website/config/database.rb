DataMapper.logger = logger
DataMapper::Property::String.length(255)

DataMapper.setup(:default, YAML.load_file(Padrino.root('config/database.yml'))[RACK_ENV])
