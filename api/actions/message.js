var BuildMail = require('buildmail');
var moment = require('moment');
var util = require('util');
var _ = require('validator');
var hostname = require('../../config/domain').api.domain;

var createMail = function (mail) {
    var message;
    switch (!!mail.text << !!mail.html) {
        case 2:
            message = new BuildMail('multipart/mixed');
            message.createChild('text/plain').setContent(mail.text);
            message.createChild('text/html').setContent(mail.html);
            break;
        case 1:
            message = new BuildMail('text/plain').setContent(mail.text);
            break;
        case 0:
            message = new BuildMail('text/html').setContent(mail.html);
            break;
    }
    message.setHeader({
        from: mail.from,
        to: mail.to,
        subject: mail.subject
    });
    mail['reply-to'] && message.setHeader('reply-to', mail['reply-to']);
    return message;
};

var buildMailBuf = function (req, mail, message, cb) {
    mail.headers = mail.headers || {};
    if (!mail.headers.received) {
        message.addHeader('received',
            util.format('from [%s] (%s [%s])\r\nby %s via HTTP for <%s>;\r\n%s',
                req.ip, 'unknown', req.ip, hostname,
                message.getAddresses()['to'][0]['address'],
                moment().format("ddd, D MMM YYYY HH:mm:ss ZZ")
            ));
    }
    for (var key in mail.headers) {
        if (mail.headers.hasOwnProperty(key)) {
            var value = mail.headers[key];
            if (key.toLowerCase() !== 'received') {
                key = 'x-' + key;
            }
            message.addHeader(key, value);
        }
    }
    message.build(cb);
};

router.all('/message', function (req, res) {
    var mail = req.$params;
    if (!(mail.from && mail.to && ( mail.html || mail.text) )) {
        return res.jsonp({code: 406, error: 'Request content not acceptable.'});
    }
    var message = createMail(mail);
    var from = message.getAddresses()['from'][0]['address'];
    var to = message.getAddresses()['to'][0]['address'];
    if (!_.isEmail(from) || !_.isEmail(to)) {
        return res.jsonp({code: 406, error: 'Request content not acceptable.'});
    }
    buildMailBuf(req, mail, message, function (err, mailBuf) {
        if (err !== null) {
            return res.jsonp({code: 500, error: 'Unknown error.'});
        }
        Mail.create({
            from: from,
            to: to,
            subject: mail.subject,
            web_hook: req.sender.web_hook,
            immediate: mail.immediate === undefined && req.sender.immediate || mail.immediate
        }).complete(function (err, mailRecord) {
            if (err) {
                return res.jsonp({code: 500, error: 'Save mail record with error.'});
            }
            req.sender.addMail(mailRecord).complete(function (err) {
                if (err) {
                    return res.jsonp({code: 500, error: 'Unknown error.'});
                }
                newCommand(req.sender.id, {command: 'sendMail', mailId: mailRecord.id}, function (err) {
                    if (err) {
                        return res.jsonp({code: 500, error: 'Unknown error.'});
                    }
                    Mail.write(mailRecord.id, mailBuf, function () {
                        res.jsonp({
                            code: 200,
                            msg: 'in queue',
                            id: mailRecord.id
                        })
                    });
                });
            });
        });
    });
});
