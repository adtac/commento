(function (global, document) {
  "use strict";

  (document);

  // Opens the import window.
  global.importOpen = function() {
    $(".view").hide();
    $("#import-view").show();
  }

  global.importDisqus = function() {
    var url = $("#disqus-url").val();
    var data = global.dashboard.$data;

    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
      "domain": data.domains[data.cd].domain,
      "url": url,
    }

    global.buttonDisable("#disqus-import-button");
    global.post(global.origin + "/api/domain/import/disqus", json, function(resp) {
      global.buttonEnable("#disqus-import-button");

      if (!resp.success) {
        global.globalErrorShow(resp.message);
        return;
      }

      $("#disqus-import-button").hide();

      global.globalOKShow("Imported " + resp.numImported + " comments!");
    });
  }

  global.importCommento = function() {
    var url = $("#commento-url").val();
    var data = global.dashboard.$data;

    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
      "domain": data.domains[data.cd].domain,
      "url": url,
    }

    global.buttonDisable("#commento-import-button");
    global.post(global.origin + "/api/domain/import/commento", json, function(resp) {
      global.buttonEnable("#commento-import-button");

      if (!resp.success) {
        global.globalErrorShow(resp.message);
        return;
      }

      $("#commento-import-button").hide();

      global.globalOKShow("Imported " + resp.numImported + " comments!");
    });
  }

} (window.commento, document));
