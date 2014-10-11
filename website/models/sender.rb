class Sender
  include DataMapper::Resource

  property :id, Serial
  property :ip, String, required: true, unique: true,
           format: /^((1?\d?\d|(2([0-4]\d|5[0-5])))\.){3}(1?\d?\d|(2([0-4]\d|5[0-5])))$/,
           messages: {
               required: 'ip record can not be blank',
               unique: 'ip record has been used before',
               format: 'ip record is not validated'
           }
  property :domain, String, required: true, unique: true,
           format: /^([0-9a-z_-]+\.)*?[a-z0-9-]+\.[a-z]{2,6}(\.[a-z]{2})?$/,
           messages: {
               required: 'domain record can not be blank',
               unique: 'domain record has been used before',
               format: 'domain record is not validated'
           }
  property :status, String, default: 'offline'
  property :secret, String, default: lambda { |r, _| r.random_secret 32 }
  property :storage_key, String, default: lambda { |r, _| r.random_secret 32 }
  property :api_key, String, default: lambda { |r, _| r.random_secret 32 }
  property :updated_at, DateTime
  property :created_at, DateTime

  public
  def random_secret(len)
    o = [('a'..'z'), ('A'..'Z'), %w(+ / ! . ?)].map { |i| i.to_a }.flatten
    (0...len).map { o[rand(o.length)] }.join
  end
end
