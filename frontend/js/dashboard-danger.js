(function (global, document) {
  "use strict";

  // Opens the danger zone.
  global.dangerOpen = function() {
    $(".view").hide();
    $("#danger-view").show();
  };


  // Deletes a domain.
  global.domainDeleteHandler = function() {
    var data = global.dashboard.$data;

    global.domainDelete(data.domains[data.cd].domain, function(success) {
      if (success) {
        document.location = global.origin + "/dashboard";
      }
    });
  }


  // Clears all comments in a domain.
  global.domainClearHandler = function() {
    var data = global.dashboard.$data;

    global.domainClear(data.domains[data.cd].domain, function(success) {
      if (success) {
        document.location = global.origin + "/dashboard";
      }
    });
  }


  // Freezes a domain.
  global.domainFreezeHandler = function() {
    var data = global.dashboard.$data;

    data.domains[data.cd].state = "frozen"
    global.domainUpdate(data.domains[data.cd])
    document.location.hash = "#modal-close";
  }


  // Unfreezes a domain.
  global.domainUnfreezeHandler = function() {
    var data = global.dashboard.$data;

    data.domains[data.cd].state = "unfrozen"
    global.domainUpdate(data.domains[data.cd])
    document.location.hash = "#modal-close";
  }


} (window.commento, document));
