class SendException
  include DataMapper::Resource

  storage_names[:default] = 'exceptions'

  property :id, Serial
  property :content, String, required: true, unique: true
  property :is_new, Boolean
  property :treatment, String

  property :updated_at, DateTime
  property :created_at, DateTime

end