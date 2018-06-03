(function (global, document) {

  // Opens the general settings window.
  global.generalOpen = function() {
    $(".view").hide();
    $("#general-view").show();
  };

  global.generalSaveHandler = function() {
    var data = global.dashboard.$data;

    global.buttonDisable("#save-general-button");
    global.domainUpdate(data.domains[data.cd], function() {
      global.globalOKShow("Settings saved!");
      global.buttonEnable("#save-general-button");
    });
  };

} (window, document));
