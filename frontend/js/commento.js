(function(global, document) {
  'use strict';

  // Do not use other files like utils.js and http.js in the Makefile to build
  // commento.js for the following reasons:
  //   - We don't use jQuery in the actual JavaScript payload because we need
  //     to be lightweight.
  //   - They pollute the global/window namespace (with global.post, etc.).
  //     That's NOT fine when we expect them source our JavaScript. For example,
  //     the user may have their own window.post defined. We don't want to
  //     override that.


  var origin = global.commento_origin;
  var cdn = global.commento_cdn;
  var root = null;
  var isAuthenticated = false;
  var comments = [];
  var commenters = [];
  var requireIdentification = true;
  var requireModeration = true;
  var isModerator = false;
  var isFrozen = false;
  var chosenAnonymous = false;
  var shownSubmitButton = {"root": false};
  var shownReply = {};


  function $(id) {
    return document.getElementById(id);
  }


  function dataGet(el, key) {
    return el.dataset[key];
  }


  function dataSet(el, key, data) {
    el.dataset[key] = data;
  }


  function append(root, el) {
    root.appendChild(el);
  }


  function prepend(root, el) {
    root.prepend(el);
  }


  function classAdd(el, cls) {
    el.classList.add("commento-" + cls);
  }


  function classRemove(el, cls) {
    el.classList.remove("commento-" + cls);
  }


  function create(el) {
    return document.createElement(el);
  }


  function remove(el) {
    el.parentNode.removeChild(el);
  }


  function attr(node, a, value) {
    node.setAttribute(a, value);
  }


  function post(url, data, callback) {
    var xmlDoc = new XMLHttpRequest();

    xmlDoc.open("POST", url, true);
    xmlDoc.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xmlDoc.onload = function() {
      callback(JSON.parse(xmlDoc.response));
    };

    xmlDoc.send(JSON.stringify(data));
  }


  function get(url, callback) {
    var xmlDoc = new XMLHttpRequest();

    xmlDoc.open('GET', url, true);
    xmlDoc.onload = function() {
      callback(JSON.parse(xmlDoc.response));
    };

    xmlDoc.send(null);
  }


  function call(callback) {
    if (typeof(callback) == "function")
      callback();
  }


  function rootSpinnerShow() {
    var spinner = create("div");

    classAdd(spinner, "loading");

    append(root, spinner);
  }


  function cookieGet(name) {
    var c = "; " + document.cookie;
    var x = c.split("; " + name + "=");
    if (x.length == 2)
      return x.pop().split(";").shift();
  }


  function cookieSet(name, value) {
    var expires = "";
    var date = new Date();
    date.setTime(date.getTime() + (365*24*60*60*1000));
    expires = "; expires=" + date.toUTCString();

    document.cookie = name + "=" + value + expires + "; path=/";
  }

  function sessionGet() {
    var session = cookieGet("session");
    if (session === undefined)
      return "anonymous";

    return session;
  }

  global.logout = function() {
    cookieSet("session", "anonymous");
    refreshAll();
  }

  function selfGet(callback) {
    var session = sessionGet();
    if (session == "anonymous") {
      isAuthenticated = false;
      call(callback);
      return;
    }

    var json = {
      session: sessionGet(),
    };

    post(origin + "/api/commenter/self", json, function(resp) {
      console.log(resp);
      if (!resp.success) {
        cookieSet("session", "anonymous");
        call(callback);
        return;
      }

      var loggedContainer = create("div");
      var loggedInAs = create("div");
      var name = create("a");
      var photo = create("img");
      var logout = create("div");

      classAdd(loggedContainer, "logged-container");
      classAdd(loggedInAs, "logged-in-as");
      classAdd(name, "name");
      classAdd(photo, "photo");
      classAdd(logout, "logout");

      name.innerText = resp.commenter.name;
      logout.innerText = "Logout";

      attr(name, "href", resp.commenter.link);
      if (resp.commenter.provider == "google") {
        attr(photo, "src", resp.commenter.photo + "?sz=50");
      } else {
        attr(photo, "src", resp.commenter.photo);
      }
      attr(logout, "onclick", "logout()");

      append(loggedInAs, photo);
      append(loggedInAs, name);
      append(loggedContainer, loggedInAs);
      append(loggedContainer, logout);
      append(root, loggedContainer);

      isAuthenticated = true;

      call(callback);
    });
  }

  function cssLoad(file) {
    var link = create("link");
    var head = document.getElementsByTagName('head')[0];

    link.type = "text/css";
    attr(link, "href", file);
    attr(link, "rel", "stylesheet");

    append(head, link);
  }

  function jsLoad(file, ready) {
    var script = document.createElement("script");
    var loaded = false;

    script.type = "application/javascript";
    script.src = file;
    script.async = true;
    script.onreadysessionchange = script.onload = function() {
      if (!loaded &&
        (!this.readySession ||
          this.readySession === "loaded" ||
          this.readySession === "complete"))
      {
        ready();
      }

      loaded = true;
      script.onload = script.onreadysessionchange = null;
    };

    append(document.body, script);
  }

  function footerLoad() {
    var footer = create("div");
    var aContainer = create("div");
    var a = create("a");
    var img = create("img");
    var text = create("span");

    classAdd(footer, "footer");
    classAdd(aContainer, "logo-container");
    classAdd(a, "logo");
    classAdd(img, "logo-svg");
    classAdd(text, "logo-text");

    attr(a, "href", "https://commento.io");
    attr(a, "target", "_blank");
    attr(img, "src", cdn + "/images/logo.svg");

    text.innerText = "Powered by Commento";

    append(a, img);
    append(a, text);
    append(aContainer, a);
    append(footer, aContainer);
    append(root, footer);
  }

  function commentsGet(callback) {
    var json = {
      session: sessionGet(),
      domain: location.host,
      path: location.pathname,
    };

    post(origin + "/api/comment/list", json, function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return;
      }

      requireModeration = resp.requireModeration;
      requireIdentification = resp.requireIdentification;
      isModerator = resp.isModerator;
      isFrozen = resp.isFrozen;
      comments = resp.comments;
      commenters = resp.commenters;

      cssLoad(cdn + "/css/commento.css");

      call(callback);
    });
  }

  function errorShow(text) {
    var el = $(ID_ERROR);

    el.innerText = text;

    attr(el, "style", "display: block;");
  }

  function errorElementCreate() {
    var el = create("div");

    el.id = ID_ERROR;

    classAdd(el, "error-box");
    attr(el, "style", "display: none;");

    append(root, el);
  }

  function autoExpander(el) {
    return function() {
      el.style.height = "";
      el.style.height = Math.min(Math.max(el.scrollHeight, 75), 400) + "px";
    }
  };

  function isMobile() {
    var mobile = false;
    if (/(android|bb\d+|meego).+mobile|avantgo|bada\/|blackberry|blazer|compal|elaine|fennec|hiptop|iemobile|ip(hone|od)|ipad|iris|kindle|Android|Silk|lge |maemo|midp|mmp|netfront|opera m(ob|in)i|palm( os)?|phone|p(ixi|re)\/|plucker|pocket|psp|series(4|6)0|symbian|treo|up\.(browser|link)|vodafone|wap|windows (ce|phone)|xda|xiino/i.test(navigator.userAgent) 
        || /1207|6310|6590|3gso|4thp|50[1-6]i|770s|802s|a wa|abac|ac(er|oo|s\-)|ai(ko|rn)|al(av|ca|co)|amoi|an(ex|ny|yw)|aptu|ar(ch|go)|as(te|us)|attw|au(di|\-m|r |s )|avan|be(ck|ll|nq)|bi(lb|rd)|bl(ac|az)|br(e|v)w|bumb|bw\-(n|u)|c55\/|capi|ccwa|cdm\-|cell|chtm|cldc|cmd\-|co(mp|nd)|craw|da(it|ll|ng)|dbte|dc\-s|devi|dica|dmob|do(c|p)o|ds(12|\-d)|el(49|ai)|em(l2|ul)|er(ic|k0)|esl8|ez([4-7]0|os|wa|ze)|fetc|fly(\-|_)|g1 u|g560|gene|gf\-5|g\-mo|go(\.w|od)|gr(ad|un)|haie|hcit|hd\-(m|p|t)|hei\-|hi(pt|ta)|hp( i|ip)|hs\-c|ht(c(\-| |_|a|g|p|s|t)|tp)|hu(aw|tc)|i\-(20|go|ma)|i230|iac( |\-|\/)|ibro|idea|ig01|ikom|im1k|inno|ipaq|iris|ja(t|v)a|jbro|jemu|jigs|kddi|keji|kgt( |\/)|klon|kpt |kwc\-|kyo(c|k)|le(no|xi)|lg( g|\/(k|l|u)|50|54|\-[a-w])|libw|lynx|m1\-w|m3ga|m50\/|ma(te|ui|xo)|mc(01|21|ca)|m\-cr|me(rc|ri)|mi(o8|oa|ts)|mmef|mo(01|02|bi|de|do|t(\-| |o|v)|zz)|mt(50|p1|v )|mwbp|mywa|n10[0-2]|n20[2-3]|n30(0|2)|n50(0|2|5)|n7(0(0|1)|10)|ne((c|m)\-|on|tf|wf|wg|wt)|nok(6|i)|nzph|o2im|op(ti|wv)|oran|owg1|p800|pan(a|d|t)|pdxg|pg(13|\-([1-8]|c))|phil|pire|pl(ay|uc)|pn\-2|po(ck|rt|se)|prox|psio|pt\-g|qa\-a|qc(07|12|21|32|60|\-[2-7]|i\-)|qtek|r380|r600|raks|rim9|ro(ve|zo)|s55\/|sa(ge|ma|mm|ms|ny|va)|sc(01|h\-|oo|p\-)|sdk\/|se(c(\-|0|1)|47|mc|nd|ri)|sgh\-|shar|sie(\-|m)|sk\-0|sl(45|id)|sm(al|ar|b3|it|t5)|so(ft|ny)|sp(01|h\-|v\-|v )|sy(01|mb)|t2(18|50)|t6(00|10|18)|ta(gt|lk)|tcl\-|tdg\-|tel(i|m)|tim\-|t\-mo|to(pl|sh)|ts(70|m\-|m3|m5)|tx\-9|up(\.b|g1|si)|utst|v400|v750|veri|vi(rg|te)|vk(40|5[0-3]|\-v)|vm40|voda|vulc|vx(52|53|60|61|70|80|81|83|85|98)|w3c(\-| )|webc|whit|wi(g |nc|nw)|wmlb|wonu|x700|yas\-|your|zeto|zte\-/i.test(navigator.userAgent.substr(0,4)))
      mobile = true;

    return mobile;
  }

  function textareaCreate(id) {
    var textareaSuperContainer = create("div");
    var textareaContainer = create("div");
    var textarea = create("textarea");

    textareaSuperContainer.id = ID_SUPER_CONTAINER + id;
    textareaContainer.id = ID_TEXTAREA_CONTAINER + id;
    textarea.id = ID_TEXTAREA + id;

    classAdd(textareaContainer, "textarea-container");

    if (!isAuthenticated && !chosenAnonymous) {
      var buttonsContainer = create("div");

      if (!isMobile()) {
        classAdd(buttonsContainer, "buttons-container");
        classAdd(textarea, "blurred-textarea");
      } else {
        classAdd(textarea, "hidden");
        classAdd(buttonsContainer, "mobile-buttons-container");
        classAdd(buttonsContainer, "opaque");
      }

      var oauths = ["google"];
      if (!requireIdentification)
        oauths.push("anonymous");

      for (var i = 0; i < oauths.length; i++) {
        var oauthButton = create("button");

        classAdd(oauthButton, "button");
        classAdd(oauthButton, oauths[i] + "-button");

        if (isMobile())
          classAdd(oauthButton, "opaque");

        attr(oauthButton, "onclick", "commentoAuth('" + oauths[i] + "', '" + id + "')");

        oauthButton.innerText = oauths[i][0].toUpperCase() + oauths[i].slice(1);

        append(buttonsContainer, oauthButton);
      }

      attr(textarea, "disabled", true);

      append(textareaContainer, buttonsContainer);
    }

    attr(textarea, "placeholder", "Join the discussion!");
    attr(textarea, "onclick", "showSubmitButton('" + id + "')");

    textarea.oninput = autoExpander(textarea);

    append(textareaContainer, textarea);
    append(textareaSuperContainer, textareaContainer);

    return textareaSuperContainer;
  }

  function rootCreate(callback) {
    var commentsArea = create("div");

    commentsArea.id = ID_COMMENTS_AREA;

    classAdd(commentsArea, "comments");

    commentsArea.innerHTML = "";

    append(root, textareaCreate("root"));
    append(root, commentsArea);

    call(callback);
  }

  function messageCreate(text) {
    var msg = create("div");

    classAdd(msg, "moderation-notice");

    msg.innerText = text;

    return msg;
  }

  global.postComment = function(id) {
    var textarea = $(ID_TEXTAREA + id);

    var comment = textarea.value;

    if (comment == "") {
      classAdd(textarea, "red-border");
      return;
    }
    else {
      classRemove(textarea, "red-border");
    }

    var json = {
      "session": sessionGet(),
      "domain": location.host,
      "path": location.pathname,
      "parentHex": id,
      "markdown": comment,
    };

    post(origin + "/api/comment/new", json, function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return;
      }

      $(ID_TEXTAREA + id).value = "";

      commentsGet(function() {
        $(ID_COMMENTS_AREA).innerHTML = "";
        commentsRender();

        if (requireModeration && !isModerator) {
          if (id == "root") {
            var body = $(ID_SUPER_CONTAINER + id);
            prepend(body, messageCreate("Your comment is under moderation."));
          } else {
            var body = $(ID_BODY + id);
            append(body, messageCreate("Your comment is under moderation."));
          }
        }
      });
    });
  }

  function colorGet(name) {
    var colors = [
      // some visually distincy
      "#35b2de", // some kind of teal/cyan
      "#62cd0a", // fresh lemon green
      "#383838", // shade of gray
      "#e4a90f", // comfy yellow
      "#f80707", // sharp red
      "#f0479c", // bright pink
    ];

    var total = 0;
    for (var i = 0; i < name.length; i++)
      total += name.charCodeAt(i);
    var color = colors[total % colors.length];

    return color;
  }

  function timeDifference(current, previous) { // thanks stackoverflow
    var msJustNow = 5000;
    var msPerMinute = 60000;
    var msPerHour = 3600000;
    var msPerDay = 86400000;
    var msPerMonth = 2592000000;
    var msPerYear = 946080000000;

    var elapsed = current - previous;

    if (elapsed < msJustNow) {
     return 'just now';
    }
    else if (elapsed < msPerMinute) {
      return Math.round(elapsed/1000) + ' seconds ago';
    }
    else if (elapsed < msPerHour) {
      return Math.round(elapsed/msPerMinute) + ' minutes ago';
    }
    else if (elapsed < msPerDay ) {
      return Math.round(elapsed/msPerHour ) + ' hours ago';
    }
    else if (elapsed < msPerMonth) {
      return Math.round(elapsed/msPerDay) + ' days ago';
    }
    else if (elapsed < msPerYear) {
      return Math.round(elapsed/msPerMonth) + ' months ago';
    }
    else {
      return Math.round(elapsed/msPerYear ) + ' years ago';
    }
  }

  function scorify(score) {
    if (score != 1)
      return score + " points";
    else
      return score + " point";
  }

  var ID_ERROR = "commento-error";
  var ID_COMMENTS_AREA = "commento-comments-area";
  var ID_SUPER_CONTAINER = "commento-textarea-super-container-";
  var ID_TEXTAREA_CONTAINER = "commento-textarea-container-";
  var ID_TEXTAREA = "commento-textarea-";
  var ID_CARD = "commento-comment-card-";
  var ID_BODY = "commento-comment-body-";
  var ID_SUBTITLE = "commento-comment-subtitle-";
  var ID_SCORE = "commento-comment-score-";
  var ID_OPTIONS = "commento-comment-options-";
  var ID_EDIT = "commento-comment-edit-";
  var ID_REPLY = "commento-comment-reply-";
  var ID_COLLAPSE = "commento-comment-collapse-";
  var ID_UPVOTE = "commento-comment-upvote-";
  var ID_DOWNVOTE = "commento-comment-downvote-";
  var ID_APPROVE = "commento-comment-approve-";
  var ID_REMOVE = "commento-comment-remove-";
  var ID_CONTENTS = "commento-comment-contents-";
  var ID_SUBMIT_BUTTON = "commento-submit-button-";

  function commentsRecurse(parentMap, parentHex) {
    var cur = parentMap[parentHex];
    if (!cur || !cur.length) {
      return null;
    }

    var cards = create("div");
    cur.forEach(function(comment) {
      var commenter = commenters[comment.commenterHex];
      var avatar;
      var card = create("div");
      var header = create("div");
      var subtitle = create("div");
      var score = create("div");
      var body = create("div");
      var options = create("div");
      var edit = create("button");
      var reply = create("button");
      var collapse = create("button");
      var upvote = create("div");
      var downvote = create("div");
      var approve = create("button");
      var remove = create("button");
      var children = commentsRecurse(parentMap, comment.commentHex);
      var contents = create("div");
      var color = colorGet(commenter.name);
      var name;
      if (commenter.link != "undefined")
        name = create("a");
      else
        name = create("div");

      card.id = ID_CARD + comment.commentHex;
      body.id = ID_BODY + comment.commentHex;
      subtitle.id = ID_SUBTITLE + comment.commentHex;
      score.id = ID_SCORE + comment.commentHex;
      options.id = ID_OPTIONS + comment.commentHex;
      edit.id = ID_EDIT + comment.commentHex;
      reply.id = ID_REPLY + comment.commentHex;
      collapse.id = ID_COLLAPSE + comment.commentHex;
      upvote.id = ID_UPVOTE + comment.commentHex;
      downvote.id = ID_DOWNVOTE + comment.commentHex;
      approve.id = ID_APPROVE + comment.commentHex;
      remove.id = ID_REMOVE + comment.commentHex;
      contents.id = ID_CONTENTS + comment.commentHex;

      collapse.title = "Collapse";
      upvote.title = "Upvote";
      downvote.title = "Downvote";
      edit.title = "Edit";
      reply.title = "Reply";
      approve.title = "Approve";
      remove.title = "Remove";

      card.style["borderLeft"] = "2px solid " + color;
      name.innerText = commenter.name;
      body.innerHTML = comment.html;
      subtitle.innerHTML = timeDifference((new Date()).getTime(), Date.parse(comment.creationDate));
      score.innerText = scorify(comment.score);

      if (commenter.photo == "undefined") {
        avatar = create("div");
        avatar.style["background"] = color;
        avatar.style["boxShadow"] = "0px 0px 0px 2px " + color;
        avatar.innerHTML = commenter.name[0].toUpperCase();
        classAdd(avatar, "avatar");
      } else {
        avatar = create("img");
        if (commenter.provider == "google") {
          attr(avatar, "src", commenter.photo + "?sz=50");
        } else {
          attr(avatar, "src", commenter.photo);
        }
        classAdd(avatar, "avatar-img");
      }

      classAdd(card, "card");
      if (isModerator && comment.state == "unapproved")
        classAdd(card, "dark-card");
      classAdd(header, "header");
      classAdd(name, "name");
      classAdd(subtitle, "subtitle");
      classAdd(score, "score");
      classAdd(body, "body");
      classAdd(options, "options");
      classAdd(edit, "option-button");
      classAdd(edit, "option-edit");
      classAdd(reply, "option-button");
      classAdd(reply, "option-reply");
      classAdd(collapse, "option-button");
      classAdd(collapse, "option-collapse");
      classAdd(upvote, "option-button");
      classAdd(upvote, "option-upvote");
      classAdd(downvote, "option-button");
      classAdd(downvote, "option-downvote");
      classAdd(approve, "option-button");
      classAdd(approve, "option-approve");
      classAdd(remove, "option-button");
      classAdd(remove, "option-remove");

      if (isAuthenticated) {
        if (comment.direction > 0)
          classAdd(upvote, "upvoted");
        else if (comment.direction < 0)
          classAdd(downvote, "downvoted");
      }

      attr(edit, "onclick", "startEdit('" + comment.commentHex + "')");
      attr(reply, "onclick", "replyShow('" + comment.commentHex + "')");
      attr(collapse, "onclick", "commentCollapse('" + comment.commentHex + "')");
      attr(approve, "onclick", "commentApprove('" + comment.commentHex + "')");
      attr(remove, "onclick", "commentDelete('" + comment.commentHex + "')");

      if (isAuthenticated) {
        if (comment.direction > 0) {
          attr(upvote, "onclick", "vote('" + comment.commentHex + "', 1, 0)");
          attr(downvote, "onclick", "vote('" + comment.commentHex + "', 1, -1)");
        }
        else if (comment.direction < 0) {
          attr(upvote, "onclick", "vote('" + comment.commentHex + "', -1, 1)");
          attr(downvote, "onclick", "vote('" + comment.commentHex + "', -1, 0)");
        }
        else {
          attr(upvote, "onclick", "vote('" + comment.commentHex + "', 0, 1)");
          attr(downvote, "onclick", "vote('" + comment.commentHex + "', 0, -1)");
        }
      } else if (!chosenAnonymous) {
        attr(upvote, "onclick", "replyShow('" + comment.commentHex + "')");
      }

      if (isAuthenticated) {
        if (isModerator) {
          if (comment.state == "unapproved")
            attr(options, "style", "width: 192px;");
          else
            attr(options, "style", "width: 160px;");
        }
        else
          attr(options, "style", "width: 128px;");
      }
      else
        attr(options, "style", "width: 32px;");

      if (commenter.link != "undefined")
        attr(name, "href", commenter.link);

      if (false) // replace when edit is implemented
        append(options, edit);

      if (isAuthenticated) {
        append(options, upvote);
        append(options, downvote);
        append(options, reply);
      }

      if (isModerator) {
        if (comment.state == "unapproved")
          append(options, approve);
        append(options, remove);
      }

      append(options, collapse);
      append(subtitle, score);
      append(header, options);
      append(header, avatar);
      append(header, name);
      append(header, subtitle);
      append(contents, body);

      if (children) {
        classAdd(children, "body");
        append(contents, children);
      }

      append(card, header);
      append(card, contents);
      append(cards, card);

      shownSubmitButton[comment.commentHex] = false;
    });

    return cards;
  }


  global.commentApprove = function(commentHex) {
    var json = {
      "session": sessionGet(),
      "commentHex": commentHex,
    }

    post(origin + "/api/comment/approve", json, function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return
      }

      var card = $(ID_CARD + commentHex);
      var options = $(ID_OPTIONS + commentHex);
      var tick = $(ID_APPROVE + commentHex);

      classRemove(card, "dark-card");
      attr(options, "style", "width: 160px;");
      remove(tick);
    });
  }


  global.commentDelete = function(commentHex) {
    var json = {
      "session": sessionGet(),
      "commentHex": commentHex,
    }

    post(origin + "/api/comment/delete", json, function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return
      }

      var card = $(ID_CARD + commentHex);
      remove(card);
    });
  }


  function nameWidthFix() {
    var els = document.getElementsByClassName("commento-name");

    for (var i = 0; i < els.length; i++)
      attr(els[i], "style", "max-width: " + (els[i].getBoundingClientRect()["width"] + 20) + "px;")
  }


  global.vote = function(commentHex, oldVote, direction) {
    var upvote = $(ID_UPVOTE + commentHex);
    var downvote = $(ID_DOWNVOTE + commentHex);
    var score = $(ID_SCORE + commentHex);

    var json = {
      "session": sessionGet(),
      "commentHex": commentHex,
      "direction": direction,
    };

    if (direction > 0) {
      attr(upvote, "onclick", "vote('" + commentHex + "', 1, 0)");
      attr(downvote, "onclick", "vote('" + commentHex + "', 1, -1)");
    }
    else if (direction < 0) {
      attr(upvote, "onclick", "vote('" + commentHex + "', -1, 1)");
      attr(downvote, "onclick", "vote('" + commentHex + "', -1, 0)");
    }
    else {
      attr(upvote, "onclick", "vote('" + commentHex + "', 0, 1)");
      attr(downvote, "onclick", "vote('" + commentHex + "', 0, -1)");
    }

    classRemove(upvote, "upvoted");
    classRemove(downvote, "downvoted");
    if (direction > 0)
      classAdd(upvote, "upvoted");
    else if (direction < 0)
      classAdd(downvote, "downvoted");

    score.innerText = scorify(parseInt(score.innerText.replace(/[^\d-.]/g, "")) + direction - oldVote);

    post(origin + "/api/comment/vote", json, function(resp) {});
  }

  global.replyShow = function(id) {
    if (id in shownReply && shownReply[id])
      return;

    var body = $(ID_BODY + id);
    append(body, textareaCreate(id));
    shownReply[id] = true;

    var replyButton = $(ID_REPLY + id);

    classRemove(replyButton, "option-reply");
    classAdd(replyButton, "option-cancel");

    replyButton.title = "Cancel reply";

    attr(replyButton, "onclick", "replyCollapse('" + id + "')")
  };

  global.replyCollapse = function(id) {
    var replyButton = $(ID_REPLY + id);
    var el = $(ID_SUPER_CONTAINER + id);

    el.remove();
    shownReply[id] = false;
    shownSubmitButton[id] = false;

    classAdd(replyButton, "option-reply");
    classRemove(replyButton, "option-cancel");

    replyButton.title = "Reply to this comment";

    attr(replyButton, "onclick", "replyShow('" + id + "')")
  }

  global.commentCollapse = function(id) {
    var contents = $(ID_CONTENTS + id);
    var button = $(ID_COLLAPSE + id);

    classAdd(contents, "hidden");

    classRemove(button, "option-collapse");
    classAdd(button, "option-uncollapse");

    button.title = "Expand";

    attr(button, "onclick", "commentUncollapse('" + id + "')");
  }

  global.commentUncollapse = function(id) {
    var contents = $(ID_CONTENTS + id);
    var button = $(ID_COLLAPSE + id);

    classRemove(contents, "hidden");

    classRemove(button, "option-uncollapse");
    classAdd(button, "option-collapse");

    button.title = "Collapse";

    attr(button, "onclick", "commentCollapse('" + id + "')");
  }

  function commentsRender() {
    var parentMap = {};
    var parentHex;

    var commentsArea = $(ID_COMMENTS_AREA);

    comments.forEach(function(comment) {
      parentHex = comment.parentHex;
      if (!(parentHex in parentMap)) {
        parentMap[parentHex] = [];
      }
      parentMap[parentHex].push(comment);
    });

    var cards = commentsRecurse(parentMap, "root");
    if (cards) {
      append(commentsArea, cards);
    }
  }

  global.showSubmitButton = function(id) {
    if (id in shownSubmitButton && shownSubmitButton[id])
      return;

    shownSubmitButton[id] = true;

    var el = $(ID_SUPER_CONTAINER + id);

    var submit = create("button");

    submit.id = ID_SUBMIT_BUTTON + id;

    submit.innerText = "Add Comment";

    classAdd(submit, "button");
    classAdd(submit, "submit-button");
    classAdd(el, "button-margin");

    attr(submit, "onclick", "postComment('" + id + "')");

    append(el, submit);
  }

  global.commentoAuth = function(provider, id) {
    if (provider == "anonymous") {
      cookieSet("session", "anonymous");
      chosenAnonymous = true;
      refreshAll(function() {
        if (id != "root")
          global.replyShow(id);
        $(ID_TEXTAREA + id).click();
        $(ID_TEXTAREA + id).focus();
      });
      return;
    }

    var popup = window.open("", "_blank");

    get(origin + "/api/commenter/session/new", function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return;
      }

      cookieSet("session", resp.session);

      popup.location = origin + "/api/oauth/" + provider + "/redirect?session=" + resp.session;

      var interval = setInterval(function() {
        if (popup.closed) {
          refreshAll(function() {
            if (id != "root")
              global.replyShow(id);
            $(ID_TEXTAREA + id).click();
            $(ID_TEXTAREA + id).focus();
          });
          clearInterval(interval);
        }
      }, 250);
    });
  }

  function refreshAll(callback) {
    $("commento").innerHTML = "";
    shownSubmitButton = {"root": false};
    shownReply = {};
    main(callback);
  }

  function main(callback) {
    root = $("commento");

    errorElementCreate();

    selfGet(function() {
      commentsGet(function() {
        rootCreate(function() {
          commentsRender();
          nameWidthFix();
          footerLoad();
          attr(root, "style", "");
          call(callback);
        });
      });
    });
  }

  document.addEventListener("DOMContentLoaded", main);

}(window, document));
