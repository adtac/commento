(function (global, document) {
  "use strict";

  (document);

  global.numberify = function(x) {
    if (x === 0) {
      return {"zeros": "000", "num": "", "units": ""}
    }

    if (x < 10) {
      return {"zeros": "00", "num": x, "units": ""}
    }

    if (x < 100) {
      return {"zeros": "0", "num": x, "units": ""}
    }

    if (x < 1000) {
      return {"zeros": "", "num": x, "units": ""}
    }

    var res;

    if (x < 1000000) {
      res = global.numberify((x/1000).toFixed(0))
      res.units = "K"
    } else if (x < 1000000000) {
      res = global.numberify((x/1000000).toFixed(0))
      res.units = "M"
    } else if (x < 1000000000000) {
      res = global.numberify((x/1000000000).toFixed(0))
      res.units = "B"
    }

    if (res.num*10 % 10 === 0) {
      res.num = Math.ceil(res.num);
    }
    
    return res;
  }

  global.statisticsOpen = function() {
    var data = global.dashboard.$data;
    
    var json = {
      "ownerToken": global.cookieGet("commentoOwnerToken"),
      "domain": data.domains[data.cd].domain,
    }

    $(".view").hide();
    global.post(global.origin + "/api/domain/statistics", json, function(resp) {
      $("#statistics-view").show();

      if (!resp.success) {
        global.globalErrorShow(resp.message);
        return;
      }

      var options = {
        showPoint: false,
        axisY: {
          onlyInteger: true,
          showGrid: false,
        },
        axisX: {
          showGrid: false,
        },
        showArea: true,
      };

      var views;
      var comments;

      views = resp.viewsLast30Days;
      // views = [0, 1, 4, 16, 14, 12, 10, 25, 13, 5, 20, 25, 12, 57, 46, 64, 4, 36, 7, 80, 43, 86, 121, 6, 74, 94, 83, 73, 140, 89, 25];

      comments = resp.commentsLast30Days;
      // comments = [0, 0, 1, 2, 3, 3, 4, 5, 7, 8, 5, 9, 9, 5, 6, 7, 8, 3, 1, 16, 3, 10, 8, 5, 12, 5, 4, 8, 4, 23, 19];

      var labels = new Array();
      for (var i = 0; i < views.length; i++) {
        if ((views.length-i) % 7 === 0) {
          var x = (views.length-i)/7;
          labels.push(x + " week" + (x > 1 ? "s" : "") + " ago");
        } else {
          labels.push("");
        }
      }

      new Chartist.Line("#views-graph", {
        labels: labels,
        series: [views],
      }, options);

      new Chartist.Line("#comments-graph", {
        labels: labels,
        series: [comments],
      }, options);

      data.domains[data.cd].viewsLast30Days = global.numberify(views.reduce(function(a, b) {
        return a + b; 
      }, 0));
      data.domains[data.cd].commentsLast30Days = global.numberify(comments.reduce(function(a, b) {
        return a + b; 
      }, 0));
    });
  }

} (window.commento, document));
