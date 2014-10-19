Command.register('sendMail', function (cmd) {
    var c = this;
    Mail.read(cmd.mailId, function (data) {
        Mail.find({
            where: {id: cmd.mailId}
        }).complete(function (err, mail) {
            if (err) return;
            c.command('sendMail', {
                id: String(mail.id),
                content: data,
                from: mail.from,
                to: mail.to,
                web_hook: c.sender.web_hook,
                immediate: mail.immediate
            })
        });
    });
});