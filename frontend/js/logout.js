(function (global, document) {
  "use strict";

  global.logout = function() {
    global.cookieDelete("commentoOwnerToken");
    document.location = global.origin + "/login";
  }

} (window.commento, document));
