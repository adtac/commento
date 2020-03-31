(function(global, document) {
  "use strict";

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
  var ID_LOGIN_BOX_FORGOT_LINK_CONTAINER = "commento-login-box-forgot-link-container";
  var ID_LOGIN_BOX_LOGIN_LINK_CONTAINER = "commento-login-box-login-link-container";
  var ID_LOGIN_BOX_SSO_PRETEXT = "commento-login-box-sso-pretext";
  var ID_LOGIN_BOX_SSO_BUTTON_CONTAINER = "commento-login-box-sso-buttton-container";
  var ID_LOGIN_BOX_HR1 = "commento-login-box-hr1";
  var ID_LOGIN_BOX_OAUTH_PRETEXT = "commento-login-box-oauth-pretext";
  var ID_LOGIN_BOX_OAUTH_BUTTONS_CONTAINER = "commento-login-box-oauth-buttons-container";
  var ID_LOGIN_BOX_HR2 = "commento-login-box-hr2";
  var ID_MOD_TOOLS = "commento-mod-tools";
  var ID_MOD_TOOLS_LOCK_BUTTON = "commento-mod-tools-lock-button";
  var ID_ERROR = "commento-error";
  var ID_LOGGED_CONTAINER = "commento-logged-container";
  var ID_PRE_COMMENTS_AREA = "commento-pre-comments-area";
  var ID_COMMENTS_AREA = "commento-comments-area";
  var ID_SUPER_CONTAINER = "commento-textarea-super-container-";
  var ID_TEXTAREA_CONTAINER = "commento-textarea-container-";
  var ID_TEXTAREA = "commento-textarea-";
  var ID_ANONYMOUS_CHECKBOX = "commento-anonymous-checkbox-";
  var ID_SORT_POLICY = "commento-sort-policy-";
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
  var ID_MARKDOWN_BUTTON = "commento-markdown-button-";
  var ID_MARKDOWN_HELP = "commento-markdown-help-";
  var ID_FOOTER = "commento-footer";


  var origin = "[[[.Origin]]]";
  var cdn = "[[[.CdnPrefix]]]";
  var root = null;
  var pageId = parent.location.pathname;
  var cssOverride;
  var noFonts;
  var hideDeleted;
  var autoInit;
  var isAuthenticated = false;
  var comments = [];
  var commentsMap = {};
  var commenters = {};
  var requireIdentification = true;
  var isModerator = false;
  var isFrozen = false;
  var chosenAnonymous = false;
  var isLocked = false;
  var stickyCommentHex = "none";
  var shownReply = {};
  var shownEdit = {};
  var configuredOauths = {};
  var anonymousOnly = false;
  var popupBoxType = "login";
  var oauthButtonsShown = false;
  var sortPolicy = "score-desc";
  var selfHex = undefined;
  var mobileView = null;


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


  function removeAllEventListeners(node) {
    if (node !== null) {
      var replacement = node.cloneNode(true);
      if (node.parentNode !== null) {
        node.parentNode.replaceChild(replacement, node);
        return replacement;
      }
    }

    return node;
  }


  function onclick(node, f, arg) {
    node.addEventListener("click", function() {
      f(arg);
    }, false);
  }


  function onload(node, f, arg) {
    node.addEventListener("load", function() {
      f(arg);
    });
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
    isAuthenticated = false;
    isModerator = false;
    selfHex = undefined;
    refreshAll();
  }


  function profileEdit() {
    window.open(origin + "/profile?commenterToken=" + commenterTokenGet(), "_blank");
  }


  function notificationSettings(unsubscribeSecretHex) {
    window.open(origin + "/unsubscribe?unsubscribeSecretHex=" + unsubscribeSecretHex, "_blank");
  }


  function selfLoad(commenter, email) {
    commenters[commenter.commenterHex] = commenter;
    selfHex = commenter.commenterHex;

    var loggedContainer = create("div");
    var loggedInAs = create("div");
    var name;
    if (commenter.link !== "undefined") {
      name = create("a");
    } else {
      name = create("div");
    }
    var avatar;
    var notificationSettingsButton = create("div");
    var profileEditButton = create("div");
    var logoutButton = create("div");
    var color = colorGet(commenter.commenterHex + "-" + commenter.name);

    loggedContainer.id = ID_LOGGED_CONTAINER;

    classAdd(loggedContainer, "logged-container");
    classAdd(loggedInAs, "logged-in-as");
    classAdd(name, "name");
    classAdd(notificationSettingsButton, "profile-button");
    classAdd(profileEditButton, "profile-button");
    classAdd(logoutButton, "profile-button");

    name.innerText = commenter.name;
    notificationSettingsButton.innerText = "Notification Settings";
    profileEditButton.innerText = "Edit Profile";
    logoutButton.innerText = "Logout";

    onclick(logoutButton, global.logout);
    onclick(notificationSettingsButton, notificationSettings, email.unsubscribeSecretHex);
    onclick(profileEditButton, profileEdit);

    attrSet(loggedContainer, "style", "display: none");
    if (commenter.link !== "undefined") {
      attrSet(name, "href", commenter.link);
    }
    if (commenter.photo === "undefined") {
      avatar = create("div");
      avatar.style["background"] = color;
      avatar.innerHTML = commenter.name[0].toUpperCase();
      classAdd(avatar, "avatar");
    } else {
      avatar = create("img");
      attrSet(avatar, "src", cdn + "/api/commenter/photo?commenterHex=" + commenter.commenterHex);
      classAdd(avatar, "avatar-img");
    }

    append(loggedInAs, avatar);
    append(loggedInAs, name);
    append(loggedContainer, loggedInAs);
    append(loggedContainer, logoutButton);
    if (commenter.provider === "commento") {
      append(loggedContainer, profileEditButton);
    }
    append(loggedContainer, notificationSettingsButton);
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

      selfLoad(resp.commenter, resp.email);
      allShow();

      call(callback);
    });
  }


  function cssLoad(file, f) {
    var link = create("link");
    var head = document.getElementsByTagName("head")[0];

    link.type = "text/css";
    attrSet(link, "href", file);
    attrSet(link, "rel", "stylesheet");
    onload(link, f);

    append(head, link);
  }


  function footerLoad() {
    var footer = create("div");
    var aContainer = create("div");
    var a = create("a");
    var text = create("span");

    footer.id = ID_FOOTER;

    classAdd(footer, "footer");
    classAdd(aContainer, "logo-container");
    classAdd(a, "logo");
    classAdd(text, "logo-text");

    attrSet(a, "href", "https://commento.io");
    attrSet(a, "target", "_blank");

    text.innerText = "Commento";

    append(a, text);
    append(aContainer, a);
    append(footer, aContainer);

    return footer;
  }


  function commentsGet(callback) {
    var json = {
      "commenterToken": commenterTokenGet(),
      "domain": parent.location.host,
      "path": pageId,
    };

    post(origin + "/api/comment/list", json, function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return;
      } else {
        errorHide();
      }

      requireIdentification = resp.requireIdentification;
      isModerator = resp.isModerator;
      isFrozen = resp.isFrozen;

      isLocked = resp.attributes.isLocked;
      stickyCommentHex = resp.attributes.stickyCommentHex;

      comments = resp.comments;
      commenters = Object.assign({}, commenters, resp.commenters)
      configuredOauths = resp.configuredOauths;

      sortPolicy = resp.defaultSortPolicy;

      call(callback);
    });
  }


  function errorShow(text) {
    var el = $(ID_ERROR);

    el.innerText = text;

    attrSet(el, "style", "display: block;");
  }


  function errorHide() {
    var el = $(ID_ERROR);

    attrSet(el, "style", "display: none;");
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


  function markdownHelpShow(id) {
    var textareaSuperContainer = $(ID_SUPER_CONTAINER + id);
    var markdownButton = $(ID_MARKDOWN_BUTTON + id);
    var markdownHelp = create("table");
    var italicsContainer = create("tr");
    var italicsLeft = create("td");
    var italicsRight = create("td");
    var boldContainer = create("tr");
    var boldLeft = create("td");
    var boldRight = create("td");
    var codeContainer = create("tr");
    var codeLeft = create("td");
    var codeRight = create("td");
    var strikethroughContainer = create("tr");
    var strikethroughLeft = create("td");
    var strikethroughRight = create("td");
    var hyperlinkContainer = create("tr");
    var hyperlinkLeft = create("td");
    var hyperlinkRight = create("td");
    var quoteContainer = create("tr");
    var quoteLeft = create("td");
    var quoteRight = create("td");

    markdownHelp.id = ID_MARKDOWN_HELP + id;

    classAdd(markdownHelp, "markdown-help");

    boldLeft.innerHTML = "<b>bold</b>";
    boldRight.innerHTML = "surround text with <pre>**two asterisks**</pre>";
    italicsLeft.innerHTML = "<i>italics</i>";
    italicsRight.innerHTML = "surround text with <pre>*asterisks*</pre>";
    codeLeft.innerHTML = "<pre>code</pre>";
    codeRight.innerHTML = "surround text with <pre>`backticks`</pre>";
    strikethroughLeft.innerHTML = "<strike>strikethrough</strike>";
    strikethroughRight.innerHTML = "surround text with <pre>~~two tilde characters~~</pre>";
    hyperlinkLeft.innerHTML = "<a href=\"https://example.com\">hyperlink</a>";
    hyperlinkRight.innerHTML = "<pre>[hyperlink](https://example.com)</pre> or just a bare URL";
    quoteLeft.innerHTML = "<blockquote>quote</blockquote>";
    quoteRight.innerHTML = "prefix with <pre>&gt;</pre>";

    markdownButton = removeAllEventListeners(markdownButton);
    onclick(markdownButton, markdownHelpHide, id);

    append(italicsContainer, italicsLeft);
    append(italicsContainer, italicsRight);
    append(markdownHelp, italicsContainer);
    append(boldContainer, boldLeft);
    append(boldContainer, boldRight);
    append(markdownHelp, boldContainer);
    append(hyperlinkContainer, hyperlinkLeft);
    append(hyperlinkContainer, hyperlinkRight);
    append(markdownHelp, hyperlinkContainer);
    append(codeContainer, codeLeft);
    append(codeContainer, codeRight);
    append(markdownHelp, codeContainer);
    append(strikethroughContainer, strikethroughLeft);
    append(strikethroughContainer, strikethroughRight);
    append(markdownHelp, strikethroughContainer);
    append(quoteContainer, quoteLeft);
    append(quoteContainer, quoteRight);
    append(markdownHelp, quoteContainer);
    append(textareaSuperContainer, markdownHelp);
  }


  function markdownHelpHide(id) {
    var markdownButton = $(ID_MARKDOWN_BUTTON + id);
    var markdownHelp = $(ID_MARKDOWN_HELP + id);

    markdownButton = removeAllEventListeners(markdownButton);
    onclick(markdownButton, markdownHelpShow, id);

    remove(markdownHelp);
  }


  function textareaCreate(id, edit) {
    var textareaSuperContainer = create("div");
    var textareaContainer = create("div");
    var textarea = create("textarea");
    var anonymousCheckboxContainer = create("div");
    var anonymousCheckbox = create("input");
    var anonymousCheckboxLabel = create("label");
    var submitButton = create("button");
    var markdownButton = create("a");

    textareaSuperContainer.id = ID_SUPER_CONTAINER + id;
    textareaContainer.id = ID_TEXTAREA_CONTAINER + id;
    textarea.id = ID_TEXTAREA + id;
    anonymousCheckbox.id = ID_ANONYMOUS_CHECKBOX + id;
    submitButton.id = ID_SUBMIT_BUTTON + id;
    markdownButton.id = ID_MARKDOWN_BUTTON + id;

    classAdd(textareaContainer, "textarea-container");
    classAdd(anonymousCheckboxContainer, "round-check");
    classAdd(anonymousCheckboxContainer, "anonymous-checkbox-container");
    classAdd(submitButton, "button");
    classAdd(submitButton, "submit-button");
    classAdd(markdownButton, "markdown-button");
    classAdd(textareaSuperContainer, "button-margin");

    attrSet(textarea, "placeholder", "Add a comment");
    attrSet(anonymousCheckbox, "type", "checkbox");
    attrSet(anonymousCheckboxLabel, "for", ID_ANONYMOUS_CHECKBOX + id);

    anonymousCheckboxLabel.innerText = "Comment anonymously";
    if (edit === true) {
      submitButton.innerText = "Save Changes";
    } else {
      submitButton.innerText = "Add Comment";
    }
    markdownButton.innerHTML = "<b>M &#8595;</b> &nbsp; Markdown";

    if (anonymousOnly) {
      anonymousCheckbox.checked = true;
      anonymousCheckbox.setAttribute("disabled", true);
    }

    textarea.oninput = autoExpander(textarea);
    if (edit === true) {
      onclick(submitButton, commentEdit, id);
    } else {
      onclick(submitButton, submitAccountDecide, id);
    }
    onclick(markdownButton, markdownHelpShow, id);

    append(textareaContainer, textarea);
    append(textareaSuperContainer, textareaContainer);
    append(anonymousCheckboxContainer, anonymousCheckbox);
    append(anonymousCheckboxContainer, anonymousCheckboxLabel);
    append(textareaSuperContainer, submitButton);
    if (!requireIdentification && edit !== true) {
      append(textareaSuperContainer, anonymousCheckboxContainer);
    }
    append(textareaSuperContainer, markdownButton);

    return textareaSuperContainer;
  }


  var sortPolicyNames = {
    "score-desc": "Upvotes",
    "creationdate-desc": "Newest",
    "creationdate-asc": "Oldest",
  };


  function sortPolicyApply(policy) {
    classRemove($(ID_SORT_POLICY + sortPolicy), "sort-policy-button-selected");

    var commentsArea = $(ID_COMMENTS_AREA);
    commentsArea.innerHTML = "";
    sortPolicy = policy;
    var cards = commentsRecurse(parentMap(comments), "root");
    if (cards) {
      append(commentsArea, cards);
    }

    classAdd($(ID_SORT_POLICY + policy), "sort-policy-button-selected");
  }


  function sortPolicyBox() {
    var sortPolicyButtonsContainer = create("div");
    var sortPolicyButtons = create("div");

    classAdd(sortPolicyButtonsContainer, "sort-policy-buttons-container");
    classAdd(sortPolicyButtons, "sort-policy-buttons");

    for (var sp in sortPolicyNames) {
      var sortPolicyButton = create("a");
      sortPolicyButton.id = ID_SORT_POLICY + sp;
      classAdd(sortPolicyButton, "sort-policy-button");
      if (sp === sortPolicy) {
        classAdd(sortPolicyButton, "sort-policy-button-selected");
      }
      sortPolicyButton.innerText = sortPolicyNames[sp];
      onclick(sortPolicyButton, sortPolicyApply, sp);
      append(sortPolicyButtons, sortPolicyButton)
    }

    append(sortPolicyButtonsContainer, sortPolicyButtons);

    return sortPolicyButtonsContainer
  }


  function rootCreate(callback) {
    var login = create("div");
    var loginText = create("div");
    var mainArea = $(ID_MAIN_AREA);
    var preCommentsArea = create("div");
    var commentsArea = create("div");

    login.id = ID_LOGIN;
    preCommentsArea.id = ID_PRE_COMMENTS_AREA;
    commentsArea.id = ID_COMMENTS_AREA;

    classAdd(login, "login");
    classAdd(loginText, "login-text");
    classAdd(commentsArea, "comments");

    loginText.innerText = "Login";
    commentsArea.innerHTML = "";

    onclick(loginText, global.loginBoxShow, null);

    var numOauthConfigured = 0;
    Object.keys(configuredOauths).forEach(function(key) {
      if (configuredOauths[key]) {
        numOauthConfigured++;
      }
    });
    if (numOauthConfigured > 0) {
      append(login, loginText);
    } else if (!requireIdentification) {
      anonymousOnly = true;
    }

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

    if (comments.length > 0) {
      append(mainArea, sortPolicyBox());
    }

    append(mainArea, preCommentsArea);

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
      "commenterToken": commenterToken,
      "domain": parent.location.host,
      "path": pageId,
      "parentHex": id,
      "markdown": markdown,
    };

    post(origin + "/api/comment/new", json, function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return;
      } else {
        errorHide();
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

      var comment = {
        "commentHex": resp.commentHex,
        "commenterHex": commenterHex,
        "markdown": markdown,
        "html": resp.html,
        "parentHex": "root",
        "score": 0,
        "state": "approved",
        "direction": 0,
        "creationDate": new Date(),
      };

      var newCard = commentsRecurse({
        "root": [comment]
      }, "root");

      commentsMap[resp.commentHex] = comment;

      if (id !== "root") {
        textareaSuperContainer.replaceWith(newCard);

        shownReply[id] = false;

        classAdd(replyButton, "option-reply");
        classRemove(replyButton, "option-cancel");

        replyButton.title = "Reply to this comment";

        onclick(replyButton, global.replyShow, id)
      } else {
        textarea.value = "";
        insertAfter($(ID_PRE_COMMENTS_AREA), newCard);
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
      "#922428",
      "#6b4c9a",
      "#535154",
    ];

    var total = 0;
    for (var i = 0; i < name.length; i++) {
      total += name.charCodeAt(i);
    }
    var color = colors[total % colors.length];

    return color;
  }


  function timeDifference(current, previous) { // thanks stackoverflow
    // Times are defined in milliseconds
    var msPerSecond = 1000;
    var msPerMinute = 60 * msPerSecond;
    var msPerHour = 60 * msPerMinute;
    var msPerDay = 24 * msPerHour;
    var msPerMonth = 30 * msPerDay;
    var msPerYear = 12 * msPerMonth;

    // Time ago thresholds
    var msJustNow = 5 * msPerSecond; // Up until 5 s
    var msMinutesAgo = 2 * msPerMinute; // Up until 2 minutes
    var msHoursAgo = 2 * msPerHour; // Up until 2 hours
    var msDaysAgo = 2 * msPerDay; // Up until 2 days
    var msMonthsAgo = 2 * msPerMonth; // Up until 2 months
    var msYearsAgo = 2 * msPerYear; // Up until 2 years

    var elapsed = current - previous;

    if (elapsed < msJustNow) {
      return "just now";
    } else if (elapsed < msMinutesAgo) {
      return Math.round(elapsed / msPerSecond) + " seconds ago";
    } else if (elapsed < msHoursAgo) {
      return Math.round(elapsed / msPerMinute) + " minutes ago";
    } else if (elapsed < msDaysAgo ) {
      return Math.round(elapsed / msPerHour ) + " hours ago";
    } else if (elapsed < msMonthsAgo) {
      return Math.round(elapsed / msPerDay) + " days ago";
    } else if (elapsed < msYearsAgo) {
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


  var sortPolicyFunctions = {
    "score-desc": function(a, b) {
      return b.score - a.score;
    },
    "creationdate-desc": function(a, b) {
      if (a.creationDate < b.creationDate) {
        return 1;
      } else {
        return -1;
      }
    },
    "creationdate-asc": function(a, b) {
      if (a.creationDate < b.creationDate) {
        return -1;
      } else {
        return 1;
      }
    },
  };


  function commentsRecurse(parentMap, parentHex) {
    var cur = parentMap[parentHex];
    if (!cur || !cur.length) {
      return null;
    }

    cur.sort(function(a, b) {
      if (!a.deleted && a.commentHex === stickyCommentHex) {
        return -Infinity;
      } else if (!b.deleted && b.commentHex === stickyCommentHex) {
        return Infinity;
      }

      return sortPolicyFunctions[sortPolicy](a, b);
    });

    var curTime = (new Date()).getTime();
    var cards = create("div");
    cur.forEach(function(comment) {
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
      timeago.title = comment.creationDate.toString();

      card.style["borderLeft"] = "2px solid " + color;
      if (comment.deleted) {
        name.innerText = "[deleted]";
      } else {
        name.innerText = commenter.name;
      }
      text.innerHTML = comment.html;
      timeago.innerHTML = timeDifference(curTime, comment.creationDate);
      score.innerText = scorify(comment.score);

      if (commenter.photo === "undefined") {
        avatar = create("div");
        avatar.style["background"] = color;

        if (comment.commenterHex === "anonymous") {
          avatar.innerHTML = "?";
          avatar.style["font-weight"] = "bold";
        } else {
          avatar.innerHTML = commenter.name[0].toUpperCase();
        }

        classAdd(avatar, "avatar");
      } else {
        avatar = create("img");
        attrSet(avatar, "src", cdn + "/api/commenter/photo?commenterHex=" + commenter.commenterHex);
        classAdd(avatar, "avatar-img");
      }

      classAdd(card, "card");
      if (isModerator && comment.state !== "approved") {
        classAdd(card, "dark-card");
      }
      if (commenter.isModerator) {
        classAdd(name, "moderator");
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

      onclick(edit, global.editShow, comment.commentHex);
      onclick(collapse, global.commentCollapse, comment.commentHex);
      onclick(approve, global.commentApprove, comment.commentHex);
      onclick(remove, global.commentDelete, comment.commentHex);
      onclick(sticky, global.commentSticky, comment.commentHex);

      if (isAuthenticated) {
        var upDown = upDownOnclickSet(upvote, downvote, comment.commentHex, comment.direction);
        upvote = upDown[0];
        downvote = upDown[1];
      } else {
        onclick(upvote, global.loginBoxShow, null);
        onclick(downvote, global.loginBoxShow, null);
      }

      onclick(reply, global.replyShow, comment.commentHex);

      if (commenter.link !== "undefined" && commenter.link !== "https://undefined" && commenter.link !== "") {
        attrSet(name, "href", commenter.link);
      }

      append(options, collapse);

      if (!comment.deleted) {
        append(options, downvote);
        append(options, upvote);
      }

      if (comment.commenterHex === selfHex) {
        append(options, edit);
      } else if (!comment.deleted) {
        append(options, reply);
      }

      if (!comment.deleted && (isModerator && parentHex === "root")) {
        append(options, sticky);
      }

      if (!comment.deleted && (isModerator || comment.commenterHex === selfHex)) {
        append(options, remove);
      }

      if (isModerator && comment.state !== "approved") {
        append(options, approve);
      }
      
      if (!comment.deleted && (!isModerator && stickyCommentHex === comment.commentHex)) {
        append(options, sticky);
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

      if (comment.deleted && (hideDeleted === "true" || children === null)) {
        return;
      }

      append(cards, card);
    });

    if (cards.childNodes.length === 0) {
      return null;
    }

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
      } else {
        errorHide();
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
    if (!confirm("Are you sure you want to delete this comment?")) {
      return;
    }

    var json = {
      "commenterToken": commenterTokenGet(),
      "commentHex": commentHex,
    }

    post(origin + "/api/comment/delete", json, function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return
      } else {
        errorHide();
      }

      var text = $(ID_TEXT + commentHex);
      text.innerText = "[deleted]";
    });
  }


  function nameWidthFix() {
    var els = document.getElementsByClassName("commento-name");

    for (var i = 0; i < els.length; i++) {
      attrSet(els[i], "style", "max-width: " + (els[i].getBoundingClientRect()["width"] + 20) + "px;")
    }
  }


  function upDownOnclickSet(upvote, downvote, commentHex, direction) {
    upvote = removeAllEventListeners(upvote);
    downvote = removeAllEventListeners(downvote);

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

    return [upvote, downvote];
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

    var upDown = upDownOnclickSet(upvote, downvote, commentHex, newDirection);
    upvote = upDown[0];
    downvote = upDown[1];

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
        classRemove(upvote, "upvoted");
        classRemove(downvote, "downvoted");
        score.innerText = scorify(parseInt(score.innerText.replace(/[^\d-.]/g, "")) - newDirection + oldDirection);
        upDownOnclickSet(upvote, downvote, commentHex, oldDirection);
        return;
      } else {
        errorHide();
      }
    });
  }


  function commentEdit(id) {
    var textarea = $(ID_TEXTAREA + id);

    var markdown = textarea.value;

    if (markdown === "") {
      classAdd(textarea, "red-border");
      return;
    } else {
      classRemove(textarea, "red-border");
    }

    var json = {
      "commenterToken": commenterTokenGet(),
      "commentHex": id,
      "markdown": markdown,
    };

    post(origin + "/api/comment/edit", json, function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return;
      } else {
        errorHide();
      }

      commentsMap[id].markdown = markdown;
      commentsMap[id].html = resp.html;

      var editButton = $(ID_EDIT + id);
      var textarea = $(ID_SUPER_CONTAINER + id);

      textarea.innerHTML = commentsMap[id].html;
      textarea.id = ID_TEXT + id;
      delete shownEdit[id];

      classAdd(editButton, "option-edit");
      classRemove(editButton, "option-cancel");

      editButton.title = "Edit comment";

      editButton = removeAllEventListeners(editButton);
      onclick(editButton, global.editShow, id)

      var message = "";
      if (resp.state === "unapproved") {
        message = "Your comment is under moderation.";
      } else if (resp.state === "flagged") {
        message = "Your comment was flagged as spam and is under moderation.";
      }

      if (message !== "") {
        prepend($(ID_SUPER_CONTAINER + id), messageCreate(message));
      }
    });
  }


  global.editShow = function(id) {
    if (id in shownEdit && shownEdit[id]) {
      return;
    }

    var text = $(ID_TEXT + id);
    shownEdit[id] = true;
    text.replaceWith(textareaCreate(id, true));

    var textarea = $(ID_TEXTAREA + id);
    textarea.value = commentsMap[id].markdown;

    var editButton = $(ID_EDIT + id);

    classRemove(editButton, "option-edit");
    classAdd(editButton, "option-cancel");

    editButton.title = "Cancel edit";

    editButton = removeAllEventListeners(editButton);
    onclick(editButton, global.editCollapse, id);
  };


  global.editCollapse = function(id) {
    var editButton = $(ID_EDIT + id);
    var textarea = $(ID_SUPER_CONTAINER + id);

    textarea.innerHTML = commentsMap[id].html;
    textarea.id = ID_TEXT + id;
    delete shownEdit[id];

    classAdd(editButton, "option-edit");
    classRemove(editButton, "option-cancel");

    editButton.title = "Edit comment";

    editButton = removeAllEventListeners(editButton);
    onclick(editButton, global.editShow, id)
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

    replyButton = removeAllEventListeners(replyButton);
    onclick(replyButton, global.replyCollapse, id);
  };


  global.replyCollapse = function(id) {
    var replyButton = $(ID_REPLY + id);
    var el = $(ID_SUPER_CONTAINER + id);

    el.remove();
    delete shownReply[id];

    classAdd(replyButton, "option-reply");
    classRemove(replyButton, "option-cancel");

    replyButton.title = "Reply to this comment";

    replyButton = removeAllEventListeners(replyButton);
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

    button = removeAllEventListeners(button);
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

    button = removeAllEventListeners(button);
    onclick(button, global.commentCollapse, id);
  }


  function parentMap(comments) {
    var m = {};
    comments.forEach(function(comment) {
      var parentHex = comment.parentHex;
      if (!(parentHex in m)) {
        m[parentHex] = [];
      }

      comment.creationDate = new Date(comment.creationDate);

      m[parentHex].push(comment);
      commentsMap[comment.commentHex] = {
        "html": comment.html,
        "markdown": comment.markdown,
      };
    });

    return m;
  }


  function commentsRender() {
    var commentsArea = $(ID_COMMENTS_AREA);
    commentsArea.innerHTML = ""

    var cards = commentsRecurse(parentMap(comments), "root");
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


  // OAuth logic
  global.commentoAuth = function(data) {
    var provider = data.provider;
    var id = data.id;
    var popup = window.open("", "_blank");

    get(origin + "/api/commenter/token/new", function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return;
      } else {
        errorHide();
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

            if (commenterTokenGet() !== "anonymous") {
              remove($(ID_LOGIN));
            }

            if (id !== null) {
              global.commentNew(id, resp.commenterToken, function() {
                global.loginBoxClose();
                commentsGet(commentsRender);
              });
            } else {
              global.loginBoxClose();
              commentsGet(commentsRender);
            }
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
    var ssoSubtitle = create("div");
    var ssoButtonContainer = create("div");
    var ssoButton = create("div");
    var hr1 = create("hr");
    var oauthSubtitle = create("div");
    var oauthButtonsContainer = create("div");
    var oauthButtons = create("div");
    var hr2 = create("hr");
    var emailSubtitle = create("div");
    var emailContainer = create("div");
    var email = create("div");
    var emailInput = create("input");
    var emailButton = create("button");
    var forgotLinkContainer = create("div");
    var forgotLink = create("a");
    var loginLinkContainer = create("div");
    var loginLink = create("a");
    var close = create("div");

    loginBox.id = ID_LOGIN_BOX;
    emailSubtitle.id = ID_LOGIN_BOX_EMAIL_SUBTITLE;
    emailInput.id = ID_LOGIN_BOX_EMAIL_INPUT;
    emailButton.id = ID_LOGIN_BOX_EMAIL_BUTTON;
    forgotLinkContainer.id = ID_LOGIN_BOX_FORGOT_LINK_CONTAINER
    loginLinkContainer.id = ID_LOGIN_BOX_LOGIN_LINK_CONTAINER;
    ssoButtonContainer.id = ID_LOGIN_BOX_SSO_BUTTON_CONTAINER;
    ssoSubtitle.id = ID_LOGIN_BOX_SSO_PRETEXT;
    hr1.id = ID_LOGIN_BOX_HR1;
    oauthSubtitle.id = ID_LOGIN_BOX_OAUTH_PRETEXT;
    oauthButtonsContainer.id = ID_LOGIN_BOX_OAUTH_BUTTONS_CONTAINER;
    hr2.id = ID_LOGIN_BOX_HR2;

    classAdd(loginBoxContainer, "login-box-container");
    classAdd(loginBox, "login-box");
    classAdd(emailSubtitle, "login-box-subtitle");
    classAdd(emailContainer, "email-container");
    classAdd(email, "email");
    classAdd(emailInput, "input");
    classAdd(emailButton, "email-button");
    classAdd(forgotLinkContainer, "forgot-link-container");
    classAdd(forgotLink, "forgot-link");
    classAdd(loginLinkContainer, "login-link-container");
    classAdd(loginLink, "login-link");
    classAdd(ssoSubtitle, "login-box-subtitle");
    classAdd(ssoButtonContainer, "oauth-buttons-container");
    classAdd(ssoButton, "oauth-buttons");
    classAdd(oauthSubtitle, "login-box-subtitle");
    classAdd(oauthButtonsContainer, "oauth-buttons-container");
    classAdd(oauthButtons, "oauth-buttons");
    classAdd(close, "login-box-close");
    classAdd(root, "root-min-height");

    forgotLink.innerText = "Forgot your password?";
    loginLink.innerText = "Don't have an account? Sign up.";
    emailSubtitle.innerText = "Login with your email address";
    emailButton.innerText = "Continue";
    oauthSubtitle.innerText = "Proceed with social login";
    ssoSubtitle.innerText = "Proceed with " + parent.location.host + " authentication";

    onclick(emailButton, global.passwordAsk, id);
    onclick(forgotLink, global.forgotPassword, id);
    onclick(loginLink, global.popupSwitch, id);
    onclick(close, global.loginBoxClose);

    attrSet(loginBoxContainer, "style", "display: none; opacity: 0;");
    attrSet(emailInput, "name", "email");
    attrSet(emailInput, "placeholder", "Email address");
    attrSet(emailInput, "type", "text");

    var numOauthConfigured = 0;
    var oauthProviders = ["google", "twitter", "github", "gitlab"];
    oauthProviders.forEach(function(provider) {
      if (configuredOauths[provider]) {
        var button = create("button");

        classAdd(button, "button");
        classAdd(button, provider + "-button");

        button.innerText = provider;

        onclick(button, global.commentoAuth, {"provider": provider, "id": id});

        append(oauthButtons, button);
        numOauthConfigured++;
      }
    });

    if (configuredOauths["sso"]) {
      var button = create("button");

      classAdd(button, "button");
      classAdd(button, "sso-button");

      button.innerText = "Single Sign-On";

      onclick(button, global.commentoAuth, {"provider": "sso", "id": id});

      append(ssoButton, button);
      append(ssoButtonContainer, ssoButton);
      append(loginBox, ssoSubtitle);
      append(loginBox, ssoButtonContainer);

      if (numOauthConfigured > 0 || configuredOauths["commento"]) {
        append(loginBox, hr1);
      }
    }

    if (numOauthConfigured > 0) {
      append(loginBox, oauthSubtitle);
      append(oauthButtonsContainer, oauthButtons);
      append(loginBox, oauthButtonsContainer);
      oauthButtonsShown = true;
    } else {
      oauthButtonsShown = false;
    }

    append(email, emailInput);
    append(email, emailButton);
    append(emailContainer, email);

    append(forgotLinkContainer, forgotLink);

    append(loginLinkContainer, loginLink);

    if (numOauthConfigured > 0 && configuredOauths["commento"]) {
      append(loginBox, hr2);
    }

    if (configuredOauths["commento"]) {
      append(loginBox, emailSubtitle);
      append(loginBox, emailContainer);
      append(loginBox, forgotLinkContainer);
      append(loginBox, loginLinkContainer);
    }

    append(loginBox, close);

    popupBoxType = "login";
    loginBoxContainer.innerHTML = "";
    append(loginBoxContainer, loginBox);
  }


  global.forgotPassword = function() {
    var popup = window.open("", "_blank");
    popup.location = origin + "/forgot?commenter=true";
    global.loginBoxClose();
  }


  global.popupSwitch = function(id) {
    var emailSubtitle = $(ID_LOGIN_BOX_EMAIL_SUBTITLE);

    if (oauthButtonsShown) {
      remove($(ID_LOGIN_BOX_OAUTH_BUTTONS_CONTAINER));
      remove($(ID_LOGIN_BOX_OAUTH_PRETEXT));
      remove($(ID_LOGIN_BOX_HR1));
      remove($(ID_LOGIN_BOX_HR2));
    }

    if (configuredOauths["sso"]) {
      remove($(ID_LOGIN_BOX_SSO_BUTTON_CONTAINER));
      remove($(ID_LOGIN_BOX_SSO_PRETEXT));
      remove($(ID_LOGIN_BOX_HR1));
      remove($(ID_LOGIN_BOX_HR2));
    }

    remove($(ID_LOGIN_BOX_LOGIN_LINK_CONTAINER));
    remove($(ID_LOGIN_BOX_FORGOT_LINK_CONTAINER));

    emailSubtitle.innerText = "Create an account";
    popupBoxType = "signup";
    global.passwordAsk(id);
    $(ID_LOGIN_BOX_EMAIL_INPUT).focus();
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
      } else {
        errorHide();
      }

      cookieSet("commentoCommenterToken", resp.commenterToken);

      selfLoad(resp.commenter, resp.email);
      allShow();

      remove($(ID_LOGIN));
      if (id !== null) {
        global.commentNew(id, resp.commenterToken, function() {
          global.loginBoxClose();
          commentsGet(commentsRender);
        });
      } else {
        global.loginBoxClose();
        commentsGet(commentsRender);
      }
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
      } else {
        errorHide();
      }

      loginUP(email.value, password.value, id);
    });
  }


  global.passwordAsk = function(id) {
    var loginBox = $(ID_LOGIN_BOX);
    var subtitle = $(ID_LOGIN_BOX_EMAIL_SUBTITLE);
    
    remove($(ID_LOGIN_BOX_EMAIL_BUTTON));
    remove($(ID_LOGIN_BOX_LOGIN_LINK_CONTAINER));
    remove($(ID_LOGIN_BOX_FORGOT_LINK_CONTAINER));
    if (oauthButtonsShown) {
      if (configuredOauths.length > 0) {
        remove($(ID_LOGIN_BOX_HR1));
        remove($(ID_LOGIN_BOX_HR2));
        remove($(ID_LOGIN_BOX_OAUTH_PRETEXT));
        remove($(ID_LOGIN_BOX_OAUTH_BUTTONS_CONTAINER));
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
      "path": pageId,
      "attributes": attributes,
    };

    post(origin + "/api/page/update", json, function(resp) {
      if (!resp.success) {
        errorShow(resp.message);
        return
      } else {
        errorHide();
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


  function loadCssOverride() {
    if (cssOverride === undefined) {
      allShow();
    } else {
      cssLoad(cssOverride, allShow);
    }
  }


  function allShow() {
    var mainArea = $(ID_MAIN_AREA);
    var modTools = $(ID_MOD_TOOLS);
    var loggedContainer = $(ID_LOGGED_CONTAINER);

    attrSet(mainArea, "style", "");

    if (isModerator) {
      attrSet(modTools, "style", "");
    }

    if (loggedContainer) {
      attrSet(loggedContainer, "style", "");
    }
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
        var pid = attrGet(scripts[i], "data-page-id");
        if (pid !== undefined) {
          pageId = pid;
        }

        cssOverride = attrGet(scripts[i], "data-css-override");

        autoInit = attrGet(scripts[i], "data-auto-init");

        ID_ROOT = attrGet(scripts[i], "data-id-root");
        if (ID_ROOT === undefined) {
          ID_ROOT = "commento";
        }

        noFonts = attrGet(scripts[i], "data-no-fonts");

        hideDeleted = attrGet(scripts[i], "data-hide-deleted");
      }
    }
  }


  function loadHash() {
    if (window.location.hash) {
      if (window.location.hash.startsWith("#commento-")) {
        var el = $(ID_CARD + window.location.hash.split("-")[1]);
        if (el === null) {
          return;
        }

        classAdd(el, "highlighted-card");
        el.scrollIntoView(true);
      } else if (window.location.hash.startsWith("#commento")) {
        root.scrollIntoView(true);
      }
    }
  }


  global.main = function(callback) {
    root = $(ID_ROOT);
    if (root === null) {
      console.log("[commento] error: no root element with ID '" + ID_ROOT + "' found");
      return;
    }

    if (mobileView === null) {
      mobileView = root.getBoundingClientRect()["width"] < 450;
    }

    classAdd(root, "root");
    if (noFonts !== "true") {
      classAdd(root, "root-font");
    }

    loginBoxCreate();

    errorElementCreate();

    mainAreaCreate();

    var footer = footerLoad();
    cssLoad(cdn + "/css/commento.css", loadCssOverride);

    selfGet(function() {
      commentsGet(function() {
        modToolsCreate();
        rootCreate(function() {
          commentsRender();
          append(root, footer);
          loadHash();
          allShow();
          nameWidthFix();
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


}(window.commento, document));
