(function(global, document) {
    'use strict';

    /**
     * Helper function /////////////////////////////
     *
     * This pattern helps with minification, because it allows Uglify to mangle function names. If we were to use
     * domprops we'd probably have to use --mangle-props, which in turn breaks the script and is hard to maintain.
     * We can rely on gzip for the rest of the repetitions.
     */
    function $(id) {
        return document.getElementById(id);
    }

    function append(root, el) {
        root.appendChild(el);
    }

    function getData(el, key){
        return el.dataset[key];
    }

    function setData(el, key, data){
        el.dataset[key] = data;
    }

    function contains(el, cls){
        return el.classList.contains(cls);
    }

    function addClass(el, cls) {
        el.classList.add(cls);
    }

    function removeClass(el, cls) {
        el.classList.remove(cls);
    }

    function create(el) {
        return document.createElement(el);
    }

    function serialize(obj) {
        var str = [];
        for(var p in obj) {
            if (obj.hasOwnProperty(p)) {
                str.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
            }
        }

        return str.join("&");
    }

    function setAttr(node, attr, value) {
        node.setAttribute(attr, value);
    }

    function post(url, data, callback) {
        var xmlDoc = new XMLHttpRequest();
        xmlDoc.open('POST', url, true);
        xmlDoc.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
        xmlDoc.onreadystatechange = function() {
            if (xmlDoc.readyState === 4 && xmlDoc.status === 200) {
                callback(xmlDoc);
            }
        };
        xmlDoc.send(serialize(data));
    }

    function startSpinner(id) {
        var button = $(id);

        button.disabled = true;
        addClass(button, "loading");
    }

    function clearSpinner(id) {
        var button = $(id);

        button.disabled = false;
        removeClass(button, "loading");
    }

    var c = 0;

    function getColor(name) {
        var colors = [
            "#c10000",
            "#c16a00",
            "#005b0c",
            "#009385",
            "#006192",
            "#3a00a8",
            "#8e00a8",
            "#e50076",
            "#383838",
        ];

        var total = 0;
        for (var i = 0; i < name.length; i++)
            total += name.charCodeAt(i);
        var color = colors[total % colors.length];

        return color;
    }

    function loadJS(file, ready) {
        var script = document.createElement("script");
        var loaded = false;

        script.type = "application/javascript";
        script.src = file;
        script.async = true;
        script.onreadystatechange = script.onload = function() {
            if(!loaded &&
                (!this.readyState ||
                    this.readyState === "loaded" ||
                    this.readyState === "complete"))
            {
                ready();
            }

            loaded = true;
            script.onload = script.onreadystatechange = null;
        };

        append(document.body, script);
    }

    function loadCSS(file) {
        var link = create("link");
        var head = document.getElementsByTagName('head')[0];

        link.type = "text/css";
        setAttr(link, "href", file);
        setAttr(link, "rel", "stylesheet");

        append(head, link);
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
            return 'approximately ' + Math.round(elapsed/msPerDay) + ' days ago';
        }
        else if (elapsed < msPerYear) {
            return 'approximately ' + Math.round(elapsed/msPerMonth) + ' months ago';
        }
        else {
            return 'approximately ' + Math.round(elapsed/msPerYear ) + ' years ago';
        }
    }

    function autoExpander(el) {
        return function() {
            el.style.height = "";
            el.style.height = Math.min(Math.max(el.scrollHeight, 75), 400) + "px";
        }
    };

    function addEvent(el, event, fn){
        el.addEventListener(event, fn);
    }

    function makeEvent(cls, key, fn){
        return function onEvent(event){
            var cur = event.target;
            var id;

            if( !( contains(cur, cls) &&
                ( id = getData(cur, key) ))){
                return;
            }

            fn(id);
        };
    }

    /**
     * Constants ///////////////////////////////////////////
     *
     * This section allows to configure Commento internally, also helps to mangle constant names.
     *
     * The convention is:
     *
     * _ID for any string that relates to a node.id
     * _CLASS for any visual class that it's defined in a CSS file
     * _JS for any JavaScript hook used to delegate events or trigger actions.
     * _DATA for any data key we wish to identify in dataset
     */

    var ROOT_COMMENT_ID = "commento-root-comment";
    var ROOT_NAME_ID = "commento-root-name";
    var ROOT_BTN_ID = "commento-root-post"
    var COMS_ID = "commento-coms";
    var HONEYPOT_ID = "commento-root-gotcha";
    var TEXTAREA_ID = "commento-textarea-";
    var CONTENTS_ID = "comment-contents-";
    var NAME_INPUT_ID = "commento-name-input-";
    var GOTCHA_ID = "commento-gotcha-";
    var CANCEL_BTN_ID = "commento-cancel-button-";
    var SUBMIT_BTN_ID = "commento-submit-button-";
    var BODY_ID = "commento-body-";
    var REPLY_BTN_ID = "commento-reply-button-";
    var COLLAPSE_BTN_ID = "commento-collapse-button-";
    var COMMENTO_ID = "commento";

    var CANCEL_JS = 'commento-js-cancel';
    var SHOW_REPLY_JS = 'commento-show-reply-js';
    var COLLAPSE_THREAD_JS = 'commento-collapse-thread-js';
    var EXPAND_THREAD_JS = 'commento-expand-thread-js';
    var SUBMIT_JS = 'commento-submit-js';

    var CARD_CLASS = 'card';
    var CARD_HEADER_CLASS = 'card-header';
    var CARD_AVATAR_CLASS = 'card-avatar';
    var CARD_NAME_CLASS = 'name-header';
    var CARD_OPTIONS_CLASS = 'card-options';
    var CARD_SUBTITLE_CLASS = 'card-subtitle';
    var CARD_BODY_CLASS = 'card-body';
    var OPTION_BTN_CLASS = 'option-btn'
    var POST_BTN_CLASS = 'post-btn';
    var POST_PRIMARY_BTN_CLASS = 'post-primary-btn';
    var POST_RED_BTN_CLASS = 'post-red-btn';
    var HIDDEN_CLASS= 'hidden';
    var OTHER_FIELDS_CLASS = 'other-fields';
    var OTHER_FIELDS_CONTAINER_CLASS = 'other-fields-container';
    var INITIAL_CLASS = 'initial';
    var TEXTAREA_CLASS = "textarea";
    var INPUT_CLASS = "input";
    var IS_ERROR_CLASS = 'is-error';

    var COMMENT_ID_DATA = 'commentId';

    /**
     * Private fields ///////////////////////////////////////////
     *
     * This section is for internal use. This hold Commento internals that shouldn't be accessible from the outside
     */

    var _showdownUrl = "/assets/vendor/showdown.min.js";
    var _commentoCssUrl = "/assets/style/commento.min.css";
    var _serverUrl = '';
    var _honeypot = false;
    var _api = {};
    var _showdownConverter;


    var _getComments = function(callback) {
        var data = {
            "url": document.location
        };
        post(_api.get, data, function(reply) {
            var response = {
                comments: []
            };

            try {
                response = JSON.parse(reply.response)
            } catch (e){}

            _redraw(response.comments);

            if (typeof(callback) == "function")
                callback();
        });
    };

    var _makeCards = function(parentMap, cur) {
        var currentParent = parentMap[cur];
        if(!currentParent || !currentParent.length) {
            return null;
        }

        var cards = create("div");
        currentParent.forEach(function(comment) {
            var card = create("div");
            var header = create("div");
            var avatar = create("div");
            var name = create("div");
            var subtitle = create("div");
            var body = create("div");
            var options = create("div");
            var reply = create("button");
            var collapse = create("button");
            var children = _makeCards(parentMap, comment.id);
            var contents = create("div");
            var color = getColor(comment.name);

            body.id = BODY_ID + comment.id;
            reply.id = REPLY_BTN_ID + comment.id;
            collapse.id = COLLAPSE_BTN_ID + comment.id;
            contents.id = CONTENTS_ID + comment.id;

            card.style["borderLeft"] = "3px solid " + color;
            name.innerHTML = comment.name;
            avatar.style["background"] = color;
            avatar.style["boxShadow"] = "0px 0px 0px 2px " + color;
            avatar.innerHTML = comment.name[0].toUpperCase();
            body.innerHTML = _showdownConverter.makeHtml(comment.comment);
            subtitle.innerHTML = timeDifference(Date.now(), Date.parse(comment.timestamp));
            reply.innerHTML = "Reply";
            collapse.innerHTML = "Collapse";

            setData(reply, COMMENT_ID_DATA, comment.id);
            setData(collapse, COMMENT_ID_DATA, comment.id);

            addClass(card, CARD_CLASS);
            addClass(header, CARD_HEADER_CLASS);
            addClass(avatar, CARD_AVATAR_CLASS);
            addClass(name, CARD_NAME_CLASS);
            addClass(subtitle, CARD_SUBTITLE_CLASS);
            addClass(body, CARD_BODY_CLASS);
            addClass(options, CARD_OPTIONS_CLASS);
            addClass(reply, OPTION_BTN_CLASS);
            addClass(reply, SHOW_REPLY_JS);
            addClass(collapse, OPTION_BTN_CLASS);
            addClass(collapse, COLLAPSE_THREAD_JS);

            append(options, reply);
            append(options, collapse);
            append(header, options);
            append(header, avatar);
            append(header, name);
            append(header, subtitle);
            append(contents, body);

            if(children) {
                addClass(children, CARD_BODY_CLASS);
                append(contents, children);
            }

            append(card, header);
            append(card, contents);
            append(cards, card);
        });

        return cards;
    };

    var _redraw = function(comments) {
        if(!comments) return;

        var $coms = $(COMS_ID);
        var parentMap = {};
        var parent;

        $coms.innerHTML = "";

        comments.forEach(function(comment) {
            parent = comment.parent;
            if(!(parent in parentMap)) {
                parentMap[parent] = [];
            }
            parentMap[parent].push(comment);
        });

        var cards = _makeCards(parentMap, -1);
        if(cards) {
            append($coms, cards);
        }
    };

    var _postRoot = function() {
        var $rootComment = $(ROOT_COMMENT_ID);
        var $rootName = $(ROOT_NAME_ID);
        var rootCommentValue = $rootComment.value;
        var rootNameValue = $rootName.value;
        var data;

        removeClass($rootComment, IS_ERROR_CLASS);
        removeClass($rootName, IS_ERROR_CLASS);

        if(!rootCommentValue || !rootCommentValue.length) {
            addClass($rootComment, IS_ERROR_CLASS);
            return;
        }

        if(!rootNameValue || !rootNameValue.length) {
            addClass($rootName, IS_ERROR_CLASS);
            return;
        }

        data = {
            url: document.location,
            comment: $rootComment.value,
            name: $rootName.value,
            parent: -1
        };

        if(_honeypot){
            data.gotcha = $(HONEYPOT_ID).value;
        }

        startSpinner(ROOT_BTN_ID);
        post(_api.create, data, function() {
            _getComments(function() {
                $rootComment.value = "";
                clearSpinner(ROOT_BTN_ID);
            });
        });
    };

    var _submitReply = function(id) {
        var $replyTextArea = $(TEXTAREA_ID + id);
        var $nameInput = $(NAME_INPUT_ID + id);
        var textareaValue = $replyTextArea.value;
        var nameInputValue = $nameInput.value;

        removeClass($replyTextArea, IS_ERROR_CLASS);
        removeClass($nameInput, IS_ERROR_CLASS);

        if(!textareaValue || !textareaValue.length) {
            addClass($replyTextArea, IS_ERROR_CLASS);
            return;
        }
        if(!nameInputValue || !nameInputValue.length) {
            addClass($nameInput, IS_ERROR_CLASS);
            return;
        }

        var data = {
            comment: textareaValue,
            name: nameInputValue,
            parent: id,
            url: document.location
        };

        if(_honeypot){
            data.gotcha = $(GOTCHA_ID + id).value;
        }

        startSpinner(SUBMIT_BTN_ID + id);
        post(_api.create, data, function() {
            _getComments(function() {
                clearSpinner(ROOT_BTN_ID);
            });
        });
    };

    var _cancelReply = function(id){
        $(TEXTAREA_ID + id).remove();
        $(SUBMIT_BTN_ID + id).remove();
        $(CANCEL_BTN_ID + id).remove();
        $(NAME_INPUT_ID + id).remove();
        addClass($(REPLY_BTN_ID + id), INITIAL_CLASS);
    };

    var _showReply = function(id) {
        var $body = $(BODY_ID + id);
        var textarea = create("textarea");
        var otherFields = create("div");
        var otherFieldsContainer = create("div");
        var name = create("input");
        var honeypot = create("input");
        var cancel = create("button");
        var submit = create("button");

        textarea.id = TEXTAREA_ID + id;
        name.id = NAME_INPUT_ID + id;
        honeypot.id = GOTCHA_ID + id;
        cancel.id = CANCEL_BTN_ID + id;
        submit.id = SUBMIT_BTN_ID + id;

        cancel.innerHTML = "Cancel";
        submit.innerHTML = "Reply";

        setData(cancel, COMMENT_ID_DATA, id);
        setData(submit, COMMENT_ID_DATA, id);

        addClass($(REPLY_BTN_ID + id), HIDDEN_CLASS);
        addClass(textarea, TEXTAREA_CLASS);
        addClass(otherFieldsContainer, OTHER_FIELDS_CONTAINER_CLASS);
        addClass(otherFields, OTHER_FIELDS_CLASS);
        addClass(name, INPUT_CLASS);
        addClass(honeypot, HIDDEN_CLASS);
        addClass(cancel, CANCEL_JS);
        addClass(cancel, POST_BTN_CLASS);
        addClass(cancel, POST_RED_BTN_CLASS);
        addClass(submit, POST_BTN_CLASS);
        addClass(submit, POST_PRIMARY_BTN_CLASS);
        addClass(submit, SUBMIT_JS);

        textarea.oninput = autoExpander(textarea);

        setAttr(name, "placeholder", "Name");

        append($body, textarea);
        append(otherFields, name);
        if(_honeypot) {
            append(otherFields, honeypot);
        }
        append(otherFields, submit);
        append(otherFields, cancel);
        append(otherFieldsContainer, otherFields);

        append($body, otherFieldsContainer);
        _expandThread(id);
    };

    var _collapseThread = function(id) {
        var $contents = $(CONTENTS_ID + id);
        var $button = $(COLLAPSE_BTN_ID + id);

        addClass($contents, HIDDEN_CLASS);
        $button.innerHTML = "Expand";

        setTimeout(function() {
            addClass($button, EXPAND_THREAD_JS);
            removeClass($button, COLLAPSE_BTN_ID);
        }, 0);
    }

    var _expandThread = function(id) {
        var $contents = $(CONTENTS_ID + id);
        var $button = $(COLLAPSE_BTN_ID + id);

        removeClass($contents, HIDDEN_CLASS);
        $button.innerHTML = "Collapse";

        setTimeout(function() {
            addClass($button, COLLAPSE_BTN_ID);
            removeClass($button, EXPAND_THREAD_JS);
        }, 0);
    }

    /**
     * Namespace definition and public members  ///////////////////////////////////////////
     *
     * This section is public, anything here can be accessed by clients.
     */

    var Commento = global.Commento || {};

    Commento.version = '0.2.0';

    Commento.init = function(configuration) {
        _serverUrl = configuration.serverUrl || _serverUrl;
        _honeypot = configuration.honeypot || _honeypot;
        _showdownUrl = configuration.showdownUrl || (_serverUrl + _showdownUrl);
        _commentoCssUrl = configuration.commentoCssUrl || (_serverUrl + _commentoCssUrl);

        _api.get = _serverUrl + '/get';
        _api.create = _serverUrl + '/create';

        loadCSS(_commentoCssUrl);

        loadJS(_showdownUrl, function() {
            _showdownConverter = new showdown.Converter();

            var commento = $(COMMENTO_ID);
            var div = create("div");
            var textarea = create("textarea");
            var otherFieldsContainer = create("div");
            var otherFields = create("div");
            var name = create("input");
            var button = create("button");
            var honeypot = create("input");
            var commentEl = create("div");

            textarea.id = ROOT_COMMENT_ID;
            name.id = ROOT_NAME_ID;
            commentEl.id = COMS_ID;
            honeypot.id = HONEYPOT_ID;
            button.id = ROOT_BTN_ID;

            button.innerHTML = "Post comment";

            addClass(textarea, TEXTAREA_CLASS);
            addClass(otherFields, OTHER_FIELDS_CLASS);
            addClass(otherFieldsContainer, OTHER_FIELDS_CONTAINER_CLASS);
            addClass(name, INPUT_CLASS);
            addClass(button, POST_BTN_CLASS);
            addClass(button, POST_PRIMARY_BTN_CLASS);
            addClass(honeypot, HIDDEN_CLASS);

            setAttr(name, "placeholder", "Name");

            textarea.oninput = autoExpander(textarea);

            addEvent(button, 'click', _postRoot);
            addEvent(name, 'keypress', function(e) { if (e.keyCode === 13) _postRoot() });
            addEvent(commento, 'click', makeEvent(SHOW_REPLY_JS, COMMENT_ID_DATA, _showReply));
            addEvent(commento, 'click', makeEvent(COLLAPSE_THREAD_JS, COMMENT_ID_DATA, _collapseThread));
            addEvent(commento, 'click', makeEvent(EXPAND_THREAD_JS, COMMENT_ID_DATA, _expandThread));
            addEvent(commento, 'click', makeEvent(CANCEL_JS, COMMENT_ID_DATA, _cancelReply));
            addEvent(commento, 'click', makeEvent(SUBMIT_JS, COMMENT_ID_DATA, _submitReply));

            append(otherFields, name);
            append(otherFields, button);
            append(otherFieldsContainer, otherFields);

            if(_honeypot){
                append(otherFields, honeypot);
            }

            append(div, textarea);
            append(div, otherFieldsContainer);
            append(div, commentEl);
            append(commento, div);

            _getComments();
        });
    };

    /**
     * Publish ///////////////////////////////////
     *
     * Publish Commento to the world.
     */
    global.Commento = Commento;

}(window, document));
