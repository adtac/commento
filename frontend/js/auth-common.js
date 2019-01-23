(function (global, document) {
  "use strict";

  // Redirect the user to the dashboard if there's a cookie. If the cookie is
  // invalid, they would be redirected back to the login page *after* the
  // cookie is deleted.
  global.loggedInRedirect = function() {
    if (global.cookieGet("commentoOwnerToken") !== undefined) {
      document.location = global.origin + "/dashboard";
    }
  }


  // Prefills the email field from the URL parameter.
  global.prefillEmail = function() {
    if (global.paramGet("email") !== undefined) {
      $("#email").val(global.paramGet("email"));
      $("#password").click();
    }
  };

} (window.commento, document));
