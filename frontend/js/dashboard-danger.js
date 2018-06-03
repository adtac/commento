(function (global, document) {

  // Opens the danger zone.
  global.dangerOpen = function() {
    $(".view").hide();
    $("#danger-view").show();
  };


  // Deletes a domain.
  global.domainDeleteHandler = function() {
    var data = global.dashboard.$data;

    domainDelete(data.domains[data.cd].domain, function(success) {
      if (success)
        document.location = '/dashboard';
    });
  }


  // Freezes a domain.
  global.domainFreezeHandler = function() {
    var data = global.dashboard.$data;

    data.domains[data.cd].state = "frozen"
    domainUpdate(data.domains[data.cd])
    document.location.hash = "#modal-close";
  }


  // Unfreezes a domain.
  global.domainUnfreezeHandler = function() {
    var data = global.dashboard.$data;

    data.domains[data.cd].state = "unfrozen"
    domainUpdate(data.domains[data.cd])
    document.location.hash = "#modal-close";
  }


} (window, document));
