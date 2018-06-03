(function (global, document) {

  // Get self details.
  global.selfGet = function(callback) {
    var json = {
      "session": global.cookieGet("session"),
    };

    global.post(global.origin + "/api/owner/self", json, function(resp) {
      if (!resp.success || !resp.loggedIn) {
        document.location = "/login";
        return;
      }

      global.owner = resp.owner;
      callback();
    });
  };

}(window, document));
