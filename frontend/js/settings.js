(function (global, document) {
  "use strict";

  (document);

  global.vueConstruct = function(callback) {
    var reactiveData = {
      hasSource: global.owner.hasSource,
      lastFour: global.owner.lastFour,
    };

    global.settings = new Vue({
      el: "#settings",
      data: reactiveData,
    });

    if (callback !== undefined) {
      callback();
    }
  };

  global.settingShow = function(setting) {
    $(".pane-setting").removeClass("selected");
    $(".view").hide();
    $("#" + setting).addClass("selected");
    $("#" + setting + "-view").show();
  };

  global.deleteOwnerHandler = function() {
    if (!confirm("Are you absolutely sure you want to delete your account?")) {
      return;
    }

    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
    }

    $("#delete-owner-button").prop("disabled", true);
    $("#delete-owner-button").text("Deleting...");
    global.post(global.origin + "/api/owner/delete", json, function(resp) {
      if (!resp.success) {
        $("#delete-owner-button").prop("disabled", false);
        $("#delete-owner-button").text("Delete Account");
        global.globalErrorShow(resp.message);
        $("#error-message").text(resp.message);
        return;
      }

      global.cookieDelete("commentoOwnerToken");
      document.location = global.origin + "/login?deleted=true";
    });
  };

} (window.commento, document));
