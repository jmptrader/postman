global.Sender = model.define('Sender', {
    ip: {
        type: Model.STRING,
        validate: {
            isIP: true
        }
    },
    private_key: Model.STRING,
    web_hook: Model.STRING,
    secret: Model.STRING,
    status: Model.STRING
}, {
    tableName: 'senders',
    updatedAt: 'updated_at',
    createdAt: 'created_at'
});