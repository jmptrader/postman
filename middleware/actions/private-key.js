Action.register('privateKey', function (args, id) {
    var c = this;
    if (!c.auth) return;
    return c.response(id, c.sender.private_key)
});