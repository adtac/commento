(function (global, document) {

  // Get self details.
  global.selfGet = function(callback) {
    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
    };

    global.post(global.commentoOrigin + "/api/owner/self", json, function(resp) {
      if (!resp.success || !resp.loggedIn) {
        document.location = "/login";
        return;
      }

      global.owner = resp.owner;
      callback();
    });
  };

}(window, document));
