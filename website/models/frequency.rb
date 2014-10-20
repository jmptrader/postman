class Frequency
  include DataMapper::Resource

  property :id, Serial
  property :domain, String, format: /^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+.?$/
  property :deliver_frequency, Integer, format: /^([0-9]|[1-5][0-9]|60)$/,
           messages: {
               format: 'frequency should not more than 60'
           }

  property :updated_at, DateTime
  property :created_at, DateTime

  belongs_to :sender, key: true
end