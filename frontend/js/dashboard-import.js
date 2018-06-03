(function (global, document) {

  // Opens the import window.
  global.importOpen = function() {
    $(".view").hide();
    $("#import-view").show();
  }


  global.importDisqus = function() {
    var url = $("#disqus-url").val();
    var data = global.dashboard.$data;

    var json = {
      session: global.cookieGet("session"),
      domain: data.domains[data.cd].domain,
      url: url,
    }

    global.buttonDisable("#disqus-import-button");
    global.post(global.origin + "/api/import/disqus", json, function(resp) {
      global.buttonEnable("#disqus-import-button");

      if (!resp.success) {
        global.globalErrorShow(resp.message);
        return;
      }

      globalOKShow("Imported " + resp.numImported + " comments!");
    });
  }

} (window, document));
