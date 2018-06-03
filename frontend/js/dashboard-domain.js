(function (global, document) {

  // Selects a domain.
  global.domainSelect = function(domain) {
    var data = global.dashboard.$data;
    var domains = data.domains;

    for (var i = 0; i < domains.length; i++) {
      if (domains[i].domain == domain) {
        vs("frozen", domains[i].state == "frozen");
        domains[i].selected = true;
        data.cd = i;
        data.importedComments = domains[i].importedComments;
      }
      else
        domains[i].selected = false;
    }

    data.showSettings = true;

    settingDeselectAll();
    $(".view").hide();
  };


  // Deselects all domains.
  global.domainDeselectAll = function() {
    var data = global.dashboard.$data;
    var domains = data.domains;

    for (var i = 0; i < domains.length; i++)
      domains[i].selected = false;
  }


  // Creates a new domain.
  global.domainNewHandler = function() {
    var json = {
      session: global.cookieGet("session"),
      name: $("#new-domain-name").val(),
      domain: $("#new-domain-domain").val(),
    }

    global.buttonDisable("#add-site-button");
    global.post(global.origin + "/api/domain/new", json, function(resp) {
      global.buttonEnable("#add-site-button");

      $("#new-domain-name").val("");
      $("#new-domain-domain").val("");
      document.location.hash = "#modal-close";

      if (!resp.success) {
        global.globalErrorShow(resp.message);
        return;
      }

      global.domainRefresh(function() {
        global.domainSelect(resp.domain);
        global.domainDeselectAll();
        global.settingSelect("installation");
      });
    });
  }


  // Refreshes the list of domains.
  global.domainRefresh = function(callback) {
    var json = {
      session: global.cookieGet("session"),
    };

    global.post(global.origin + "/api/domain/list", json, function(resp) {
      if (!resp.success) {
        global.globalErrorShow(resp.message);
        return;
      }

      resp.domains = resp.domains.sort(function(a, b) {
          var x = a.creationDate; var y = b.creationDate;
          return ((x < y) ? -1 : ((x > y) ? 1 : 0));
      });

      for (var i = 0; i < resp.domains.length; i++) {
        resp.domains[i].show = true;
        resp.domains[i].selected = false;

        resp.domains[i].origName = resp.domains[i].name;
        resp.domains[i].origDomain = resp.domains[i].domain;

        resp.domains[i].viewsLast30Days = global.numberify(0);
        resp.domains[i].commentsLast30Days = global.numberify(0);
        
        for (var j = 0; j < resp.domains[i].moderators.length; j++) {
          resp.domains[i].moderators[j].timeAgo = global.timeSince(
            Date.parse(resp.domains[i].moderators[j].addDate));
        }
      }

      global.vs("domains", resp.domains);

      if (callback !== undefined)
        callback();
    });
  };


  // Updates a domain with the backend.
  global.domainUpdate = function(domain, callback) {
    var json = {
      session: global.cookieGet("session"),
      domain: domain,
    };

    global.post(global.origin + "/api/domain/update", json, function(resp) {
      if (callback !== undefined)
        callback(resp.success);

      if (!resp.success) {
        global.globalErrorShow(resp.message);
        return;
      }
    });
  }


  // Deletes a domain.
  global.domainDelete = function(domain, callback) {
    var json = {
      session: global.cookieGet("session"),
      domain: domain,
    };

    global.post(global.origin + "/api/domain/delete", json, function(resp) {
      if (!resp.success) {
        global.globalErrorShow(resp.message);
        return;
      }

      if (callback !== undefined)
        callback(resp.success);
    });
  }

} (window, document));
