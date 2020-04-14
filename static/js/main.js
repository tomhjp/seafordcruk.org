$(function()
{
    var start = Date.now()

    // For some reason, redirects don't always seem to clear the form...
    $('.text-input').each(function() {this.value=""});

    var contactName = $('#contact-name');
    var contactEmail = $('#contact-email');
    var contactMessage = $('#contact-message');
    var submitButton = $('#contact-form-submit');
    submitButton.prop('disabled', true);

    $('#contact-form').submit(function() {
        $('#time-input').val(Date.now() - start);
    })

    function lengthNotExcessive(val) {
        return val.length < 10000;
    }
    function contactEmailValid() {
        var email = $('#contact-email').val();
        var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
        return lengthNotExcessive(email) && re.test(String(email).toLowerCase());
    }
    function contactNameValid() {
        return lengthNotExcessive(contactName.val()) && contactName.val() != '';
    }
    function contactMessageValid() {
        return lengthNotExcessive(contactMessage.val());
    }
    function updateSubmitButton() {
        if (contactNameValid() && contactEmailValid() && contactMessageValid()) {
            submitButton.prop('disabled', false);
        }
        else {
            submitButton.prop('disabled', true);
        }
    }

    var contactNameCheck = $('#contact-name-check');
    var contactNameTimes = $('#contact-name-times');
    var contactEmailCheck = $('#contact-email-check');
    var contactEmailTimes = $('#contact-email-times');
    function display(item) {
        item.css("display", "inherit");
    }
    function hide(item) {
        item.css("display", "none");
    }
    contactName.keyup(function() {
        if (contactNameValid()) {
            display(contactNameCheck);
            hide(contactNameTimes);
        }
        else {
            hide(contactNameCheck);
            display(contactNameTimes);
        }

        updateSubmitButton();
    });

    contactEmail.keyup(function() {
        if (contactEmailValid()) {
            display(contactEmailCheck);
            hide(contactEmailTimes);
        }
        else {
            hide(contactEmailCheck);
            display(contactEmailTimes);
        }

        updateSubmitButton();
    });

    contactMessage.keyup(function() {
        updateSubmitButton();
    });
});