Action.register('log', function (args) {
    var mailId = args.id;
    var c = this;
    if (!c.auth) return;
    Mail.find({where: {id: mailId}}).complete(function (err, m) {
        if (err) return;
        var result = {
            log: args.log,
            status: args.success ? 'delivered' : 'dropped'
        };
        Log.create(result).complete(function (err, l) {
            if (err) return;
            MailSync.emit(mailId, result);
            m.addLog(l);
        });
    });
});