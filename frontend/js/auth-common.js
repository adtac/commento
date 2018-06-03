(function (global, document) {

  // Prefills the email field from the URL parameter.
  global.prefillEmail = function() {
    if (paramGet("email") != undefined) {
      $("#email").val(paramGet("email"));
      $("#password").click();
    }
  };

} (window, document));
