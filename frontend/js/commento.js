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


  var ID_ROOT = "commento";
  var ID_MAIN_AREA = "commento-main-area";
  var ID_LOGIN_BOX_CONTAINER = "commento-login-box-container";
  var ID_LOGIN_BOX = "commento-login-box";
  var ID_LOGIN_BOX_HEADER = "commento-login-box-header";
  var ID_LOGIN_BOX_SUBTITLE = "commento-login-box-subtitle";
  var ID_LOGIN_BOX_EMAIL_INPUT = "commento-login-box-email-input";
  var ID_LOGIN_BOX_PASSWORD_INPUT = "commento-login-box-password-input";
  var ID_LOGIN_BOX_NAME_INPUT = "commento-login-box-name-input";
  var ID_LOGIN_BOX_WEBSITE_INPUT = "commento-login-box-website-input";
  var ID_LOGIN_BOX_EMAIL_BUTTON = "commento-login-box-email-button";
  var ID_LOGIN_BOX_LOGIN_LINK_CONTAINER = "commento-login-box-login-link-container";
  var ID_LOGIN_BOX_LOGIN_LINK = "commento-login-box-login-link";
  var ID_LOGIN_BOX_HR = "commento-login-box-hr";
  var ID_LOGIN_BOX_OAUTH_PRETEXT = "commento-login-box-oauth-pretext";
  var ID_LOGIN_BOX_OAUTH_BUTTONS_CONTAINER = "commento-login-box-oauth-buttons-container";
  var ID_ERROR = "commento-error";
  var ID_LOGGED_CONTAINER = "commento-logged-container";
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
  var ID_FOOTER = "commento-footer";


  var origin = global.commentoOrigin;
  var cdn = global.commentoCdn;
  var root = null;
  var cssOverride = undefined;
  var autoInit = undefined;
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
  var configuredOauths = [];
  var loginBoxType = "signup";


  function $(id) {
    return document.getElementById(id);
  }


  function tags(tag) {
    return document.getElementsByTagName(tag);
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


  function attrGet(node, a) {
    var attr = node.attributes[a];

    if (attr === undefined)
      return undefined;
    
    return attr.value;
  }


  function attrSet(node, a, value) {
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


  function commenterTokenGet() {
    var commenterToken = cookieGet("commentoCommenterToken");
    if (commenterToken === undefined)
      return "anonymous";

    return commenterToken;
  }


  global.logout = function() {
    cookieSet("commentoCommenterToken", "anonymous");
    refreshAll();
  }


  function selfGet(callback) {
    var commenterToken = commenterTokenGet();
    if (commenterToken == "anonymous") {
      isAuthenticated = false;
      call(callback);
      return;
    }

    var json = {
      "commenterToken": commenterTokenGet(),
    };

    post(origin + "/api/commenter/self", json, function(resp) {
      if (!resp.success) {
        cookieSet("commentoCommenterToken", "anonymous");
        call(callback);
        return;
      }

      var loggedContainer = create("div");
      var loggedInAs = create("div");
      var name = create("a");
      var avatar;
      var logout = create("div");
      var color = colorGet(resp.commenter.name);

      loggedContainer.id = ID_LOGGED_CONTAINER;

      classAdd(loggedContainer, "logged-container");
      classAdd(loggedInAs, "logged-in-as");
      classAdd(name, "name");
      classAdd(logout, "logout");

      name.innerText = resp.commenter.name;
      logout.innerText = "Logout";

      attrSet(loggedContainer, "style", "display: none");
      attrSet(logout, "onclick", "logout()");
      attrSet(name, "href", resp.commenter.link);
      if (resp.commenter.photo == "undefined") {
        avatar = create("div");
        avatar.style["background"] = color;
        avatar.style["boxShadow"] = "0px 0px 0px 2px " + color;
        avatar.innerHTML = resp.commenter.name[0].toUpperCase();
        classAdd(avatar, "avatar");
      } else {
        avatar = create("img");
        if (resp.commenter.provider == "google") {
          attrSet(avatar, "src", resp.commenter.photo + "?sz=50");
        } else {
          attrSet(avatar, "src", resp.commenter.photo);
        }
        classAdd(avatar, "avatar-img");
      }

      append(loggedInAs, avatar);
      append(loggedInAs, name);
      append(loggedContainer, loggedInAs);
      append(loggedContainer, logout);
      append(root, loggedContainer);

      isAuthenticated = true;

      call(callback);
    });
  }


  function cssLoad(file, onload) {
    var link = create("link");
    var head = document.getElementsByTagName('head')[0];

    link.type = "text/css";
    attrSet(link, "href", file);
    attrSet(link, "rel", "stylesheet");
    attrSet(link, "onload", onload);

    append(head, link);
  }


  function footerLoad() {
    var footer = create("div");
    var aContainer = create("div");
    var a = create("a");
    var img = create("img");
    var text = create("span");

    footer.id = ID_FOOTER;

    classAdd(footer, "footer");
    classAdd(aContainer, "logo-container");
    classAdd(a, "logo");
    classAdd(img, "logo-svg");
    classAdd(text, "logo-text");

    attrSet(footer, "style", "display: none");
    attrSet(a, "href", "https://commento.io");
    attrSet(a, "target", "_blank");
    attrSet(img, "src", cdn + "/images/logo.svg");

    text.innerText = "Powered by Commento";

    append(a, img);
    append(a, text);
    append(aContainer, a);
    append(footer, aContainer);
    append(root, footer);
  }


  function commentsGet(callback) {
    var json = {
      "commenterToken": commenterTokenGet(),
      "domain": location.host,
      "path": location.pathname,
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
      configuredOauths = resp.configuredOauths;

      if (!resp.requireIdentification)
        configuredOauths.push("anonymous");

      cssLoad(cdn + "/css/commento.css", "window.loadCssOverride()");

      call(callback);
    });
  }


  function errorShow(text) {
    var el = $(ID_ERROR);

    el.innerText = text;

    attrSet(el, "style", "display: block;");
  }


  function errorElementCreate() {
    var el = create("div");

    el.id = ID_ERROR;

    classAdd(el, "error-box");
    attrSet(el, "style", "display: none;");

    append(root, el);
  }


  function autoExpander(el) {
    return function() {
      el.style.height = "";
      el.style.height = Math.min(Math.max(el.scrollHeight, 75), 400) + "px";
    }
  };


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
      var createButton = create("div");

      classAdd(buttonsContainer, "create-container");

      classAdd(createButton, "button");
      classAdd(createButton, "create-button");

      attrSet(createButton, "onclick", "loginBoxShow()");
      attrSet(textarea, "disabled", true);

      createButton.innerText = "Create an Account";

      append(buttonsContainer, createButton);
      append(textareaContainer, buttonsContainer);
    }
    else {
      attrSet(textarea, "placeholder", "Join the discussion!");
      attrSet(textarea, "onclick", "showSubmitButton('" + id + "')");
    }

    textarea.oninput = autoExpander(textarea);

    append(textareaContainer, textarea);
    append(textareaSuperContainer, textareaContainer);

    return textareaSuperContainer;
  }


  function rootCreate(callback) {
    var mainArea = $(ID_MAIN_AREA);
    var commentsArea = create("div");

    commentsArea.id = ID_COMMENTS_AREA;

    classAdd(commentsArea, "comments");

    commentsArea.innerHTML = "";

    append(mainArea, textareaCreate("root"));
    append(mainArea, commentsArea);
    append(root, mainArea);

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
      "commenterToken": commenterTokenGet(),
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

        if (!resp.approved) {
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
          attrSet(avatar, "src", commenter.photo + "?sz=50");
        } else {
          attrSet(avatar, "src", commenter.photo);
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

      attrSet(edit, "onclick", "startEdit('" + comment.commentHex + "')");
      attrSet(collapse, "onclick", "commentCollapse('" + comment.commentHex + "')");
      attrSet(approve, "onclick", "commentApprove('" + comment.commentHex + "')");
      attrSet(remove, "onclick", "commentDelete('" + comment.commentHex + "')");

      if (isAuthenticated) {
        if (comment.direction > 0) {
          attrSet(upvote, "onclick", "vote('" + comment.commentHex + "', 1, 0)");
          attrSet(downvote, "onclick", "vote('" + comment.commentHex + "', 1, -1)");
        }
        else if (comment.direction < 0) {
          attrSet(upvote, "onclick", "vote('" + comment.commentHex + "', -1, 1)");
          attrSet(downvote, "onclick", "vote('" + comment.commentHex + "', -1, 0)");
        }
        else {
          attrSet(upvote, "onclick", "vote('" + comment.commentHex + "', 0, 1)");
          attrSet(downvote, "onclick", "vote('" + comment.commentHex + "', 0, -1)");
        }
      }
      else {
        attrSet(upvote, "onclick", "loginBoxShow()");
        attrSet(downvote, "onclick", "loginBoxShow()");
      }

      if (isAuthenticated || chosenAnonymous)
        attrSet(reply, "onclick", "replyShow('" + comment.commentHex + "')");
      else
        attrSet(reply, "onclick", "loginBoxShow()");

      if (commenter.link != "undefined")
        attrSet(name, "href", commenter.link);

      append(options, collapse);

      // append(options, edit); // uncomment when implemented
      append(options, downvote);
      append(options, upvote);
      append(options, reply);
      if (isModerator) {
        append(options, remove);
        if (comment.state == "unapproved")
          append(options, approve);
      }

      attrSet(options, "style", "width: " + ((options.childNodes.length+1)*32) + "px;");
      for (var i = 0; i < options.childNodes.length; i++)
        attrSet(options.childNodes[i], "style", "right: " + (i*32) + "px;");

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
      "commenterToken": commenterTokenGet(),
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
      remove(tick);
    });
  }


  global.commentDelete = function(commentHex) {
    var json = {
      "commenterToken": commenterTokenGet(),
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
      attrSet(els[i], "style", "max-width: " + (els[i].getBoundingClientRect()["width"] + 20) + "px;")
  }


  global.vote = function(commentHex, oldVote, direction) {
    var upvote = $(ID_UPVOTE + commentHex);
    var downvote = $(ID_DOWNVOTE + commentHex);
    var score = $(ID_SCORE + commentHex);

    var json = {
      "commenterToken": commenterTokenGet(),
      "commentHex": commentHex,
      "direction": direction,
    };

    if (direction > 0) {
      attrSet(upvote, "onclick", "vote('" + commentHex + "', 1, 0)");
      attrSet(downvote, "onclick", "vote('" + commentHex + "', 1, -1)");
    }
    else if (direction < 0) {
      attrSet(upvote, "onclick", "vote('" + commentHex + "', -1, 1)");
      attrSet(downvote, "onclick", "vote('" + commentHex + "', -1, 0)");
    }
    else {
      attrSet(upvote, "onclick", "vote('" + commentHex + "', 0, 1)");
      attrSet(downvote, "onclick", "vote('" + commentHex + "', 0, -1)");
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

    attrSet(replyButton, "onclick", "replyCollapse('" + id + "')")
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

    attrSet(replyButton, "onclick", "replyShow('" + id + "')")
  }


  global.commentCollapse = function(id) {
    var contents = $(ID_CONTENTS + id);
    var button = $(ID_COLLAPSE + id);

    classAdd(contents, "hidden");

    classRemove(button, "option-collapse");
    classAdd(button, "option-uncollapse");

    button.title = "Expand";

    attrSet(button, "onclick", "commentUncollapse('" + id + "')");
  }


  global.commentUncollapse = function(id) {
    var contents = $(ID_CONTENTS + id);
    var button = $(ID_COLLAPSE + id);

    classRemove(contents, "hidden");

    classRemove(button, "option-uncollapse");
    classAdd(button, "option-collapse");

    button.title = "Collapse";

    attrSet(button, "onclick", "commentCollapse('" + id + "')");
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

    attrSet(submit, "onclick", "postComment('" + id + "')");

    append(el, submit);
  }


  global.commentoAuth = function(provider) {
    if (provider == "anonymous") {
      cookieSet("commentoCommenterToken", "anonymous");
      chosenAnonymous = true;
      refreshAll();
      return;
    }

    var popup = window.open("", "_blank");

    get(origin + "/api/commenter/token/new", function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return;
      }

      cookieSet("commentoCommenterToken", resp.commenterToken);

      popup.location = origin + "/api/oauth/" + provider + "/redirect?commenterToken=" + resp.commenterToken;

      var interval = setInterval(function() {
        if (popup.closed) {
          refreshAll();
          clearInterval(interval);
        }
      }, 250);
    });
  }


  function refreshAll(callback) {
    $(ID_ROOT).innerHTML = "";
    shownSubmitButton = {"root": false};
    shownReply = {};
    main(callback);
  }


  function loginBoxCreate() {
    var loginBoxContainer = create("div");

    loginBoxContainer.id = ID_LOGIN_BOX_CONTAINER;

    append(root, loginBoxContainer);
  }


  global.signupRender = function() {
    var loginBoxContainer = $(ID_LOGIN_BOX_CONTAINER);
    var loginBox = create("div");
    var header = create("div");
    var subtitle = create("div");
    var emailContainer = create("div");
    var email = create("div");
    var emailInput = create("input");
    var emailButton = create("button");
    var loginLinkContainer = create("div");
    var loginLink = create("a");
    var hr = create("hr");
    var oauthPretext = create("div");
    var oauthButtonsContainer = create("div");
    var oauthButtons = create("div");
    var close = create("div");

    loginBox.id = ID_LOGIN_BOX;
    header.id = ID_LOGIN_BOX_HEADER;
    subtitle.id = ID_LOGIN_BOX_SUBTITLE;
    emailInput.id = ID_LOGIN_BOX_EMAIL_INPUT;
    emailButton.id = ID_LOGIN_BOX_EMAIL_BUTTON;
    loginLink.id = ID_LOGIN_BOX_LOGIN_LINK;
    loginLinkContainer.id = ID_LOGIN_BOX_LOGIN_LINK_CONTAINER;
    hr.id = ID_LOGIN_BOX_HR;
    oauthPretext.id = ID_LOGIN_BOX_OAUTH_PRETEXT;
    oauthButtonsContainer.id = ID_LOGIN_BOX_OAUTH_BUTTONS_CONTAINER;

    header.innerText = "Create an account to join the discussion!";

    classAdd(loginBoxContainer, "login-box-container");
    classAdd(loginBox, "login-box");
    classAdd(header, "login-box-header");
    classAdd(subtitle, "login-box-subtitle");
    classAdd(emailContainer, "email-container");
    classAdd(email, "email");
    classAdd(emailInput, "input");
    classAdd(emailButton, "email-button");
    classAdd(loginLinkContainer, "login-link-container");
    classAdd(loginLink, "login-link");
    classAdd(oauthPretext, "login-box-subtitle");
    classAdd(oauthButtonsContainer, "oauth-buttons-container");
    classAdd(oauthButtons, "oauth-buttons");
    classAdd(close, "login-box-close");

    emailButton.innerText = "Continue";
    loginLink.innerText = "Already have an account? Log in.";
    subtitle.innerText = "Sign up with your email to vote and comment.";
    oauthPretext.innerText = "Or proceed with social login.";

    attrSet(loginBoxContainer, "style", "display: none; opacity: 0;");
    attrSet(emailInput, "name", "email");
    attrSet(emailInput, "placeholder", "Email address");
    attrSet(emailInput, "type", "text");
    attrSet(emailButton, "onclick", "passwordAsk()");
    attrSet(loginLink, "onclick", "loginSwitch()");
    attrSet(close, "onclick", "loginBoxClose()");

    for (var i = 0; i < configuredOauths.length; i++) {
      var button = create("button");

      classAdd(button, "button");
      classAdd(button, configuredOauths[i] + "-button");

      button.innerText = configuredOauths[i];

      attrSet(button, "onclick", "commentoAuth('" + configuredOauths[i] + "')");

      append(oauthButtons, button);
    }

    append(loginBox, header);
    append(loginBox, subtitle);

    append(email, emailInput);
    append(email, emailButton);
    append(emailContainer, email);
    append(loginBox, emailContainer);

    append(loginLinkContainer, loginLink);
    append(loginBox, loginLinkContainer);

    if (configuredOauths.length > 0) {
      append(loginBox, hr);
      append(loginBox, oauthPretext);
      append(oauthButtonsContainer, oauthButtons);
      append(loginBox, oauthButtonsContainer);
    }

    append(loginBox, close);

    loginBoxType = "signup";
    loginBoxContainer.innerHTML = "";
    append(loginBoxContainer, loginBox);
  }


  global.loginSwitch = function() {
    var header = $(ID_LOGIN_BOX_HEADER);
    var subtitle = $(ID_LOGIN_BOX_SUBTITLE);
    var loginLink = $(ID_LOGIN_BOX_LOGIN_LINK);
    var hr = $(ID_LOGIN_BOX_HR);
    var oauthButtonsContainer = $(ID_LOGIN_BOX_OAUTH_BUTTONS_CONTAINER);
    var oauthPretext = $(ID_LOGIN_BOX_OAUTH_PRETEXT);

    header.innerText = "Login to continue";
    loginLink.innerText = "Don't have an account? Sign up.";
    subtitle.innerText = "Enter your email address to start with.";

    attrSet(loginLink, "onclick", "signupSwitch()");

    loginBoxType = "login";

    if (configuredOauths.length > 0) {
      remove(hr);
      remove(oauthPretext);
      remove(oauthButtonsContainer);
    }
  }


  global.signupSwitch = function() {
    loginBoxClose();
    loginBoxShow();
  }


  function loginUP(username, password) {
    var json = {
      "email": username,
      "password": password,
    };

    post(origin + "/api/commenter/login", json, function(resp) {
      if (!resp.success) {
        loginBoxClose();
        errorShow(resp.message);
        return
      }

      cookieSet("commentoCommenterToken", resp.commenterToken);
      refreshAll();
    });
  }


  global.login = function() {
    var email = $(ID_LOGIN_BOX_EMAIL_INPUT);
    var password = $(ID_LOGIN_BOX_PASSWORD_INPUT);

    loginUP(email.value, password.value);
  }


  global.signup = function() {
    var email = $(ID_LOGIN_BOX_EMAIL_INPUT);
    var name = $(ID_LOGIN_BOX_NAME_INPUT);
    var website = $(ID_LOGIN_BOX_WEBSITE_INPUT);
    var password = $(ID_LOGIN_BOX_PASSWORD_INPUT);

    var json = {
      "email": email.value,
      "name": name.value,
      "website": website.value,
      "password": password.value,
    };

    post(origin + "/api/commenter/new", json, function(resp) {
      if (!resp.success) {
        loginBoxClose();
        errorShow(resp.message);
        return
      }

      loginUP(email.value, password.value);
    });
  }


  global.passwordAsk = function() {
    var loginBox = $(ID_LOGIN_BOX);
    var subtitle = $(ID_LOGIN_BOX_SUBTITLE);
    var emailInput = $(ID_LOGIN_BOX_EMAIL_INPUT);
    var emailButton = $(ID_LOGIN_BOX_EMAIL_BUTTON);
    var loginLinkContainer = $(ID_LOGIN_BOX_LOGIN_LINK_CONTAINER);
    var hr = $(ID_LOGIN_BOX_HR);
    var oauthButtonsContainer = $(ID_LOGIN_BOX_OAUTH_BUTTONS_CONTAINER);
    var oauthPretext = $(ID_LOGIN_BOX_OAUTH_PRETEXT);
    
    remove(emailButton);
    remove(loginLinkContainer);
    if (loginBoxType == "signup") {
      if (configuredOauths.length > 0) {
        remove(hr);
        remove(oauthPretext);
        remove(oauthButtonsContainer);
      }
    }

    var order, id, type, placeholder;

    if (loginBoxType == "signup") {
      var order = ["name", "website", "password"];
      var id = [ID_LOGIN_BOX_NAME_INPUT, ID_LOGIN_BOX_WEBSITE_INPUT, ID_LOGIN_BOX_PASSWORD_INPUT];
      var type = ["text", "text", "password"];
      var placeholder = ["Real Name", "Website (Optional)", "Password"];
    }
    else {
      var order = ["password"];
      var id = [ID_LOGIN_BOX_PASSWORD_INPUT];
      var type = ["password"];
      var placeholder = ["Password"];
    }

    subtitle.innerText = "Finish the rest of your profile to complete."

    for (var i = 0; i < order.length; i++) {
      var fieldContainer = create("div");
      var field = create("div");
      var fieldInput = create("input");

      fieldInput.id = id[i];

      classAdd(fieldContainer, "email-container");
      classAdd(field, "email");
      classAdd(fieldInput, "input");

      attrSet(fieldInput, "name", order[i]);
      attrSet(fieldInput, "type", type[i]);
      attrSet(fieldInput, "placeholder", placeholder[i]);

      append(field, fieldInput);
      append(fieldContainer, field);

      if (order[i] == "password") {
        var fieldButton = create("button");
        classAdd(fieldButton, "email-button");
        fieldButton.innerText = loginBoxType;

        if (loginBoxType == "signup")
          attrSet(fieldButton, "onclick", "signup()");
        else
          attrSet(fieldButton, "onclick", "login()");

        append(field, fieldButton);
      }

      append(loginBox, fieldContainer);
    }
  }


  function mainAreaCreate() {
    var mainArea = create("div");

    mainArea.id = ID_MAIN_AREA;

    classAdd(mainArea, "main-area");

    attrSet(mainArea, "style", "display: none");

    append(root, mainArea);
  }


  global.loadCssOverride = function() {
    if (cssOverride === undefined)
      global.allShow();
    else
      cssLoad(cssOverride, "window.allShow()");
  }


  global.allShow = function() {
    var mainArea = $(ID_MAIN_AREA);
    var loggedContainer = $(ID_LOGGED_CONTAINER);
    var footer = $(ID_FOOTER);

    attrSet(mainArea, "style", "");
    if (loggedContainer)
      attrSet(loggedContainer, "style", "");
    attrSet(footer, "style", "");

    nameWidthFix();
  }


  global.loginBoxClose = function() {
    var mainArea = $(ID_MAIN_AREA);
    var loginBoxContainer = $(ID_LOGIN_BOX_CONTAINER);

    classRemove(mainArea, "blurred");

    attrSet(loginBoxContainer, "style", "display: none");
  }


  global.loginBoxShow = function() {
    var mainArea = $(ID_MAIN_AREA);
    var loginBoxContainer = $(ID_LOGIN_BOX_CONTAINER);

    global.signupRender();

    classAdd(mainArea, "blurred");
    
    attrSet(loginBoxContainer, "style", "");

    window.location.hash = ID_LOGIN_BOX_CONTAINER;

    $(ID_LOGIN_BOX_EMAIL_INPUT).focus();
  }


  function dataTagsLoad() {
    var scripts = tags("script")
    for (var i = 0; i < scripts.length; i++) {
      if (scripts[i].src.match(/\/js\/commento\.js$/)) {
        cssOverride = attrGet(scripts[i], "data-css-override");

        autoInit = attrGet(scripts[i], "data-auto-init");

        ID_ROOT = attrGet(scripts[i], "data-id-root");
        if (ID_ROOT === undefined)
          ID_ROOT = "commento";
      }
    }
  }


  global.main = function(callback) {
    root = $(ID_ROOT);
    classAdd(root, "root");

    loginBoxCreate();

    errorElementCreate();

    mainAreaCreate();

    selfGet(function() {
      commentsGet(function() {
        rootCreate(function() {
          commentsRender();
          footerLoad();
          attrSet(root, "style", "");
          call(callback);
        });
      });
    });
  }


  var initted = false;
  function init() {
    if (initted)
      return;
    initted = true;

    dataTagsLoad();

    if (autoInit == "true" || autoInit === undefined)
      main(undefined);
    else if (autoInit != "false")
      console.log("[commento] error: invalid value for data-auto-init; allowed values: true, false");
  }


  var readyLoad = function() {
    var readyState = document.readyState;

    if (readyState == "loading") {
      // The document is still loading. The div we need to fill might not have
      // been parsed yet, so let's wait and retry when the readyState changes.
      // If there is more than one state change, we aren't affected because we
      // have a double-call protection in init().
      document.addEventListener("readystatechange", readyLoad);
    }
    else if (readyState == "interactive") {
      // The document has been parsed and DOM objects are now accessible. While
      // JS, CSS, and images are still loading, we don't need to wait.
      init();
    }
    else if (readyState == "complete") {
      // The page has fully loaded (including JS, CSS, and images). From our
      // point of view, this is practically no different from interactive.
      init();
    }
  };

  readyLoad();

}(window, document));
