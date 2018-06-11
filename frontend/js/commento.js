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
  var configuredOauths = [];
  var loginBoxType = "signup";


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
      if (!resp.success) {
        cookieSet("session", "anonymous");
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

      attr(loggedContainer, "style", "display: none");
      attr(logout, "onclick", "logout()");
      attr(name, "href", resp.commenter.link);
      if (resp.commenter.photo == "undefined") {
        avatar = create("div");
        avatar.style["background"] = color;
        avatar.style["boxShadow"] = "0px 0px 0px 2px " + color;
        avatar.innerHTML = resp.commenter.name[0].toUpperCase();
        classAdd(avatar, "avatar");
      } else {
        avatar = create("img");
        if (resp.commenter.provider == "google") {
          attr(avatar, "src", resp.commenter.photo + "?sz=50");
        } else {
          attr(avatar, "src", resp.commenter.photo);
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

  function cssLoad(file) {
    var link = create("link");
    var head = document.getElementsByTagName('head')[0];

    link.type = "text/css";
    attr(link, "href", file);
    attr(link, "rel", "stylesheet");
    attr(link, "onload", "window.allShow()");

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

    footer.id = ID_FOOTER;

    classAdd(footer, "footer");
    classAdd(aContainer, "logo-container");
    classAdd(a, "logo");
    classAdd(img, "logo-svg");
    classAdd(text, "logo-text");

    attr(footer, "style", "display: none");
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
      configuredOauths = resp.configuredOauths;

      if (!resp.requireIdentification)
        configuredOauths.push("anonymous");

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
      var createButton = create("div");

      classAdd(buttonsContainer, "create-container");

      classAdd(createButton, "button");
      classAdd(createButton, "create-button");

      attr(createButton, "onclick", "loginBoxShow()");
      attr(textarea, "disabled", true);

      createButton.innerText = "Create an Account";

      append(buttonsContainer, createButton);
      append(textareaContainer, buttonsContainer);
    }
    else {
      attr(textarea, "placeholder", "Join the discussion!");
      attr(textarea, "onclick", "showSubmitButton('" + id + "')");
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
      }
      else {
        attr(upvote, "onclick", "loginBoxShow()");
        attr(downvote, "onclick", "loginBoxShow()");
      }

      if (isAuthenticated || chosenAnonymous)
        attr(reply, "onclick", "replyShow('" + comment.commentHex + "')");
      else
        attr(reply, "onclick", "loginBoxShow()");

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

      append(options, upvote);
      append(options, downvote);
      append(options, reply);

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

  global.commentoAuth = function(provider) {
    if (provider == "anonymous") {
      cookieSet("session", "anonymous");
      chosenAnonymous = true;
      refreshAll();
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

    attr(loginBoxContainer, "style", "display: none; opacity: 0;");
    attr(emailInput, "name", "email");
    attr(emailInput, "placeholder", "Email address");
    attr(emailInput, "type", "text");
    attr(emailButton, "onclick", "passwordAsk()");
    attr(loginLink, "onclick", "loginSwitch()");
    attr(close, "onclick", "loginBoxClose()");

    for (var i = 0; i < configuredOauths.length; i++) {
      var button = create("button");

      classAdd(button, "button");
      classAdd(button, configuredOauths[i] + "-button");

      button.innerText = configuredOauths[i];

      attr(button, "onclick", "commentoAuth('" + configuredOauths[i] + "')");

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

    attr(loginLink, "onclick", "signupSwitch()");

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
      email: username,
      password: password,
    };

    post(origin + "/api/commenter/login", json, function(resp) {
      if (!resp.success) {
        loginBoxClose();
        errorShow(resp.message);
        return
      }

      cookieSet("session", resp.session);
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
      email: email.value,
      name: name.value,
      website: website.value,
      password: password.value,
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

      attr(fieldInput, "name", order[i]);
      attr(fieldInput, "type", type[i]);
      attr(fieldInput, "placeholder", placeholder[i]);

      append(field, fieldInput);
      append(fieldContainer, field);

      if (order[i] == "password") {
        var fieldButton = create("button");
        classAdd(fieldButton, "email-button");
        fieldButton.innerText = loginBoxType;

        if (loginBoxType == "signup")
          attr(fieldButton, "onclick", "signup()");
        else
          attr(fieldButton, "onclick", "login()");

        append(field, fieldButton);
      }

      append(loginBox, fieldContainer);
    }
  }

  function mainAreaCreate() {
    var mainArea = create("div");

    mainArea.id = ID_MAIN_AREA;

    classAdd(mainArea, "main-area");

    attr(mainArea, "style", "display: none");

    append(root, mainArea);
  }

  global.allShow = function() {
    var mainArea = $(ID_MAIN_AREA);
    var loggedContainer = $(ID_LOGGED_CONTAINER);
    var footer = $(ID_FOOTER);

    attr(mainArea, "style", "");
    if (loggedContainer)
      attr(loggedContainer, "style", "");
    attr(footer, "style", "");

    nameWidthFix();
  }

  global.loginBoxClose = function() {
    var mainArea = $(ID_MAIN_AREA);
    var loginBoxContainer = $(ID_LOGIN_BOX_CONTAINER);

    classRemove(mainArea, "blurred");

    attr(loginBoxContainer, "style", "display: none");
  }

  global.loginBoxShow = function() {
    var mainArea = $(ID_MAIN_AREA);
    var loginBoxContainer = $(ID_LOGIN_BOX_CONTAINER);

    global.signupRender();

    classAdd(mainArea, "blurred");
    
    attr(loginBoxContainer, "style", "");

    window.location.hash = ID_LOGIN_BOX_CONTAINER;

    $(ID_LOGIN_BOX_EMAIL_INPUT).focus();
  }

  function main(callback) {
    root = $(ID_ROOT);

    loginBoxCreate();

    errorElementCreate();

    mainAreaCreate();

    selfGet(function() {
      commentsGet(function() {
        rootCreate(function() {
          commentsRender();
          footerLoad();
          attr(root, "style", "");
          call(callback);
        });
      });
    });
  }

  document.addEventListener("DOMContentLoaded", main);

}(window, document));
