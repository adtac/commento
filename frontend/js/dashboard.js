(function (global, document) {

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
        "text": "Installation",
        "meaning": "Install Commento with HTML",
        "selected": false,
        "open": installationOpen,
      },
      {
        "id": "general",
        "text": "General Settings",
        "meaning": "Names, domains and the rest",
        "selected": false,
        "open": generalOpen,
      },
      {
        "id": "moderation",
        "text": "Moderation Settings",
        "meaning": "Approve and delete comments",
        "selected": false,
        "open": moderationOpen,
      },
      {
        "id": "statistics",
        "text": "Statistics",
        "meaning": "Usage and comment statistics",
        "selected": false,
        "open": statisticsOpen,
      },
      {
        "id": "import",
        "text": "Import Comments",
        "meaning": "Import from a different service",
        "selected": false,
        "open": importOpen,
      },
      {
        "id": "danger",
        "text": "Danger Zone",
        "meaning": "Delete or freeze domain",
        "selected": false,
        "open": dangerOpen,
      },
    ];

    var reactiveData = {
      // list of panes; mutable because selection information is stored within
      settings: settings,

      // list of domains dynamically loaded; obviously mutable
      domains: [{show: false, viewsLast30Days: global.numberify(0), commentsLast30Days: global.numberify(0), moderators: []}],

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

    if (callback !== undefined)
      callback();
  };

} (window, document));
