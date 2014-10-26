(function ($) {
    var $LogList = $('#mailLogsList');
    _.templateSettings = {
        interpolate: /\{\{(.+?)\}\}/g,
        evaluate: /{%([\s\S]+?)%}/g
    };
    var template = _.template($("#mailLogItem").html());

    window.getAvatar = function (email) {
        return 'http://www.gravatar.com/avatar/' + md5(email.toLowerCase())
    };
    var itemsPerpage = 15;
    var updateLogs = function () {
        var pages = parseInt(window.location.hash.substr(1), 10) || 1;
        var params = {
            limit: itemsPerpage,
            offset: (pages - 1) * itemsPerpage,
            expire: Math.round(new Date().getTime() / 1000) + 10 * 60
        };
        $.getJSON(window.api_url + '?callback=?', {
            params: JSON.stringify(params),
            secret: md5(JSON.stringify(params).toLowerCase() + window.api_key)
        }, function (res) {
            if (res.code !== 200) {
                return swal("Oops...", "Network error, please try later!", "error");
            }
            if (res.logs.length === 0) {
                return $LogList.html('<p>No record found.</p>')
            }
            var result = {current_page: pages, total: Math.ceil(res.total / itemsPerpage), max: res.logs[0].id, min: res.logs[res.logs.length - 1].id, logs: {}};
            _.each(res.logs, function (log) {
                if (result.logs.hasOwnProperty(log.id)) {
                    result.logs[log.id].logs.push({
                        log: log.log,
                        status: log.status,
                        created_at: log.log_created_at
                    });
                    if (log.status === 'delivered') {
                        result.logs[log.id].status = 'delivered';
                    }
                    return;
                }
                result.logs[log.id] = {
                    from: log.from,
                    to: log.to,
                    subject: log.subject,
                    created_at: log.created_at,
                    status: log.status,
                    logs: !log.status ? [] : [
                        {
                            log: log.log,
                            status: log.status,
                            created_at: log.log_created_at
                        }
                    ]
                }
            });
            $LogList.html(template(result));
            $LogList.find('.item').popup({
                position: 'bottom center'
            });
        })
    };
    window.onhashchange = updateLogs;
    updateLogs();
})(jQuery);