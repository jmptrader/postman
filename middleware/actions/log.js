Action.register('log', function (args) {
    var mailId = args.id;
    var c = this;
    if (!c.auth) return;
    Mail.find({where: {id: mailId}}).complete(function (err, m) {
        if (err) return;
        Log.create({
            log: args.log,
            status: args.success ? 'delivered' : 'dropped'
        }).complete(function (err, l) {
            if (err) return;
            m.addLog(l);
        });
    });
});