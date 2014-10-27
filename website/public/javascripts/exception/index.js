(function ($) {
    var current_type = 'unknown';
    window.location.search.substr(1).split('?').forEach(function (v) {
        var pair = v.split('=');
        if (pair[0] === 'type') {
            current_type = pair[1];
        }
    });
    $('.type-' + current_type).addClass('active');

    $('form').submit(function (env) {
        env.preventDefault();
        var $form = $(this);
        $.ajax($form[0].action, {
            type: 'POST',
            data: $form.serialize(),
            dataType: 'json'
        }).done(function (data) {
            if (data.code !== 200) {
                swal("Oops...", data.errors, "error");
                return;
            }
            swal("Success!", "treatment is set success.", "success");
            setTimeout(function () {
                window.location.reload();
            }, 1000);
        }).error(function () {
            swal("Oops...", "Network error, please try later!", "error");
        });
    });
})(jQuery);