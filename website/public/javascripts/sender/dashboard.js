(function ($) {

    $.get(window.location.pathname + '/dns-records', function (data) {
        $('#txt-currentSPF').val(data['spf']);
        $('#txt-currentDKIM').val(data['dkim']);
    }, 'json');

    $('#btn-checkDNS').click(function () {
        var $btn = $(this), $form = $('#form-checkDNS');
        $form.addClass('loading');
        $btn.addClass('disabled');
        $.ajax(window.location.pathname + '/check-dns', {
            type: 'POST',
            dataType: 'json'
        }).always(function () {
            $form.removeClass('loading');
            $btn.removeClass('disabled');
        }).done(function (data) {
            $form.find('.negative').addClass('negative');
            if (data['spf'] && data['dkim']) {
                window.location.reload();
                return;
            }
            $('#txt-currentSPF').val(data.records['spf']);
            $('#txt-currentDKIM').val(data.records['dkim']);
            if (!data['spf']) {
                $form.find('.spf').addClass('negative');
            }
            if (!data['dkim']) {
                $form.find('.dkim').addClass('negative');
            }
            swal("Oops...", 'Some of DNS records have not propagated', "error");
        }, 'json');
    });
})(jQuery);