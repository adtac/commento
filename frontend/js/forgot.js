(function (global, document) {
  "use strict";

  (document);

  // Talks to the API and sends an reset email.
  global.sendResetHex = function() {
    var allOk = global.unfilledMark(["#email"], function(el) {
      el.css("border-bottom", "1px solid red");
    });

    if (!allOk) {
      global.textSet("#err", "Please make sure all fields are filled.");
      return;
    }

    var json = {
      "email": $("#email").val(),
    };

    global.buttonDisable("#reset-button");
    global.post(global.origin + "/api/owner/send-reset-hex", json, function(resp) {
      global.buttonEnable("#reset-button");

      global.textSet("#err", "");
      if (!resp.success) {
        global.textSet("#err", resp.message);
        return
      }

      $("#msg").html("If that email is a registered account, you will receive an email with instructions on how to reset your password.");
      $("#reset-button").hide();
    });
  }

} (window.commento, document));
