global.Sender = model.define('Sender', {
    ip: {type: Model.STRING, unique: true},
    domain: {type: Model.STRING, unique: true},
    status: Model.STRING,
    secret: Model.STRING,
    storage_key: Model.STRING,
    api_key: Model.STRING,
    deliver_frequency: Model.INTEGER,
    private_key: Model.TEXT,
    public_key: Model.TEXT,
    web_hook: Model.STRING,
    immediate: Model.BOOLEAN
}, {
    tableName: 'senders',
    updatedAt: 'updated_at',
    createdAt: 'created_at'
});

Sender.hasMany(Mail, {as: 'Mails', foreignKey: 'sender_id'});
Sender.hasMany(Frequency, {as: 'Frequencies', foreignKey: 'sender_id'});