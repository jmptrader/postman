const MAIL_SYNC_KEY = 'jianxin:mailSync';

global.MailSync = {
    _syncMap: {},
    listen: function (id, next) {
        this._syncMap[id] = next;
    },
    emit: function (id, result) {
        redisClient.LPUSH(MAIL_SYNC_KEY, [id, JSON.stringify(result)].join(':'));
    },
    loop: function () {
        var loop = this.loop;
        var self = this;
        redisBlockClient.BRPOP(MAIL_SYNC_KEY, 0, function (err, resultArr) {
            var mail_id = resultArr[1].split(':')[0];
            var result = JSON.parse(resultArr[1].substr(mail_id.length + 1));
            if (self._syncMap.hasOwnProperty(mail_id)) {
                self._syncMap[mail_id](result);
            }
            delete(self._syncMap[mail_id]);
            loop.call(self);
        });
    }
};
