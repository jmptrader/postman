global.Client = model.define('Client', {
    ip: {
        type: Model.STRING,
        validate: {
            isIP: true
        }
    },
    secret: Model.STRING,
    status: Model.STRING
});