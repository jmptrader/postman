var crypto = require('crypto');

Action.register('auth', function (args) {
    var md5 = crypto.createHash('md5');
    var c = this;
    md5.update(this.sender.secret + this.auth_key);
    delete this['auth_key'];
    if (args.result !== md5.digest('hex')) {
        this.command('exit', {
            "msg": 'secret not match'
        });
        console.log('sender auth fail');
        return
    }
    c.sender.status = 'online';
    c.sender.save(['status']).success(function () {
        console.log('sender: ' + c.sender.ip + ' auth success.');
        this.auth = true;
    });
});