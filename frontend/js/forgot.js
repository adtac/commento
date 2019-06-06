(function (global, document) {
  "use strict";

  (document);

  // Talks to the API and sends an reset email.
  global.sendResetHex = function(event) {
    event.preventDefault();

    var allOk = global.unfilledMark(["#email"], function(el) {
      el.css("border-bottom", "1px solid red");
    });

    if (!allOk) {
      global.textSet("#err", "Please make sure all fields are filled.");
      return;
    }

    var entity = "owner";
    if (global.paramGet("commenter") === "true") {
      entity = "commenter";
    }

    var json = {
      "email": $("#email").val(),
      "entity": entity,
    };

    global.buttonDisable("#reset-button");
    global.post(global.origin + "/api/forgot", json, function(resp) {
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
