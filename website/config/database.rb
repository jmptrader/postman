DataMapper.logger = logger
DataMapper::Property::String.length(255)

DataMapper.setup(:default, JSON.parse(IO.read(Padrino.root('../config/database.json')))['mysql'])
