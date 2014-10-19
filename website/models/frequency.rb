class Frequency
  include DataMapper::Resource

  property :id, Serial
  property :domain, String
  property :deliver_frequency, Integer

  property :updated_at, DateTime
  property :created_at, DateTime

  belongs_to :sender, key: true
end