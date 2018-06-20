(function (global, document) {

  global.logout = function() {
    global.cookieSet("commentoOwnerToken", "");
    document.location = "/login";
  }

} (window, document));
