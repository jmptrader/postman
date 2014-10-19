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
  property :status, String, default: 'unverified'
  property :deliver_frequency, Integer, default: 10

  # secret keys
  property :secret, String
  property :storage_key, String
  property :api_key, String
  property :private_key, Text
  property :public_key, Text
  # settings
  property :web_hook, String
  property :immediate, Boolean, default: false

  property :updated_at, DateTime
  property :created_at, DateTime

  has n, :frequencies, child_key: [:sender_id]

  private
  def random_secret(len)
    o = [('a'..'z'), ('A'..'Z'), %w(+ / ! . ?)].map { |i| i.to_a }.flatten
    (0...len).map { o[rand(o.length)] }.join
  end

  def generate_keys
    key = OpenSSL::PKey::RSA.new 1024
    self.private_key = key.to_pem
    self.public_key = key.public_key.to_pem
  end

  public
  def init
    self.secret = random_secret(32)
    self.storage_key = random_secret(32)
    self.api_key = random_secret(32)
    generate_keys
    save!
  end

  def dkim_record
    public_key = self.public_key.lines.map { |line| line.chomp }
    "v=DKIM1; k=rsa; p=#{public_key[1...-1].join}"
  end

end
