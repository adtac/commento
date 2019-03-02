(function(global, document) {
  "use strict";

  var origin = "[[[.Origin]]]";

  function post(url, data, callback) {
    var xmlDoc = new XMLHttpRequest();

    xmlDoc.open("POST", url, true);
    xmlDoc.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xmlDoc.onload = function() {
      callback(JSON.parse(xmlDoc.response));
    };

    xmlDoc.send(JSON.stringify(data));
  }

  function main() {
    var paths = [];
    var doms = [];
    var as = document.getElementsByTagName("a");
    for (var i = 0; i < as.length; i++) {
      var href = as[i].href;
      if (href === undefined) {
        return;
      }

      href = href.replace(/^.*\/\/[^\/]+/, "");

      if (href.endsWith("#commento")) {
        var path = href.substr(0, href.indexOf("#commento"));
        if (path.startsWith(parent.location.host)) {
          path = path.substr(parent.location.host.length);
        }

        paths.push(path);
        doms.push(as[i]);
      }
    }

    var json = {
      "domain": parent.location.host,
      "paths": paths,
    };

    post(origin + "/api/comment/count", json, function(resp) {
      if (!resp.success) {
        console.log("[commento] error: " + resp.message);
        return;
      }

      for (var i = 0; i < paths.length; i++) {
        var count = 0;
        if (paths[i] in resp.commentCounts) {
          count = resp.commentCounts[paths[i]];
        }

        doms[i].innerText = count + " " + (count === 1 ? "comment" : "comments");
      }
    });
  }

  var initted = false;

  function init() {
    if (initted) {
      return;
    }
    initted = true;

    main(undefined);
  }

  var readyLoad = function() {
    var readyState = document.readyState;

    if (readyState === "loading") {
      // The document is still loading. The div we need to fill might not have
      // been parsed yet, so let's wait and retry when the readyState changes.
      // If there is more than one state change, we aren't affected because we
      // have a double-call protection in init().
      document.addEventListener("readystatechange", readyLoad);
    } else if (readyState === "interactive") {
      // The document has been parsed and DOM objects are now accessible. While
      // JS, CSS, and images are still loading, we don't need to wait.
      init();
    } else if (readyState === "complete") {
      // The page has fully loaded (including JS, CSS, and images). From our
      // point of view, this is practically no different from interactive.
      init();
    }
  };

  readyLoad();

}(window, document));
