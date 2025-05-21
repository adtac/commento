  global.commentSpoiler = function(commentHex) {
    var text = $(ID_TEXT + commentHex);
    var spoilerButton = $(ID_SPOILER + commentHex);
    var isSpoiler = !text.classList.contains("spoiler");

    var json = {
      "commenterToken": commenterTokenGet(),
      "commentHex": commentHex,
      "spoiler": isSpoiler
    };

    post(origin + "/api/comment/spoiler", json, function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return;
      } else {
        errorHide();
      }

      if (isSpoiler) {
        classAdd(text, "spoiler");
        spoilerButton.title = "Remove spoiler tag";
      } else {
        classRemove(text, "spoiler");
        spoilerButton.title = "Mark as spoiler";
      }
    });
  }
