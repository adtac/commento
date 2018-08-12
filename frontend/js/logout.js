(function (global, document) {

  global.logout = function() {
    global.cookieDelete("commentoOwnerToken");
    document.location = global.commentoOrigin + "/login";
  }

} (window, document));
