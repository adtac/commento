(function (global, document) {
  "use strict";

  // Gets a GET parameter in the current URL.
  global.paramGet = function(param) {
    var pageURL = decodeURIComponent(window.location.search.substring(1));
    var urlVariables = pageURL.split("&");

    for (var i = 0; i < urlVariables.length; i++) {
      var paramURL = urlVariables[i].split("=");
      if (paramURL[0] === param) {
        return paramURL[1] === undefined ? true : paramURL[1];
      }
    }

    return null;
  }


  // Sets the disabled attribute in a button.
  global.buttonDisable = function(id) {
    var el = $(id);

    el.attr("disabled", true);
  }


  // Unsets the disabled attribute in a button.
  global.buttonEnable = function(id) {
    var el = $(id);

    el.attr("disabled", false);
  }

  // Sets the text on the given label ID.
  global.textSet = function(id, text) {
    var el = $(id);

    el.show();
    el.text(text);
  }


  // Given an array of input IDs, this function calls a callback function with
  // the first unfilled ID.
  global.unfilledMark = function(fields, callback) {
    var allOk = true;

    for (var i = 0; i < fields.length; i++) {
      var el = $(fields[i]);
      if (el.val() === "") {
        callback(el);
      }
    }

    return allOk;
  }


  // Gets the value of a cookie.
  global.cookieGet = function(name) {
    var c = "; " + document.cookie;
    var x = c.split("; " + name + "=");
    if (x.length === 2) {
      return x.pop().split(";").shift();
    }
  };


  // Sets the value of a cookie.
  global.cookieSet = function(name, value) {
    var expires = "";
    var date = new Date();
    date.setTime(date.getTime() + (365*24*60*60*1000));
    expires = "; expires=" + date.toUTCString();

    var cookieString = name + "=" + value + expires + "; path=/";
    if (/^https:\/\//i.test(global.origin)) {
      cookieString += "; secure";
    }

    document.cookie = cookieString;
  }


  // Deletes a cookie.
  global.cookieDelete = function(name) {
    document.cookie = name + "=;expires=Thu, 01 Jan 1970 00:00:01 GMT;";
  }


  // Converts a date in the past to a human-friendly duration relative to now.
  global.timeSince = function(date) {
    var seconds = Math.floor((new Date() - date) / 1000);
    var interval = Math.floor(seconds / 31536000);

    if (interval > 1) {
      return interval + " years ago";
    }

    interval = Math.floor(seconds / 2592000);
    if (interval > 1) {
      return interval + " months ago";
    }

    interval = Math.floor(seconds / 86400);
    if (interval > 1) {
      return interval + " days ago";
    }

    interval = Math.floor(seconds / 3600);
    if (interval > 1) {
      return interval + " hours ago";
    }

    interval = Math.floor(seconds / 60);
    if (interval > 1) {
      return interval + " minutes ago";
    }

    if (seconds > 5) {
      return Math.floor(seconds) + " seconds ago";
    } else {
      return "just now";
    }
  }

} (window.commento, document));
