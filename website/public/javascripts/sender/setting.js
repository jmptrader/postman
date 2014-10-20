(function ($) {
    var $RoleBody = $('#tbody-frequenciesList');
    var updateRoleList = function () {
        $.get('/sender/' + $RoleBody.attr('data-senderId') + '/frequencies.html', function (html) {
            $RoleBody.html(html);
        });
    };
    updateRoleList();

    $('.updateForm').submit(function (e) {
        e.preventDefault();
        var $form = $(this);
        $form.addClass('loading');
        $.ajax(this.action, {
            type: 'POST',
            data: $form.serialize(),
            dataType: 'json'
        }).always(function () {
            $form.removeClass('loading');
        }).done(function (data) {
            if (data.code !== 200) {
                $form[0].reset();
                swal("Oops...", data.errors, "error");
                return;
            }
            swal("Success!", "Basic setting save success.", "success");
        }).error(function () {
            swal("Oops...", "Network error, please try later!", "error");
        });
    });

    $('#form-createRole').form({
        domain: {
            identifier: 'frequency[domain]',
            rules: [
                {
                    type: 'empty',
                    prompt: 'Please enter a domain'
                },
                {
                    type: 'regExp[^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+.?$]',
                    prompt: 'Please enter a valid domain'
                }
            ]
        }, deliver_frequency: {
            identifier: 'frequency[deliver_frequency]',
            rules: [
                {
                    type: 'empty',
                    prompt: 'Please enter a number'
                },
                {
                    type: 'regExp[^([0-9]|[1-5][0-9]|60)$]',
                    prompt: 'frequency should be a number not more than 60.'
                }
            ]
        }
    }, {
        inline: true
    }).submit(function (env) {
        env.preventDefault();
        $('#btn-submitNewRole').trigger('click');
    });

    var createSenderModal = $('.create-role.modal')
        .modal('setting', 'transition', 'pulse')
        .modal('setting', 'closable', false)
        .modal('setting', 'onApprove', function () {
            return submitForm();
        });


    var submitForm = function () {
        var $form = $('#form-createRole');
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
            swal("Success!", "New role create success.", "success");
            updateRoleList();
        }).error(function () {
            swal("Oops...", "Network error, please try later!", "error");
        });
    };

    $('#btn-createNewRole').on('click', function () {
        createSenderModal.modal('show');
    });

    $RoleBody.on('click', '.remove', function (e) {
        e.preventDefault();
        var $this = $(this);
        swal({
            title: "Are you sure?",
            text: "Role will be disabled!",
            type: "warning",
            showCancelButton: true,
            confirmButtonColor: "#DD6B55",
            confirmButtonText: "Yes, delete it!",
            closeOnConfirm: false
        }, function () {
            $.post($this.attr('href'), function (result) {
                if (result.code === 200) {
                    swal("Success!", "Role destroy success.", "success");
                    return  updateRoleList();
                }
                swal("Oops...", result.error, "error");
            }, 'json');
        });
    });

})(jQuery);