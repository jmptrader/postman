Action.register('frequency', function (args, id) {
    var domain = args.domain;
    var c = this;
    if (!c.auth) return;
    if (domain === 'default') {
        return c.response(id, String(c.sender.deliver_frequency))
    }
    Frequency.find({where: {domain: domain}}).complete(function (err, f) {
        if (!err) {
            return c.response(id, String(f.deliver_frequency))
        }
    });
});