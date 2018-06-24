(function (global, document) {
  "use strict";

  (document);

  // Registers a given ID for a fade out after 5 seconds.
  global.registerHide = function(id) {
    setTimeout(function() {
      $(id).fadeOut("fast");
    }, 5000);
  }


  // Shows a global message on the given label ID and registers it for hiding.
  global.showGlobalMessage = function(id, text) {
    global.textSet(id, text);
    global.registerHide(id);
  }


  // Shows a global error message.
  global.globalErrorShow = function(text) {
    global.showGlobalMessage("#global-error", text);
  }


  // Shows a global success message.
  global.globalOKShow = function(text) {
    global.showGlobalMessage("#global-ok", text);
  }

} (window.commento, document));
