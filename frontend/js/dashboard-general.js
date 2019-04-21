(function (global, document) {
  "use strict";

  (document);

  // Opens the general settings window.
  global.generalOpen = function() {
    $(".view").hide();
    $("#general-view").show();
  };

  global.generalSaveHandler = function() {
    var data = global.dashboard.$data;

    global.buttonDisable("#save-general-button");
    global.domainUpdate(data.domains[data.cd], function() {
      global.globalOKShow("Settings saved!");
      global.buttonEnable("#save-general-button");
    });
  };

  global.ssoProviderChangeHandler = function() {
    var data = global.dashboard.$data;

    if (data.domains[data.cd].ssoSecret === "") {
      var json = {
        "ownerToken": global.cookieGet("commentoOwnerToken"),
        "domain": data.domains[data.cd].domain,
      };

      global.post(global.origin + "/api/domain/sso/new", json, function(resp) {
        if (!resp.success) {
          global.globalErrorShow(resp.message);
          return;
        }

        data.domains[data.cd].ssoSecret = resp.ssoSecret;
        $("#sso-secret").val(data.domains[data.cd].ssoSecret);
      });
    } else {
      $("#sso-secret").val(data.domains[data.cd].ssoSecret);
    }
  };

} (window.commento, document));
