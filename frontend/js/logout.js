(function (global, document) {

  global.logout = function() {
    global.cookieSet("commentoOwnerToken", "");
    document.location = global.commentoOrigin + "/login";
  }

} (window, document));
