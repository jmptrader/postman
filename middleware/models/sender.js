global.Sender = model.define('Sender', {
    ip: {
        type: Model.STRING,
        validate: {
            isIP: true
        }
    },
    secret: Model.STRING,
    status: Model.STRING
}, {
    tableName: 'senders',
    updatedAt: 'updated_at',
    createdAt: 'created_at'
});