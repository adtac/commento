(function (global, document) {
  "use strict"

  // Signs up the user and redirects to either the login page or the email
  // confirmation, depending on whether or not SMTP is configured in the
  // backend.
  global.signup = function(event) {
    event.preventDefault();

    if ($("#password").val() !== $("#password2").val()) {
      global.textSet("#err", "The two passwords don't match");
      return;
    }

    var allOk = global.unfilledMark(["#email", "#name", "#password", "#password2"], function(el) {
      el.css("border-bottom", "1px solid red");
    });

    if (!allOk) {
      global.textSet("#err", "Please make sure all fields are filled");
      return;
    }

    var json = {
      "email": $("#email").val(),
      "name": $("#name").val(),
      "password": $("#password").val(),
    };

    global.buttonDisable("#signup-button");
    global.post(global.origin + "/api/owner/new", json, function(resp) {
      global.buttonEnable("#signup-button")

      if (!resp.success) {
        global.textSet("#err", resp.message);
        return;
      }

      if (resp.confirmEmail) {
        document.location = global.origin + "/confirm-email";
      } else {
        document.location = global.origin + "/login?signedUp=true";
      }
    });
  };

} (window.commento, document));
