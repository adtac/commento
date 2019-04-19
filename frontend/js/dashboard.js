(function (global, document) {
  "use strict";

  (document);

  // Sets a vue.js field. Short for "vue set".
  function vs(field, value) {
    Vue.set(global.dashboard, field, value);
  }


  global.vs = vs;

  
  // Sets the owner's name in the navbar.
  global.navbarFill = function() {
    $("#owner-name").text(global.owner.name);
  };


  // Constructs the vue.js object
  global.vueConstruct = function(callback) {
    var settings = [
      {
        "id": "installation",
        "text": "Installation Guide",
        "meaning": "Install Commento with HTML",
        "selected": false,
        "open": global.installationOpen,
      },
      {
        "id": "general",
        "text": "General",
        "meaning": "Names, authentication, and export",
        "selected": false,
        "open": global.generalOpen,
      },
      {
        "id": "moderation",
        "text": "Moderation Settings",
        "meaning": "Manage moderators, spam filtering",
        "selected": false,
        "open": global.moderationOpen,
      },
      {
        "id": "statistics",
        "text": "Analytics",
        "meaning": "Anonymous statistics and graphs",
        "selected": false,
        "open": global.statisticsOpen,
      },
      {
        "id": "import",
        "text": "Import Comments",
        "meaning": "Import from a different service",
        "selected": false,
        "open": global.importOpen,
      },
      {
        "id": "danger",
        "text": "Danger Zone",
        "meaning": "Here be dragons",
        "selected": false,
        "open": global.dangerOpen,
      },
    ];

    var reactiveData = {
      // list of panes; mutable because selection information is stored within
      settings: settings,

      // list of domains dynamically loaded; obviously mutable
      domains: [{show: false, viewsLast30Days: global.numberify(0), commentsLast30Days: global.numberify(0), moderators: []}],

      // configured oauth providers that will be filled in after a backend request
      configuredOauths: {},

      // whether or not to show the settings column; mutable because we do not
      // show the column until a domain has been selected
      showSettings: false,

      // currently selected domain index; obviously mutable
      cd: 0, // stands for "current domain"
    };

    global.dashboard = new Vue({
      el: "#dashboard",
      data: reactiveData,
    });

    if (callback !== undefined) {
      callback();
    }
  };

} (window.commento, document));
