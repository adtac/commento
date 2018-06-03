(function (global, document) {

  // Shows messages produced from email confirmation attempts.
  function displayConfirmedEmail() {
    var confirmed = global.paramGet("confirmed");

    if (confirmed == "true") {
      $("#msg").html("Successfully confirmed! Login to continue.")
    }
    else if (confirmed == "false") {
      $("#err").html("That link has expired.")
    }
  }


  // Shows messages produced from password reset attempts.
  function displayChangedPassword() {
    var changed = paramGet("changed");

    if (changed == "true") {
      $("#msg").html("Password changed successfully! Login to continue.")
    }
  }

  // Shows messages produced from completed signups.
  function displaySignedUp() {
    var signedUp = paramGet("signedUp");

    if (signedUp == "true") {
      $("#msg").html("Registration successful! Login to continue.")
    }
  }


  // Shows email confirmation and password reset messages, if any.
  global.displayMessages = function() {
    displayConfirmedEmail();
    displayChangedPassword();
    displaySignedUp();
  };


  // Logs the user in and redirects to the dashboard.
  global.login = function() {
    var all_ok = global.unfilledMark(["#email", "#password"], function(el) {
      el.css("border-bottom", "1px solid red");
    });

    if (!all_ok) {
      global.textSet("#err", "Please make sure all fields are filled");
      return;
    }

    var json = {
      "email": $("#email").val(),
      "password": $("#password").val(),
    };

    global.buttonDisable("#login-button");
    global.post(global.origin + "/api/owner/login", json, function(resp) {
      global.buttonEnable("#login-button");

      if (!resp.success) {
        global.textSet("#err", resp.message);
        return;
      }

      global.cookieSet("session", resp.session);
      document.location = "/dashboard";
    });
  };

} (window, document));
