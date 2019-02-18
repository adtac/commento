(function(global, document) {
  "use strict";

  if (global.commento === undefined) {
    console.log("[commento] error: window.commento namespace not defined; maybe there's a mismatch in version between the backend and the frontend?");
    return;
  } else {
    global = global.commento;
  }


  // Do not use other files like utils.js and http.js in the gulpfile to build
  // commento.js for the following reasons:
  //   - We don't use jQuery in the actual JavaScript payload because we need
  //     to be lightweight.
  //   - They pollute the global/window namespace (with global.post, etc.).
  //     That's NOT fine when we expect them to source our JavaScript. For example,
  //     the user may have their own window.post defined. We don't want to
  //     override that.


  var ID_ROOT = "commento";
  var ID_MAIN_AREA = "commento-main-area";
  var ID_LOGIN = "commento-login";
  var ID_LOGIN_BOX_CONTAINER = "commento-login-box-container";
  var ID_LOGIN_BOX = "commento-login-box";
  var ID_LOGIN_BOX_EMAIL_SUBTITLE = "commento-login-box-email-subtitle";
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
  var ID_MOD_TOOLS = "commento-mod-tools";
  var ID_MOD_TOOLS_LOCK_BUTTON = "commento-mod-tools-lock-button";
  var ID_ERROR = "commento-error";
  var ID_LOGGED_CONTAINER = "commento-logged-container";
  var ID_COMMENTS_AREA = "commento-comments-area";
  var ID_SUPER_CONTAINER = "commento-textarea-super-container-";
  var ID_TEXTAREA_CONTAINER = "commento-textarea-container-";
  var ID_TEXTAREA = "commento-textarea-";
  var ID_ANONYMOUS_CHECKBOX = "commento-anonymous-checkbox-";
  var ID_CARD = "commento-comment-card-";
  var ID_BODY = "commento-comment-body-";
  var ID_TEXT = "commento-comment-text-";
  var ID_SUBTITLE = "commento-comment-subtitle-";
  var ID_TIMEAGO = "commento-comment-timeago-";
  var ID_SCORE = "commento-comment-score-";
  var ID_OPTIONS = "commento-comment-options-";
  var ID_EDIT = "commento-comment-edit-";
  var ID_REPLY = "commento-comment-reply-";
  var ID_COLLAPSE = "commento-comment-collapse-";
  var ID_UPVOTE = "commento-comment-upvote-";
  var ID_DOWNVOTE = "commento-comment-downvote-";
  var ID_APPROVE = "commento-comment-approve-";
  var ID_REMOVE = "commento-comment-remove-";
  var ID_STICKY = "commento-comment-sticky-";
  var ID_CHILDREN = "commento-comment-children-";
  var ID_CONTENTS = "commento-comment-contents-";
  var ID_NAME = "commento-comment-name-";
  var ID_SUBMIT_BUTTON = "commento-submit-button-";
  var ID_FOOTER = "commento-footer";


  var origin = "[[[.Origin]]]";
  var cdn = "[[[.CdnPrefix]]]";
  var root = null;
  var cssOverride;
  var autoInit;
  var isAuthenticated = false;
  var comments = [];
  var commenters = {};
  var requireIdentification = true;
  var isModerator = false;
  var isFrozen = false;
  var chosenAnonymous = false;
  var isLocked = false;
  var stickyCommentHex = "none";
  var shownReply = {};
  var configuredOauths = [];
  var popupBoxType = "login";
  var oauthButtonsShown = false;
  var selfHex = undefined;


  function $(id) {
    return document.getElementById(id);
  }


  function tags(tag) {
    return document.getElementsByTagName(tag);
  }


  function prepend(root, el) {
    root.prepend(el);
  }


  function append(root, el) {
    root.appendChild(el);
  }


  function insertAfter(el1, el2) {
    el1.parentNode.insertBefore(el2, el1.nextSibling);
  }


  function classAdd(el, cls) {
    el.classList.add("commento-" + cls);
  }


  function classRemove(el, cls) {
    if (el !== null) {
      el.classList.remove("commento-" + cls);
    }
  }


  function create(el) {
    return document.createElement(el);
  }


  function remove(el) {
    if (el !== null) {
      el.parentNode.removeChild(el);
    }
  }


  function attrGet(node, a) {
    var attr = node.attributes[a];

    if (attr === undefined) {
      return undefined;
    }
    
    return attr.value;
  }


  function onclick(node, f, arg) {
    node.addEventListener("click", function() {
      f(arg);
    }, false);
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

    xmlDoc.open("GET", url, true);
    xmlDoc.onload = function() {
      callback(JSON.parse(xmlDoc.response));
    };

    xmlDoc.send(null);
  }


  function call(callback) {
    if (typeof(callback) === "function") {
      callback();
    }
  }


  function cookieGet(name) {
    var c = "; " + document.cookie;
    var x = c.split("; " + name + "=");
    if (x.length === 2) {
      return x.pop().split(";").shift();
    }
  }


  function cookieSet(name, value) {
    var expires = "";
    var date = new Date();
    date.setTime(date.getTime() + (365 * 24 * 60 * 60 * 1000));
    expires = "; expires=" + date.toUTCString();

    document.cookie = name + "=" + value + expires + "; path=/";
  }


  function commenterTokenGet() {
    var commenterToken = cookieGet("commentoCommenterToken");
    if (commenterToken === undefined) {
      return "anonymous";
    }

    return commenterToken;
  }


  global.logout = function() {
    cookieSet("commentoCommenterToken", "anonymous");
    refreshAll();
  }

  function selfLoad(commenter) {
    commenters[commenter.commenterHex] = commenter;
    selfHex = commenter.commenterHex;

    var loggedContainer = create("div");
    var loggedInAs = create("div");
    var name = create("a");
    var avatar;
    var logout = create("div");
    var color = colorGet(commenter.commenterHex + "-" + commenter.name);

    loggedContainer.id = ID_LOGGED_CONTAINER;

    classAdd(loggedContainer, "logged-container");
    classAdd(loggedInAs, "logged-in-as");
    classAdd(name, "name");
    classAdd(logout, "logout");

    name.innerText = commenter.name;
    logout.innerText = "Logout";

    onclick(logout, global.logout);

    attrSet(loggedContainer, "style", "display: none");
    attrSet(name, "href", commenter.link);
    if (commenter.photo === "undefined") {
      avatar = create("div");
      avatar.style["background"] = color;
      avatar.innerHTML = commenter.name[0].toUpperCase();
      classAdd(avatar, "avatar");
    } else {
      avatar = create("img");
      if (commenter.provider === "google") {
        attrSet(avatar, "src", commenter.photo + "?sz=50");
      } else if (commenter.provider === "github") {
        attrSet(avatar, "src", commenter.photo + "&s=50");
      } else {
        attrSet(avatar, "src", commenter.photo);
      }
      classAdd(avatar, "avatar-img");
    }

    append(loggedInAs, avatar);
    append(loggedInAs, name);
    append(loggedContainer, loggedInAs);
    append(loggedContainer, logout);
    prepend(root, loggedContainer);

    isAuthenticated = true;
  }


  function selfGet(callback) {
    var commenterToken = commenterTokenGet();
    if (commenterToken === "anonymous") {
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

      selfLoad(resp.commenter);

      call(callback);
    });
  }


  function cssLoad(file, onload) {
    var link = create("link");
    var head = document.getElementsByTagName("head")[0];

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
      "domain": parent.location.host,
      "path": parent.location.pathname,
    };

    post(origin + "/api/comment/list", json, function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return;
      }

      requireIdentification = resp.requireIdentification;
      isModerator = resp.isModerator;
      isFrozen = resp.isFrozen;

      isLocked = resp.attributes.isLocked;
      stickyCommentHex = resp.attributes.stickyCommentHex;

      comments = resp.comments;
      commenters = Object.assign({}, commenters, resp.commenters)
      configuredOauths = resp.configuredOauths;

      cssLoad(cdn + "/css/commento.css", "window.commento.loadCssOverride()");

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
    var anonymousCheckboxContainer = create("div");
    var anonymousCheckbox = create("input");
    var anonymousCheckboxLabel = create("label");
    var submitButton = create("button");

    textareaSuperContainer.id = ID_SUPER_CONTAINER + id;
    textareaContainer.id = ID_TEXTAREA_CONTAINER + id;
    textarea.id = ID_TEXTAREA + id;
    anonymousCheckbox.id = ID_ANONYMOUS_CHECKBOX + id;
    submitButton.id = ID_SUBMIT_BUTTON + id;

    classAdd(textareaContainer, "textarea-container");
    classAdd(anonymousCheckboxContainer, "round-check");
    classAdd(anonymousCheckboxContainer, "anonymous-checkbox-container");
    classAdd(submitButton, "button");
    classAdd(submitButton, "submit-button");
    classAdd(textareaSuperContainer, "button-margin");

    attrSet(textarea, "placeholder", "Add a comment");
    attrSet(anonymousCheckbox, "type", "checkbox");
    attrSet(anonymousCheckboxLabel, "for", ID_ANONYMOUS_CHECKBOX + id);

    anonymousCheckboxLabel.innerText = "Comment anonymously";
    submitButton.innerText = "Add Comment";

    textarea.oninput = autoExpander(textarea);
    onclick(submitButton, submitAccountDecide, id);

    append(textareaContainer, textarea);
    append(textareaSuperContainer, textareaContainer);
    append(anonymousCheckboxContainer, anonymousCheckbox);
    append(anonymousCheckboxContainer, anonymousCheckboxLabel);
    append(textareaSuperContainer, submitButton);
    if (!requireIdentification) {
      append(textareaSuperContainer, anonymousCheckboxContainer);
    }

    return textareaSuperContainer;
  }


  function rootCreate(callback) {
    var login = create("div");
    var loginText = create("div");
    var mainArea = $(ID_MAIN_AREA);
    var commentsArea = create("div");

    login.id = ID_LOGIN;
    commentsArea.id = ID_COMMENTS_AREA;

    classAdd(login, "login");
    classAdd(loginText, "login-text");
    classAdd(commentsArea, "comments");

    loginText.innerText = "Login";
    commentsArea.innerHTML = "";

    onclick(loginText, global.loginBoxShow, null);

    append(login, loginText);

    if (isLocked || isFrozen) {
      if (isAuthenticated || chosenAnonymous) {
        append(mainArea, messageCreate("This thread is locked. You cannot add new comments."));
        remove($(ID_LOGIN));
      } else {
        append(mainArea, login);
        append(mainArea, textareaCreate("root"));
      }
    } else {
      if (!isAuthenticated) {
        append(mainArea, login);
      } else {
        remove($(ID_LOGIN));
      }
      append(mainArea, textareaCreate("root"));
    }

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


  global.commentNew = function(id, commenterToken, callback) {
    var textareaSuperContainer = $(ID_SUPER_CONTAINER + id);
    var textarea = $(ID_TEXTAREA + id);
    var replyButton = $(ID_REPLY + id);

    var markdown = textarea.value;

    if (markdown === "") {
      classAdd(textarea, "red-border");
      return;
    } else {
      classRemove(textarea, "red-border");
    }

    var json = {
      "commenterToken": commenterTokenGet(),
      "domain": parent.location.host,
      "path": parent.location.pathname,
      "parentHex": id,
      "markdown": markdown,
    };

    post(origin + "/api/comment/new", json, function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return;
      }

      var message = "";
      if (resp.state === "unapproved") {
        message = "Your comment is under moderation.";
      } else if (resp.state === "flagged") {
        message = "Your comment was flagged as spam and is under moderation.";
      }

      if (message !== "") {
        prepend($(ID_SUPER_CONTAINER + id), messageCreate(message));
      }
      
      var commenterHex = selfHex;
      if (commenterHex === undefined || commenterToken === "anonymous") {
        commenterHex = "anonymous";
      }

      var newCard = commentsRecurse({
        "root": [{
          "commentHex": resp.commentHex,
          "commenterHex": commenterHex,
          "markdown": markdown,
          "html": resp.html,
          "parentHex": "root",
          "score": 0,
          "state": "approved",
          "direction": 0,
          "creationDate": (new Date()).toISOString(),
        }],
      }, "root")

      if (id !== "root") {
        textareaSuperContainer.replaceWith(newCard);

        shownReply[id] = false;

        classAdd(replyButton, "option-reply");
        classRemove(replyButton, "option-cancel");

        replyButton.title = "Reply to this comment";

        onclick(replyButton, global.replyShow, id)
      } else {
        textarea.value = "";
        insertAfter(textareaSuperContainer, newCard);
      }

      call(callback);
    });
  }


  function colorGet(name) {
    var colors = [
      "#396ab1",
      "#da7c30",
      "#3e9651",
      "#cc2529",
      "#535154",
      "#6b4c9a",
      "#922428",
    ];

    var total = 0;
    for (var i = 0; i < name.length; i++) {
      total += name.charCodeAt(i);
    }
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
      return "just now";
    } else if (elapsed < msPerMinute) {
      return Math.round(elapsed / 1000) + " seconds ago";
    } else if (elapsed < msPerHour) {
      return Math.round(elapsed / msPerMinute) + " minutes ago";
    } else if (elapsed < msPerDay ) {
      return Math.round(elapsed / msPerHour ) + " hours ago";
    } else if (elapsed < msPerMonth) {
      return Math.round(elapsed / msPerDay) + " days ago";
    } else if (elapsed < msPerYear) {
      return Math.round(elapsed / msPerMonth) + " months ago";
    } else {
      return Math.round(elapsed / msPerYear ) + " years ago";
    }
  }


  function scorify(score) {
    if (score !== 1) {
      return score + " points";
    } else {
      return score + " point";
    }
  }


  function commentsRecurse(parentMap, parentHex) {
    var cur = parentMap[parentHex];
    if (!cur || !cur.length) {
      return null;
    }

    cur.sort(function(a, b) {
      if (a.commentHex === stickyCommentHex) {
        return -Infinity;
      }
      if (b.commentHex === stickyCommentHex) {
        return Infinity;
      }
      return b.score - a.score;
    });

    var cards = create("div");
    cur.forEach(function(comment) {
      var mobileView = root.getBoundingClientRect()["width"] < 450;
      var commenter = commenters[comment.commenterHex];
      var avatar;
      var card = create("div");
      var header = create("div");
      var subtitle = create("div");
      var timeago = create("div");
      var score = create("div");
      var body = create("div");
      var text = create("div");
      var options = create("div");
      var edit = create("button");
      var reply = create("button");
      var collapse = create("button");
      var upvote = create("button");
      var downvote = create("button");
      var approve = create("button");
      var remove = create("button");
      var sticky = create("button");
      var children = commentsRecurse(parentMap, comment.commentHex);
      var contents = create("div");
      var color = colorGet(comment.commenterHex + "-" + commenter.name);
      var name;
      if (commenter.link !== "undefined" && commenter.link !== "https://undefined" && commenter.link !== "") {
        name = create("a");
      } else {
        name = create("div");
      }

      card.id = ID_CARD + comment.commentHex;
      body.id = ID_BODY + comment.commentHex;
      text.id = ID_TEXT + comment.commentHex;
      subtitle.id = ID_SUBTITLE + comment.commentHex;
      timeago.id = ID_TIMEAGO + comment.commentHex;
      score.id = ID_SCORE + comment.commentHex;
      options.id = ID_OPTIONS + comment.commentHex;
      edit.id = ID_EDIT + comment.commentHex;
      reply.id = ID_REPLY + comment.commentHex;
      collapse.id = ID_COLLAPSE + comment.commentHex;
      upvote.id = ID_UPVOTE + comment.commentHex;
      downvote.id = ID_DOWNVOTE + comment.commentHex;
      approve.id = ID_APPROVE + comment.commentHex;
      remove.id = ID_REMOVE + comment.commentHex;
      sticky.id = ID_STICKY + comment.commentHex;
      if (children) {
        children.id = ID_CHILDREN + comment.commentHex;
      }
      contents.id = ID_CONTENTS + comment.commentHex;
      name.id = ID_NAME + comment.commentHex;

      collapse.title = "Collapse children";
      upvote.title = "Upvote";
      downvote.title = "Downvote";
      edit.title = "Edit";
      reply.title = "Reply";
      approve.title = "Approve";
      remove.title = "Remove";
      if (stickyCommentHex === comment.commentHex) {
        if (isModerator) {
          sticky.title = "Unsticky";
        } else {
          sticky.title = "This comment has been stickied";
        }
      } else {
        sticky.title = "Sticky";
      }

      card.style["borderLeft"] = "2px solid " + color;
      name.innerText = commenter.name;
      text.innerHTML = comment.html;
      timeago.innerHTML = timeDifference((new Date()).getTime(), Date.parse(comment.creationDate));
      score.innerText = scorify(comment.score);

      if (commenter.photo === "undefined") {
        avatar = create("div");
        avatar.style["background"] = color;
        avatar.innerHTML = commenter.name[0].toUpperCase();
        classAdd(avatar, "avatar");
      } else {
        avatar = create("img");
        if (commenter.provider === "google") {
          attrSet(avatar, "src", commenter.photo + "?sz=50");
        } else if (commenter.provider === "github") {
          attrSet(avatar, "src", commenter.photo + "&s=50");
        } else {
          attrSet(avatar, "src", commenter.photo);
        }
        classAdd(avatar, "avatar-img");
      }

      classAdd(card, "card");
      if (isModerator && comment.state !== "approved") {
        classAdd(card, "dark-card");
      }
      if (comment.state === "flagged") {
        classAdd(name, "flagged");
      }
      classAdd(header, "header");
      classAdd(name, "name");
      classAdd(subtitle, "subtitle");
      classAdd(timeago, "timeago");
      classAdd(score, "score");
      classAdd(body, "body");
      classAdd(options, "options");
      if (mobileView) {
        classAdd(options, "options-mobile");
      }
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
      classAdd(sticky, "option-button");
      if (stickyCommentHex === comment.commentHex) {
        classAdd(sticky, "option-unsticky");
      } else {
        classAdd(sticky, "option-sticky");
      }

      if (isAuthenticated) {
        if (comment.direction > 0) {
          classAdd(upvote, "upvoted");
        } else if (comment.direction < 0) {
          classAdd(downvote, "downvoted");
        }
      }

      onclick(collapse, global.commentCollapse, comment.commentHex);
      onclick(approve, global.commentApprove, comment.commentHex);
      onclick(remove, global.commentDelete, comment.commentHex);
      onclick(sticky, global.commentSticky, comment.commentHex);

      if (isAuthenticated) {
        upDownOnclickSet(upvote, downvote, comment.commentHex, comment.direction);
      } else {
        onclick(upvote, global.loginBoxShow, null);
        onclick(downvote, global.loginBoxShow, null);
      }

      onclick(reply, global.replyShow, comment.commentHex);

      if (commenter.link !== "undefined" && commenter.link !== "https://undefined" && commenter.link !== "") {
        attrSet(name, "href", commenter.link);
      }

      append(options, collapse);

      // append(options, edit); // uncomment when implemented
      append(options, downvote);
      append(options, upvote);

      append(options, reply);

      if (isModerator) {
        if (parentHex === "root") {
          append(options, sticky);
        }
        append(options, remove);
        if (comment.state !== "approved") {
          append(options, approve);
        }
      } else {
        if (stickyCommentHex === comment.commentHex) {
          append(options, sticky);
        }
      }

      attrSet(options, "style", "width: " + ((options.childNodes.length+1)*32) + "px;");
      for (var i = 0; i < options.childNodes.length; i++) {
        attrSet(options.childNodes[i], "style", "right: " + (i*32) + "px;");
      }

      append(subtitle, score);
      append(subtitle, timeago);

      if (!mobileView) {
        append(header, options);
      }
      append(header, avatar);
      append(header, name);
      append(header, subtitle);
      append(body, text);
      append(contents, body);
      if (mobileView) {
        append(contents, options);
        var optionsClearfix = create("div");
        classAdd(optionsClearfix, "options-clearfix");
        append(contents, optionsClearfix);
      }

      if (children) {
        classAdd(children, "body");
        append(contents, children);
      }

      append(card, header);
      append(card, contents);
      append(cards, card);
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
      var name = $(ID_NAME + commentHex);
      var tick = $(ID_APPROVE + commentHex);

      classRemove(card, "dark-card");
      classRemove(name, "flagged");
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

    for (var i = 0; i < els.length; i++) {
      attrSet(els[i], "style", "max-width: " + (els[i].getBoundingClientRect()["width"] + 20) + "px;")
    }
  }


  function upDownOnclickSet(upvote, downvote, commentHex, direction) {
    if (direction > 0) {
      onclick(upvote, global.vote, [commentHex, [1, 0]]);
      onclick(downvote, global.vote, [commentHex, [1, -1]]);
    } else if (direction < 0) {
      onclick(upvote, global.vote, [commentHex, [-1, 1]]);
      onclick(downvote, global.vote, [commentHex, [-1, 0]]);
    } else {
      onclick(upvote, global.vote, [commentHex, [0, 1]]);
      onclick(downvote, global.vote, [commentHex, [0, -1]]);
    }
  }


  global.vote = function(data) {
    var commentHex = data[0];
    var oldDirection = data[1][0];
    var newDirection = data[1][1];

    var upvote = $(ID_UPVOTE + commentHex);
    var downvote = $(ID_DOWNVOTE + commentHex);
    var score = $(ID_SCORE + commentHex);

    var json = {
      "commenterToken": commenterTokenGet(),
      "commentHex": commentHex,
      "direction": newDirection,
    };

    upDownOnclickSet(upvote, downvote, commentHex, newDirection);

    classRemove(upvote, "upvoted");
    classRemove(downvote, "downvoted");
    if (newDirection > 0) {
      classAdd(upvote, "upvoted");
    } else if (newDirection < 0) {
      classAdd(downvote, "downvoted");
    }

    score.innerText = scorify(parseInt(score.innerText.replace(/[^\d-.]/g, "")) + newDirection - oldDirection);

    post(origin + "/api/comment/vote", json, function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return;
      }
    });
  }


  global.replyShow = function(id) {
    if (id in shownReply && shownReply[id]) {
      return;
    }

    var text = $(ID_TEXT + id);
    insertAfter(text, textareaCreate(id));
    shownReply[id] = true;

    var replyButton = $(ID_REPLY + id);

    classRemove(replyButton, "option-reply");
    classAdd(replyButton, "option-cancel");

    replyButton.title = "Cancel reply";

    onclick(replyButton, global.replyCollapse, id);
  };


  global.replyCollapse = function(id) {
    var replyButton = $(ID_REPLY + id);
    var el = $(ID_SUPER_CONTAINER + id);

    el.remove();
    shownReply[id] = false;

    classAdd(replyButton, "option-reply");
    classRemove(replyButton, "option-cancel");

    replyButton.title = "Reply to this comment";

    onclick(replyButton, global.replyShow, id)
  }


  global.commentCollapse = function(id) {
    var children = $(ID_CHILDREN + id);
    var button = $(ID_COLLAPSE + id);

    if (children) {
      classAdd(children, "hidden");
    }

    classRemove(button, "option-collapse");
    classAdd(button, "option-uncollapse");

    button.title = "Expand children";

    onclick(button, global.commentUncollapse, id);
  }


  global.commentUncollapse = function(id) {
    var children = $(ID_CHILDREN + id);
    var button = $(ID_COLLAPSE + id);

    if (children) {
      classRemove(children, "hidden");
    }

    classRemove(button, "option-uncollapse");
    classAdd(button, "option-collapse");

    button.title = "Collapse children";

    onclick(button, global.commentCollapse, id);
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


  function submitAuthenticated(id) {
    if (isAuthenticated) {
      global.commentNew(id, commenterTokenGet());
      return;
    }

    global.loginBoxShow(id);
    return;
  }


  function submitAnonymous(id) {
    chosenAnonymous = true;
    global.commentNew(id, "anonymous");
  }


  function submitAccountDecide(id) {
    if (requireIdentification) {
      submitAuthenticated(id);
      return;
    }

    var anonymousCheckbox = $(ID_ANONYMOUS_CHECKBOX + id);
    var textarea = $(ID_TEXTAREA + id);
    var markdown = textarea.value;

    if (markdown === "") {
      classAdd(textarea, "red-border");
      return;
    } else {
      classRemove(textarea, "red-border");
    }

    if (!anonymousCheckbox.checked) {
      submitAuthenticated(id);
      return;
    } else {
      submitAnonymous(id);
      return;
    }
  }


  global.commentoAuth = function(data) {
    var provider = data.provider;
    var id = data.id;
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
          clearInterval(interval);
          selfGet(function() {
            var loggedContainer = $(ID_LOGGED_CONTAINER);
            if (loggedContainer) {
              attrSet(loggedContainer, "style", "");
            }

            global.commentNew(id, resp.commenterToken, function() {
              global.loginBoxClose();
            });
          });
        }
      }, 250);
    });
  }


  function refreshAll(callback) {
    $(ID_ROOT).innerHTML = "";
    shownReply = {};
    global.main(callback);
  }


  function loginBoxCreate() {
    var loginBoxContainer = create("div");

    loginBoxContainer.id = ID_LOGIN_BOX_CONTAINER;

    append(root, loginBoxContainer);
  }


  global.popupRender = function(id) {
    var loginBoxContainer = $(ID_LOGIN_BOX_CONTAINER);
    var loginBox = create("div");
    var oauthSubtitle = create("div");
    var oauthButtonsContainer = create("div");
    var oauthButtons = create("div");
    var hr = create("hr");
    var emailSubtitle = create("div");
    var emailContainer = create("div");
    var email = create("div");
    var emailInput = create("input");
    var emailButton = create("button");
    var loginLinkContainer = create("div");
    var loginLink = create("a");
    var close = create("div");

    loginBox.id = ID_LOGIN_BOX;
    emailSubtitle.id = ID_LOGIN_BOX_EMAIL_SUBTITLE;
    emailInput.id = ID_LOGIN_BOX_EMAIL_INPUT;
    emailButton.id = ID_LOGIN_BOX_EMAIL_BUTTON;
    loginLink.id = ID_LOGIN_BOX_LOGIN_LINK;
    loginLinkContainer.id = ID_LOGIN_BOX_LOGIN_LINK_CONTAINER;
    hr.id = ID_LOGIN_BOX_HR;
    oauthSubtitle.id = ID_LOGIN_BOX_OAUTH_PRETEXT;
    oauthButtonsContainer.id = ID_LOGIN_BOX_OAUTH_BUTTONS_CONTAINER;

    classAdd(loginBoxContainer, "login-box-container");
    classAdd(loginBox, "login-box");
    classAdd(emailSubtitle, "login-box-subtitle");
    classAdd(emailContainer, "email-container");
    classAdd(email, "email");
    classAdd(emailInput, "input");
    classAdd(emailButton, "email-button");
    classAdd(loginLinkContainer, "login-link-container");
    classAdd(loginLink, "login-link");
    classAdd(oauthSubtitle, "login-box-subtitle");
    classAdd(oauthButtonsContainer, "oauth-buttons-container");
    classAdd(oauthButtons, "oauth-buttons");
    classAdd(close, "login-box-close");
    classAdd(root, "root-min-height");

    loginLink.innerText = "Don't have an account? Sign up.";
    emailSubtitle.innerText = "Login with your email address";
    emailButton.innerText = "Continue";
    oauthSubtitle.innerText = "Proceed with social login";

    onclick(emailButton, global.passwordAsk, id);
    onclick(loginLink, global.popupSwitch);
    onclick(close, global.loginBoxClose);

    attrSet(loginBoxContainer, "style", "display: none; opacity: 0;");
    attrSet(emailInput, "name", "email");
    attrSet(emailInput, "placeholder", "Email address");
    attrSet(emailInput, "type", "text");

    for (var i = 0; i < configuredOauths.length; i++) {
      var button = create("button");

      classAdd(button, "button");
      classAdd(button, configuredOauths[i] + "-button");

      button.innerText = configuredOauths[i];

      onclick(button, global.commentoAuth, {"provider": configuredOauths[i], "id": id});

      append(oauthButtons, button);
    }

    if (configuredOauths.length > 0) {
      append(loginBox, oauthSubtitle);
      append(oauthButtonsContainer, oauthButtons);
      append(loginBox, oauthButtonsContainer);
      append(loginBox, hr);
      oauthButtonsShown = true;
    } else {
      oauthButtonsShown = false;
    }

    append(loginBox, emailSubtitle);
    append(email, emailInput);
    append(email, emailButton);
    append(emailContainer, email);
    append(loginBox, emailContainer);

    append(loginLinkContainer, loginLink);
    append(loginBox, loginLinkContainer);

    append(loginBox, close);

    popupBoxType = "login";
    loginBoxContainer.innerHTML = "";
    append(loginBoxContainer, loginBox);
  }


  global.popupSwitch = function() {
    var emailSubtitle = $(ID_LOGIN_BOX_EMAIL_SUBTITLE);
    var loginLink = $(ID_LOGIN_BOX_LOGIN_LINK);

    if (popupBoxType === "login") {
      loginLink.innerText = "Already have an account? Log in.";
      emailSubtitle.innerText = "Create an account";
      popupBoxType = "signup";
    } else {
      loginLink.innerText = "Don't have an account? Sign up.";
      emailSubtitle.innerText = "Login with your email address";
      popupBoxType = "login";
    }
  }


  function loginUP(username, password, id) {
    var json = {
      "email": username,
      "password": password,
    };

    post(origin + "/api/commenter/login", json, function(resp) {
      if (!resp.success) {
        global.loginBoxClose();
        errorShow(resp.message);
        return
      }

      cookieSet("commentoCommenterToken", resp.commenterToken);

      selfLoad(resp.commenter);
      var loggedContainer = $(ID_LOGGED_CONTAINER);
      if (loggedContainer) {
        attrSet(loggedContainer, "style", "");
      }

      global.commentNew(id, resp.commenterToken, function() {
        global.loginBoxClose();
      });
    });
  }


  global.login = function(id) {
    var email = $(ID_LOGIN_BOX_EMAIL_INPUT);
    var password = $(ID_LOGIN_BOX_PASSWORD_INPUT);

    loginUP(email.value, password.value, id);
  }


  global.signup = function(id) {
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
        global.loginBoxClose();
        errorShow(resp.message);
        return
      }

      loginUP(email.value, password.value, id);
    });
  }


  global.passwordAsk = function(id) {
    var loginBox = $(ID_LOGIN_BOX);
    var subtitle = $(ID_LOGIN_BOX_EMAIL_SUBTITLE);
    var emailButton = $(ID_LOGIN_BOX_EMAIL_BUTTON);
    var loginLinkContainer = $(ID_LOGIN_BOX_LOGIN_LINK_CONTAINER);
    var hr = $(ID_LOGIN_BOX_HR);
    var oauthButtonsContainer = $(ID_LOGIN_BOX_OAUTH_BUTTONS_CONTAINER);
    var oauthPretext = $(ID_LOGIN_BOX_OAUTH_PRETEXT);
    
    remove(emailButton);
    remove(loginLinkContainer);
    if (oauthButtonsShown) {
      if (configuredOauths.length > 0) {
        remove(hr);
        remove(oauthPretext);
        remove(oauthButtonsContainer);
      }
    }

    var order, fid, type, placeholder;

    if (popupBoxType === "signup") {
      order = ["name", "website", "password"];
      fid = [ID_LOGIN_BOX_NAME_INPUT, ID_LOGIN_BOX_WEBSITE_INPUT, ID_LOGIN_BOX_PASSWORD_INPUT];
      type = ["text", "text", "password"];
      placeholder = ["Real Name", "Website (Optional)", "Password"];
    } else {
      order = ["password"];
      fid = [ID_LOGIN_BOX_PASSWORD_INPUT];
      type = ["password"];
      placeholder = ["Password"];
    }

    if (popupBoxType === "signup") {
      subtitle.innerText = "Finish the rest of your profile to complete."
    } else {
      subtitle.innerText = "Enter your password to log in."
    }

    for (var i = 0; i < order.length; i++) {
      var fieldContainer = create("div");
      var field = create("div");
      var fieldInput = create("input");

      fieldInput.id = fid[i];

      classAdd(fieldContainer, "email-container");
      classAdd(field, "email");
      classAdd(fieldInput, "input");

      attrSet(fieldInput, "name", order[i]);
      attrSet(fieldInput, "type", type[i]);
      attrSet(fieldInput, "placeholder", placeholder[i]);

      append(field, fieldInput);
      append(fieldContainer, field);

      if (order[i] === "password") {
        var fieldButton = create("button");
        classAdd(fieldButton, "email-button");
        fieldButton.innerText = popupBoxType;

        if (popupBoxType === "signup") {
          onclick(fieldButton, global.signup, id);
        } else {
          onclick(fieldButton, global.login, id);
        }

        append(field, fieldButton);
      }

      append(loginBox, fieldContainer);
    }

    if (popupBoxType === "signup") {
      $(ID_LOGIN_BOX_NAME_INPUT).focus();
    } else {
      $(ID_LOGIN_BOX_PASSWORD_INPUT).focus();
    }
  }


  function pageUpdate(callback) {
    var attributes = {
      "isLocked": isLocked,
      "stickyCommentHex": stickyCommentHex,
    };

    var json = {
      "commenterToken": commenterTokenGet(),
      "domain": parent.location.host,
      "path": parent.location.pathname,
      "attributes": attributes,
    };

    post(origin + "/api/page/update", json, function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return
      }

      call(callback);
    });
  }


  global.threadLockToggle = function() {
    var lock = $(ID_MOD_TOOLS_LOCK_BUTTON);

    isLocked = !isLocked;

    lock.disabled = true;
    pageUpdate(function() {
      lock.disabled = false;
      refreshAll();
    });
  }


  global.commentSticky = function(commentHex) {
    if (stickyCommentHex !== "none") {
      var sticky = $(ID_STICKY + stickyCommentHex);
      classRemove(sticky, "option-unsticky");
      classAdd(sticky, "option-sticky");
    }

    if (stickyCommentHex === commentHex) {
      stickyCommentHex = "none";
    } else {
      stickyCommentHex = commentHex;
    }

    pageUpdate(function() {
      var sticky = $(ID_STICKY + commentHex);
      if (stickyCommentHex === commentHex) {
        classRemove(sticky, "option-sticky");
        classAdd(sticky, "option-unsticky");
      } else {
        classRemove(sticky, "option-unsticky");
        classAdd(sticky, "option-sticky");
      }
    });
  }


  function mainAreaCreate() {
    var mainArea = create("div");

    mainArea.id = ID_MAIN_AREA;

    classAdd(mainArea, "main-area");

    attrSet(mainArea, "style", "display: none");

    append(root, mainArea);
  }


  function modToolsCreate() {
    var modTools = create("div");
    var lock = create("button");

    modTools.id = ID_MOD_TOOLS;
    lock.id = ID_MOD_TOOLS_LOCK_BUTTON;

    classAdd(modTools, "mod-tools");

    if (isLocked) {
      lock.innerHTML = "Unlock Thread";
    } else {
      lock.innerHTML = "Lock Thread";
    }

    onclick(lock, global.threadLockToggle);

    attrSet(modTools, "style", "display: none");

    append(modTools, lock);
    append(root, modTools);
  }


  global.loadCssOverride = function() {
    if (cssOverride === undefined) {
      global.allShow();
    } else {
      cssLoad(cssOverride, "window.allShow()");
    }
  }


  global.allShow = function() {
    var mainArea = $(ID_MAIN_AREA);
    var modTools = $(ID_MOD_TOOLS);
    var loggedContainer = $(ID_LOGGED_CONTAINER);
    var footer = $(ID_FOOTER);

    attrSet(mainArea, "style", "");

    if (isModerator) {
      attrSet(modTools, "style", "");
    }

    if (loggedContainer) {
      attrSet(loggedContainer, "style", "");
    }

    attrSet(footer, "style", "");

    nameWidthFix();
    loadHash();
  }


  global.loginBoxClose = function() {
    var mainArea = $(ID_MAIN_AREA);
    var loginBoxContainer = $(ID_LOGIN_BOX_CONTAINER);

    classRemove(mainArea, "blurred");
    classRemove(root, "root-min-height");

    attrSet(loginBoxContainer, "style", "display: none");
  }


  global.loginBoxShow = function(id) {
    var mainArea = $(ID_MAIN_AREA);
    var loginBoxContainer = $(ID_LOGIN_BOX_CONTAINER);

    global.popupRender(id);

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
        if (ID_ROOT === undefined) {
          ID_ROOT = "commento";
        }
      }
    }
  }


  function loadHash() {
    if (window.location.hash && window.location.hash.startsWith("#commento-")) {
      var el = $(ID_CARD + window.location.hash.split("-")[1]);
      if (el === null) {
        return;
      }

      classAdd(el, "highlighted-card");
      el.scrollIntoView(true);
    }
  }


  global.main = function(callback) {
    root = $(ID_ROOT);
    if (root === null) {
      console.log("[commento] error: no root element with ID '" + ID_ROOT + "' found");
      return;
    }

    classAdd(root, "root");

    loginBoxCreate();

    errorElementCreate();

    mainAreaCreate();

    selfGet(function() {
      commentsGet(function() {
        modToolsCreate();
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
    if (initted) {
      return;
    }
    initted = true;

    dataTagsLoad();

    if (autoInit === "true" || autoInit === undefined) {
      global.main(undefined);
    } else if (autoInit !== "false") {
      console.log("[commento] error: invalid value for data-auto-init; allowed values: true, false");
    }
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
