(function (global, document) {

  global.logout = function() {
    global.cookieSet("ownerToken", "");
    document.location = "/login";
  }

} (window, document));
