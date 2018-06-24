(function (global, document) {
  "use strict";

  (document);

  // Performs a JSON POST request to the given url with the given data and
  // calls the callback function with the JSON response.
  global.post = function(url, json, callback) {
    $.ajax({
      url: url,
      type: "POST",
      data: JSON.stringify(json),
      success: function(data) {
        var resp = JSON.parse(data);
        callback(resp);
      },
    });
  }


  // Performs a GET request and calls the callback function with the JSON
  // response.
  global.get = function(url, callback) {
    $.ajax({
      url: url,
      type: "GET",
      success: function(data) {
        var resp = JSON.parse(data);
        callback(resp);
      },
    });
  }

} (window.commento, document));
