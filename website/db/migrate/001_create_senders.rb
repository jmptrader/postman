migration 1, :create_senders do
  up do
    create_table :senders do
      column :id, Integer, :serial => true
      column :ip, DataMapper::Property::String, :length => 255
      column :domain, DataMapper::Property::String, :length => 255
      column :status, DataMapper::Property::String, :length => 255
      column :secret, DataMapper::Property::String, :length => 255
      column :storage_key, DataMapper::Property::String, :length => 255
      column :api_key, DataMapper::Property::String, :length => 255
      column :private_key, DataMapper::Property::Text
      column :public_key, DataMapper::Property::Text
      column :web_hook, DataMapper::Property::String, :length => 255
      column :immediate, DataMapper::Property::Boolean
      column :updated_at, DataMapper::Property::DateTime
      column :created_at, DataMapper::Property::DateTime
    end
  end

  down do
    drop_table :senders
  end
end
