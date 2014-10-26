global.Exception = model.define('Exception', {
    content: {type: Model.STRING, unique: true},
    is_new: {type: Model.BOOLEAN, defaultValue: true},
    treatment: Model.STRING
}, {
    tableName: 'exceptions',
    updatedAt: 'updated_at',
    createdAt: 'created_at'
});
