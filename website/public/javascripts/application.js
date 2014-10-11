!function ($) {
    $('.pp').popup();

    $('.message .close').on('click', function () {
        $(this).closest('.message').fadeOut();
    });

    // auto choose current nav by pathname
    var doc = document;
    var path = window.location.pathname,
        nav = path.substr(1).split('/')[0],
        $nav = doc.getElementById('js:nav-' + nav);
    if (!$nav) {
        $nav = doc.getElementById('js:nav-index');
    }
    $($nav).addClass('active');

    $.fn.form.settings.rules.regExp = function (value, re) {
        var reg = new RegExp(re);
        return reg.test(value);
    };
}(jQuery);