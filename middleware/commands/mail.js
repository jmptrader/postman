Command.register('sendMail', function (cmd) {
    var c = this;
    Mail.find({
        where: {id: cmd.mailId}
    }).complete(function (err, mail) {
        mail.read(function (data) {
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