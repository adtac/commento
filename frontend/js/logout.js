(function (global, document) {

  global.logout = function() {
    global.cookieSet("session", "");
    document.location = "/login";
  }

} (window, document));
