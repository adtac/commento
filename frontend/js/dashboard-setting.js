(function (global, document) {
  "use strict";

  (document);

  // Sets the vue.js toggle to select and deselect panes visually.
  function settingSelectCSS(id) {
    var data = global.dashboard.$data;
    var settings = data.settings;

    for (var i = 0; i < settings.length; i++) {
      settings[i].selected = settings[i].id === id;
    }
  }


  // Selects a setting.
  global.settingSelect = function(id) {
    var data = global.dashboard.$data;
    var settings = data.settings;

    settingSelectCSS(id);

    $("ul.tabs li").removeClass("current");
    $(".content").removeClass("current");
    $(".original").addClass("current");

    for (var i = 0; i < settings.length; i++) {
      if (id === settings[i].id) {
        settings[i].open();
      }
    }
  };


  // Deselects all settings.
  global.settingDeselectAll = function() {
    var data = global.dashboard.$data;
    var settings = data.settings;

    for (var i = 0; i < settings.length; i++) {
      settings[i].selected = false;
    }
  }

} (window.commento, document));
