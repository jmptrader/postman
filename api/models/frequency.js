global.Frequency = model.define('Frequency', {
    domain: {type: Model.STRING, unique: 'senderDomainUnique'},
    sender_id: {type: Model.INTEGER, unique: 'senderDomainUnique'},
    deliver_frequency: Model.INTEGER
}, {
    tableName: 'frequencies',
    updatedAt: 'updated_at',
    createdAt: 'created_at'
});
