router.all('/logs', function (req, res) {
    var limit = req.param('limit') || req.body.limit || req.$params.limit || 0;
    var offset = req.param('offset') || req.body.offset || req.$params.offset || 0;
    model.query('SELECT COUNT(id) FROM jianxin.mails WHERE 1=1;').complete(function (err, result) {
        var length = err ? 0 : result[0]['COUNT(id)'];
        if (!length) {
            return res.jsonp({code: 200, total: length, logs: []});
        }
        model.query('SELECT mails.from, mails.to, mails.subject, mails.id, mails.created_at, logs.log, logs.status, logs.created_at AS log_created_at\
        FROM jianxin.mails\
        LEFT JOIN jianxin.logs\
        ON logs.mail_id= mails.id\
        WHERE mails.id in (\
            SELECT * FROM (\
                SELECT mails.id FROM jianxin.mails WHERE mails.sender_id = :senderId ORDER BY mails.id DESC LIMIT :offset, :limit\
            ) AS t\
        )\
        ORDER BY mails.id DESC', null,
            { raw: true }, { offset: offset, limit: limit, senderId: req.sender.id }).complete(function (_, logs) {
                return res.jsonp({code: 200, total: length, logs: logs});
            });
    });
});