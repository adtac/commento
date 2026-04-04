(function (global, document) {
  "use strict";

  (document);

  global.adminVueConstruct = function(callback) {
    global.adminApp = new Vue({
      el: "#admin",
      data: {
        owners: [],
        domains: [],
        stats: {},
        plans: [],
        selectedOwner: null,
        selectedOwnerDetails: null,
        view: "stats",
        loaded: false,
      },
    });

    if (callback !== undefined) {
      callback();
    }
  };

  global.adminLoadStats = function() {
    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
    };

    global.post(global.origin + "/api/admin/stats", json, function(resp) {
      if (!resp.success) {
        global.globalError(resp.message);
        return;
      }

      Vue.set(global.adminApp, "stats", resp);
      Vue.set(global.adminApp, "loaded", true);
    });
  };

  global.adminLoadUsers = function() {
    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
    };

    Vue.set(global.adminApp, "view", "users");

    global.post(global.origin + "/api/admin/users/list", json, function(resp) {
      if (!resp.success) {
        global.globalError(resp.message);
        return;
      }

      Vue.set(global.adminApp, "owners", resp.owners);
    });
  };

  global.adminLoadDomains = function() {
    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
    };

    Vue.set(global.adminApp, "view", "domains");

    global.post(global.origin + "/api/admin/domains/list", json, function(resp) {
      if (!resp.success) {
        global.globalError(resp.message);
        return;
      }

      Vue.set(global.adminApp, "domains", resp.domains);
    });
  };

  global.adminShowStats = function() {
    Vue.set(global.adminApp, "view", "stats");
    global.adminLoadStats();
  };

  global.adminViewUser = function(ownerHex) {
    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
      "ownerHex": ownerHex,
    };

    global.post(global.origin + "/api/admin/users/get", json, function(resp) {
      if (!resp.success) {
        global.globalError(resp.message);
        return;
      }

      Vue.set(global.adminApp, "selectedOwnerDetails", resp);
      Vue.set(global.adminApp, "view", "user-detail");
    });
  };

  global.adminChangePlan = function(targetOwnerHex, planId) {
    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
      "targetOwnerHex": targetOwnerHex,
      "planId": planId,
    };

    global.post(global.origin + "/api/admin/subscription/update", json, function(resp) {
      if (!resp.success) {
        global.globalError(resp.message);
        return;
      }

      global.globalOk("Plan updated successfully!");
      global.adminViewUser(targetOwnerHex);
    });
  };

} (window.commento, document));
