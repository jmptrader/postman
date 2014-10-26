var path = require('path');
var moment = require('moment');
var fs = require('fs');

const MAIL_FILE_EXT = '.ml';

global.Mail = model.define('Mail', {
    from: Model.STRING,
    to: Model.STRING,
    subject: Model.STRING,
    block: Model.BOOLEAN,
    immediate: Model.BOOLEAN
}, {
    tableName: 'mails',
    updatedAt: 'updated_at',
    createdAt: 'created_at',
    instanceMethods: {
        write: function (content, cb) {
            var self = this;
            var filePath = path.join(archiveDir, moment(new Date(self.created_at)).format("YYYYMMDD"));
            ensureExists(filePath, 0744, function (err) {
                if (err !== null) throw err;
                fs.writeFile(path.join(filePath, self.id + MAIL_FILE_EXT), content, function (err) {
                    if (err) throw err;
                    cb && cb();
                });
            });
        }, read: function (cb) {
            var self = this;
            var filePath = path.join(archiveDir, moment(new Date(self.created_at)).format("YYYYMMDD"), self.id + MAIL_FILE_EXT);
            fs.readFile(filePath, function (err, data) {
                if (err) return cb(null);
                cb(data.toString());
            });
        }
    }
});

global.Log = model.define('Log', {
    status: Model.STRING,
    log: Model.STRING
}, {
    tableName: 'logs',
    updatedAt: 'updated_at',
    createdAt: 'created_at'
});

var archiveDir = process.env['POSTMAN_CONFIG_DIR'] || '../mail_archive';
var ensureExists = function (path, mask, cb) {
    if (typeof mask == 'function') {
        cb = mask;
        mask = 0777;
    }
    fs.mkdir(path, mask, function (err) {
        if (err) {
            if (err.code == 'EEXIST') cb(null);
            else cb(err);
        } else cb(null);
    });
};

Mail.hasMany(Log, {as: 'Logs', foreignKey: 'mail_id'});