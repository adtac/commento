(function (global, document) {
  "use strict";

  global.resetPassword = function() {
    var allOk = global.unfilledMark(["#password", "#password2"], function(el) {
      el.css("border-bottom", "1px solid red");
    });

    if (!allOk) {
      global.textSet("#err", "Please make sure all fields are filled.");
      return;
    }

    if ($("#password").val() !== $("#password2").val()) {
      global.textSet("#err", "The two passwords do not match.");
      return;
    }

    var json = {
      "resetHex": global.paramGet("hex"),
      "password": $("#password").val(),
    };

    global.buttonDisable("#reset-button");
    global.post(global.origin + "/api/owner/reset-password", json, function(resp) {
      global.buttonEnable("#reset-button");

      global.textSet("#err", "");
      if (!resp.success) {
        global.textSet("#err", resp.message);
        return
      }

      document.location = global.origin + "/login?changed=true";
    });
  }

} (window.commento, document));
