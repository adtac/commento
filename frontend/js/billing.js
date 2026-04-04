(function (global, document) {
  "use strict";

  (document);

  global.billingVueConstruct = function(callback) {
    global.billingApp = new Vue({
      el: "#billing",
      data: {
        subscription: {},
        plan: {},
        usage: {},
        plans: [],
        domainCount: 0,
        loaded: false,
      },
    });

    if (callback !== undefined) {
      callback();
    }
  };

  global.billingLoad = function() {
    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
    };

    global.post(global.origin + "/api/subscription/get", json, function(resp) {
      if (!resp.success) {
        global.globalError(resp.message);
        return;
      }

      Vue.set(global.billingApp, "subscription", resp.subscription);
      Vue.set(global.billingApp, "plan", resp.plan);
    });

    global.post(global.origin + "/api/owner/usage", json, function(resp) {
      if (!resp.success) {
        global.globalError(resp.message);
        return;
      }

      Vue.set(global.billingApp, "usage", resp.usage);
      Vue.set(global.billingApp, "domainCount", resp.domainCount);
    });

    global.post(global.origin + "/api/plan/list", {}, function(resp) {
      if (!resp.success) {
        global.globalError(resp.message);
        return;
      }

      Vue.set(global.billingApp, "plans", resp.plans);
      Vue.set(global.billingApp, "loaded", true);
    });
  };

  global.planUpgrade = function(planId) {
    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
      "planId": planId,
    };

    global.post(global.origin + "/api/subscription/update", json, function(resp) {
      if (!resp.success) {
        global.globalError(resp.message);
        return;
      }

      global.globalOk("Plan updated successfully!");
      global.billingLoad();
    });
  };

} (window.commento, document));
