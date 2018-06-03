(function (global, document) {

  // Opens the installation view.
  global.installationOpen = function() {
    var data = global.dashboard.$data;

    var html = '' +
               '<div id="commento"></div>\n' +
               '<script src="' + window.cdn + '/js/commento.js"><\/script>\n' +
               '';

    $("#code-div").text(html);

    $('pre code').each(function(i, block) {
      hljs.highlightBlock(block);
    });

    $(".view").hide();
    $("#installation-view").show();
  };

} (window, document));
