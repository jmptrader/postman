var keepIf = function (objs, attrs) {
    var newObjs = [];
    objs.forEach(function (obj) {
        var newObj = {};
        attrs.forEach(function (attr) {
            newObj[attr] = obj[attr];
        });
        newObjs.push(newObj);
    });
    return newObjs;
};

router.all('/logs', function (req, res) {
    var limit = req.param('limit') || req.body.limit || req.$params.limit;
    var offset = req.param('offset') || req.body.offset || req.$params.offset;
    var result = [];
    req.sender.getMails({
        order: [
            ["created_at", "DESC"]
        ],
        limit: limit,
        offset: offset
    }).success(function (mails) {
        var len = mails.length;
        mails.forEach(function (mail) {
            mail.getLogs({order: [
                ["created_at", "DESC"]
            ]}).complete(function (_, logs) {
                result.push({
                    logs: keepIf(logs, ['created_at', 'status', 'log']),
                    created_at: mail.created_at,
                    from: mail.from,
                    to: mail.to,
                    subject: mail.subject
                });
                len -= 1;
                if (len <= 0) {
                    res.jsonp(result.sort(function (a, b) {
                        return (new Date(a)).getTime() - (new Date(b)).getTime();
                    }));
                }
            })
        });
    });
});