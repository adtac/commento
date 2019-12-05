(function (global, document) {
  "use strict";

  (document);

  // Update the email records.
  global.update = function(event) {
    event.preventDefault();

    $(".err").text("");
    $(".msg").text("");

    var allOk = global.unfilledMark(["#name", "#email"], function(el) {
      el.css("border-bottom", "1px solid red");
    });

    if (!allOk) {
      global.textSet("#err", "Please make sure all fields are filled");
      return;
    }

    var json = {
      "commenterToken": global.paramGet("commenterToken"),
      "name": $("#name").val(),
      "email": $("#email").val(),
      "link": $("#link").val(),
      "photo": $("#photo").val(),
    };

    global.buttonDisable("#save-button");
    global.post(global.origin + "/api/commenter/update", json, function(resp) {
      global.buttonEnable("#save-button");
      if (!resp.success) {
        $(".err").text(resp.message);
        return;
      }

      $(".msg").text("Successfully updated!");
    });
  }

  global.profilePrefill = function() {
    $(".err").text("");
    $(".msg").text("");
    var json = {
      "commenterToken": global.paramGet("commenterToken"),
    };

    global.post(global.origin + "/api/commenter/self", json, function(resp) {
      $("#loading").hide();
      $("#form").show();
      if (!resp.success) {
        $(".err").text(resp.message);
        return;
      }

      $("#name").val(resp.commenter.name);
      $("#email").val(resp.commenter.email);
      $("#unsubscribe").attr("href", global.origin + "/unsubscribe?unsubscribeSecretHex=" + resp.email.unsubscribeSecretHex);

      if (resp.commenter.provider === "commento") {
        $("#link-row").attr("style", "")
        if (resp.commenter.link !== "undefined") {
          $("#link").val(resp.commenter.link);
        }

        $("#photo-row").attr("style", "")
        $("#photo-subtitle").attr("style", "")
        if (resp.commenter.photo !== "undefined") {
          $("#photo").val(resp.commenter.photo);
        }
      }
    });
  };

} (window.commento, document));
