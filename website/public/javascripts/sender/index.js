(function ($) {
    // add form validation
    var $form = $('#form-createSender');
    var $senderManager = $('#sender-manager');

    var submitForm = function () {
        $form.form('validate form');
        if ($form.find('.error').length > 0)return false;
        $.ajax($form[0].action, {
            type: 'POST',
            data: $form.serialize(),
            dataType: 'json'
        }).always(function () {
            $form[0].reset();
        }).done(function (data) {
            if (data.code !== 200) {
                swal("Oops...", data.errors, "error");
                return;
            }
            swal("Success!", "New sender create success.", "success");
            setTimeout(function () {
                window.location.href = '/sender/' + data['sender_id'];
            }, 1000);
        }).error(function () {
            swal("Oops...", "Network error, please try later!", "error");
        });
    };

    $form.form({
        ip: {
            identifier: 'sender[ip]',
            rules: [
                {
                    type: 'empty',
                    prompt: 'Please enter a ip'
                },
                {
                    type: 'regExp[^((1?\\d?\\d|(2([0-4]\\d|5[0-5])))\\.){3}(1?\\d?\\d|(2([0-4]\\d|5[0-5])))$]',
                    prompt: 'Please enter a valid ip record'
                }
            ]
        },
        domain: {
            identifier: 'sender[domain]',
            rules: [
                {
                    type: 'empty',
                    prompt: 'Please enter a domain'
                },
                {
                    type: 'regExp[^([0-9a-z_-]+\\.)*?[a-z0-9-]+\\.[a-z]{2,6}(\\.[a-z]{2})?$]',
                    prompt: 'Please enter a valid domain record'
                }
            ]
        }
    }, {
        inline: true
    }).submit(function (env) {
        env.preventDefault();
        $('#btn-submitSender').trigger('click');
    });


    var createSenderModal = $('.create-sender.modal')
        .modal('setting', 'transition', 'pulse')
        .modal('setting', 'closable', false)
        .modal('setting', 'onApprove', function () {
            return submitForm();
        });

    $senderManager.on('click', '#btn-createSender', function () {
        createSenderModal.modal('show');
    });

    window.onhashchange = function () {
        var currentStatus = window.location.hash.substr(1);
        if (['online', 'offline'].indexOf(currentStatus) === -1) {
            currentStatus = 'all';
        }
        var $statusNav = $('#nav-currentStatus');
        $statusNav.find('.active').removeClass('active');
        $statusNav.find('.status-' + currentStatus).addClass('active');
        $.get('/sender/list.html?status=' + currentStatus, function (html) {
            $('#holder-senderList').html(html);
        });
    };

    window.onhashchange();

    $senderManager.on('click', '.remove', function (e) {
        e.preventDefault();
        var $this = $(this);
        swal({
            title: "Are you sure?",
            text: "Sender will be disabled!",
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "Yes, delete it!",
            closeOnConfirm: false
        }, function () {
            $.post($this.attr('href'), function (result) {
                if (result.code === 200) {
                    swal("Success!", "Sender destroy success.", "success");
                    return  window.onhashchange();
                }
                swal("Oops...", result.error, "error");
            }, 'json');
        });
    });

})(jQuery);