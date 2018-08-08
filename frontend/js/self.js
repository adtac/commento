(function (global, document) {

  // Get self details.
  global.selfGet = function(callback) {
    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
    };

    if (json.ownerToken === undefined) {
      document.location = "/login";
      return;
    }

    global.post(global.commentoOrigin + "/api/owner/self", json, function(resp) {
      if (!resp.success || !resp.loggedIn) {
        global.cookieDelete("commentoOwnerToken");
        document.location = "/login";
        return;
      }

      global.owner = resp.owner;
      callback();
    });
  };

}(window, document));
