(function (global, document) {
  "use strict";

  (document);

  // Opens the installation view.
  global.installationOpen = function() {
    var html = "" +
               "<div id=\"commento\"></div>\n" +
               "<script src=\"" + global.cdn + "/js/commento.js\"><\/script>\n" +
               "";

    $("#code-div").text(html);

    $("pre code").each(function(i, block) {
      hljs.highlightBlock(block);
    });

    $(".view").hide();
    $("#installation-view").show();
  };

} (window.commento, document));
