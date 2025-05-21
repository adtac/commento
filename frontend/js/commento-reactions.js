  global.commentReaction = function(data) {
    var commentHex = data[0];
    var reactionType = data[1];
    var reactionButton;
    
    // Get the correct reaction button based on type
    switch(reactionType) {
      case "funny":
        reactionButton = $(ID_REACTION_FUNNY + commentHex);
        break;
      case "interesting":
        reactionButton = $(ID_REACTION_INTERESTING + commentHex);
        break;
      case "upsetting":
        reactionButton = $(ID_REACTION_UPSETTING + commentHex);
        break;
      case "sad":
        reactionButton = $(ID_REACTION_SAD + commentHex);
        break;
    }

    // Check if the button is already active (has the active class)
    var isActive = reactionButton.classList.contains("active-reaction");
    
    var json = {
      "commenterToken": commenterTokenGet(),
      "commentHex": commentHex,
      "reactionType": reactionType,
      "remove": isActive
    };

    post(origin + "/api/comment/reaction", json, function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return;
      } else {
        errorHide();
      }

      // Get current count from button text (emoji + number)
      var currentText = reactionButton.innerHTML;
      var currentCount = parseInt(currentText.replace(/[^0-9]/g, '')) || 0;
      var emoji = "";
      
      switch(reactionType) {
        case "funny":
          emoji = "😄 ";
          break;
        case "interesting":
          emoji = "🤔 ";
          break;
        case "upsetting":
          emoji = "😠 ";
          break;
        case "sad":
          emoji = "😢 ";
          break;
      }
      
      if (isActive) {
        // Remove reaction
        classRemove(reactionButton, "active-reaction");
        reactionButton.innerHTML = emoji + (currentCount - 1);
      } else {
        // Add reaction
        classAdd(reactionButton, "active-reaction");
        reactionButton.innerHTML = emoji + (currentCount + 1);
      }
    });
  }
