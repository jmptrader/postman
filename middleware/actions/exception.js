var noExistKeywords1 = 'address|user|recipient|account|mailbox|alias'.split('|');
var noExistKeywords2 = 'unknown|exist|found|unavailable|disabled|no such|have'.split('|');
var errKeywords = 'blocked'.split('|');


var defaultTreatment = function (log) {
    var needResend = false;
    var treatment = 'resendLater';
    errKeywords.forEach(function (kw) {
        if (log.indexOf(kw) !== -1) {
            needResend = true;
        }
    });
    if (needResend) return 'resendLater';
    var accountKeywordExist = false;
    noExistKeywords1.forEach(function (kw) {
        if (log.indexOf(kw) !== -1) {
            accountKeywordExist = true;
        }
    });
    if (!accountKeywordExist) return 'resendLater';
    noExistKeywords2.forEach(function (kw) {
        if (log.indexOf(kw) !== -1) {
            treatment = 'ignore';
        }
    });
    return treatment;
};

Action.register('exception', function (args, id) {
    var log = args.log;
    var c = this;
    if (!c.auth) return;
    Exception.findOrCreate({where: {content: log}, defaults: {
        treatment: defaultTreatment(log)
    }}).complete(function (err, e) {
        if (err) return;
        return c.response(id, e[0].treatment)
    });
});
