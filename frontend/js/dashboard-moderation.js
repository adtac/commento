(function (global, document) {
  "use strict";

  (document);

  // Opens the moderatiosn settings window.
  global.moderationOpen = function() {
    $(".view").hide();
    $("#moderation-view").show();
  };

  
  // Adds a moderator.
  global.moderatorNewHandler = function() {
    var data = global.dashboard.$data;
    var email = $("#new-mod").val();
    
    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
      "domain": data.domains[data.cd].domain,
      "email": email,
    }

    var idx = -1;
    for (var i = 0; i < data.domains[data.cd].moderators.length; i++) {
      if (data.domains[data.cd].moderators[i].email === email) {
        idx = i;
        break;
      }
    }

    if (idx === -1) {
      data.domains[data.cd].moderators.push({"email": email, "timeAgo": "just now"});
      global.buttonDisable("#new-mod-button");
      global.post(global.origin + "/api/domain/moderator/new", json, function(resp) {
        global.buttonEnable("#new-mod-button");

        if (!resp.success) {
          global.globalErrorShow(resp.message);
          return
        }

        global.globalOKShow("Added a new moderator!");
        $("#new-mod").val("");
        $("#new-mod").focus();
      });
    } else {
      global.globalErrorShow("Already a moderator.");
    }
  }


  // Deletes a moderator.
  global.moderatorDeleteHandler = function(email) {
    var data = global.dashboard.$data;
    
    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
      "domain": data.domains[data.cd].domain,
      "email": email,
    }

    var idx = -1;
    for (var i = 0; i < data.domains[data.cd].moderators.length; i++) {
      if (data.domains[data.cd].moderators[i].email === email) {
        idx = i;
        break;
      }
    }

    if (idx !== -1) {
      data.domains[data.cd].moderators.splice(idx, 1);
      global.post(global.origin + "/api/domain/moderator/delete", json, function(resp) {
        if (!resp.success) {
          global.globalErrorShow(resp.message);
          return
        }

        global.globalOKShow("Removed!");
      });
    }
  }

} (window.commento, document));
