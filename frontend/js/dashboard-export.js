(function (global, document) {
  "use strict";

  (document);

  global.domainExportBegin = function() {
    var data = global.dashboard.$data;

    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
      "domain": data.domains[data.cd].domain,
    }

    global.buttonDisable("#domain-export-button");
    global.post(global.origin + "/api/domain/export/begin", json, function(resp) {
      global.buttonEnable("#domain-export-button");
      if (!resp.success) {
        global.globalErrorShow(resp.message);
        return;
      }

      global.globalOKShow("Data export operation has been successfully queued. You will receive an email.");
    });
  };

} (window.commento, document));
