global.Sender = model.define('Sender', {
    domain: Model.STRING,
    api_key: Model.STRING,
    status: Model.STRING,
    immediate: Model.BOOLEAN
}, {
    tableName: 'senders',
    updatedAt: 'updated_at',
    createdAt: 'created_at'
});