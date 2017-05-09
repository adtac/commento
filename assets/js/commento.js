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
        var msPerMinute = 60000;
        var msPerHour = 3600000;
        var msPerDay = 86400000;
        var msPerMonth = 2592000000;
        var msPerYear = 946080000000;

        var elapsed = current - previous;

        if (elapsed < msPerMinute) {
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
     * _CLASS for any visual class that it's defined in a CSS file (spectre or Commento)
     * _JS for any JavaScript hook used to delegate events or trigger actions.
     * _DATA for any data key we wish to identify in dataset
     */

    var ROOT_COMMENT_ID = "commento-root-comment";
    var ROOT_NAME_ID = "commento-root-name";
    var COMS_ID = "commento-coms";
    var HONEYPOT_ID = "commento-root-gotcha";
    var REPLY_ID = "commento-reply-textarea-";
    var NAME_INPUT_ID = "commento-name-input-";
    var GOTCHA_ID = "commento-gotcha-";
    var CANCEL_BTN_ID = "commento-cancel-button-";
    var SUBMIT_BTN_ID = "commento-submit-button-";
    var BODY_ID = "commento-body-";
    var REPLY_BTN_ID = "commento-reply-button-";
    var COMMENTO_ID = "commento";

    var CANCEL_JS = 'commento-js-cancel';
    var SHOW_REPLY_JS = 'commento-show-reply-js';
    var SUBMIT_JS = 'commento-submit-js';

    var CARD_CLASS = 'card';
    var CARD_HEADER_CLASS = 'card-header';
    var CARD_FOOTER_CLASS = 'card-footer';
    var CARD_SUBTITLE_CLASS = 'card-subtitle';
    var H5_CLASS = 'h5';
    var CARD_BODY_CLASS = 'card-body';
    var EXTRA_MARGIN_CLASS = 'extra-margin';
    var HIDDEN_CLASS= 'hidden';
    var BUTTON_HOLDER_CLASS = 'button-holder';
    var MARGIN_CLASS = 'margin';
    var COMMENTS_CLASS = 'comments';
    var SUBMIT_AREA_CLASS = 'submit-area';
    var ROOT_ELEMENT_CLASS = 'root-elem';
    var INITIAL_CLASS = 'initial';
    var BUTTON_CLASS = 'btn';
    var BUTTON_PRIMARY_CLASS = 'btn-primary';
    var FORM_INPUT_CLASS = 'form-input';
    var IS_ERROR_CLASS = 'is-error';

    var COMMENT_ID_DATA = 'commentId';

    /**
     * Private fields ///////////////////////////////////////////
     *
     * This section is for internal use. This hold Commento internals that shouldn't be accessible from the outside
     */

    var _showdownUrl = "/assets/vendor/showdown.min.js";
    var _spectreUrl = "/assets/vendor/spectre.min.css";
    var _commentoCssUrl = "/assets/style/commento.min.css";
    var _serverUrl = '';
    var _honeypot = false;
    var _api = {};
    var _showdownConverter;


    var _getComments = function() {
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
            var title = create("div");
            var h5 = create("h5");
            var subtitle = create("div");
            var body = create("div");
            var footer = create("div");
            var button = create("button");
            var children = _makeCards(parentMap, comment.id);

            body.id = BODY_ID + comment.id;
            button.id = REPLY_BTN_ID + comment.id;

            h5.innerHTML = comment.name;
            body.innerHTML = _showdownConverter.makeHtml(comment.comment);
            subtitle.innerHTML = timeDifference(Date.now(), Date.parse(comment.timestamp));
            button.innerHTML = "Reply";

            setData(button, COMMENT_ID_DATA, comment.id);

            addClass(card, CARD_CLASS);
            addClass(title, CARD_HEADER_CLASS);
            addClass(h5, H5_CLASS);
            addClass(subtitle, CARD_SUBTITLE_CLASS);
            addClass(subtitle, EXTRA_MARGIN_CLASS);
            addClass(body, CARD_BODY_CLASS);
            addClass(footer, CARD_FOOTER_CLASS);
            addClass(button, BUTTON_CLASS);
            addClass(button, SHOW_REPLY_JS);

            append(footer, button);
            append(title, h5);
            append(title, subtitle);
            append(card, title);
            append(card, body);
            append(card, footer);

            if(children) {
                addClass(children, CARD_BODY_CLASS);
                append(card, children);
            }
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

        post(_api.create, data, function() {
            $rootComment.value = "";
            _getComments();
        });
    };

    var _submitReply = function(id) {
        var $replyTextArea = $(REPLY_ID + id);
        var $nameInput = $(NAME_INPUT_ID + id);
        var textAreaValue = $replyTextArea.value;
        var nameInputValue = $nameInput.value;

        removeClass($replyTextArea, IS_ERROR_CLASS);
        removeClass($nameInput, IS_ERROR_CLASS);

        if(!textAreaValue || !textAreaValue.length) {
            addClass($replyTextArea, IS_ERROR_CLASS);
            return;
        }
        if(!nameInputValue || !nameInputValue.length) {
            addClass($nameInput, IS_ERROR_CLASS);
            return;
        }

        var data = {
            comment: textAreaValue,
            name: nameInputValue,
            parent: id,
            url: document.location
        };

        if(_honeypot){
            data.gotcha = $(GOTCHA_ID + id).value;
        }

        post(_api.create, data, _getComments);
    };

    var _cancelReply = function(id){
        $(REPLY_ID + id).remove();
        $(SUBMIT_BTN_ID + id).remove();
        $(CANCEL_BTN_ID + id).remove();
        $(NAME_INPUT_ID + id).remove();
        addClass($(REPLY_BTN_ID + id), INITIAL_CLASS);
    };

    var _showReply = function(id) {
        var $body = $(BODY_ID + id);
        var textArea = create("textarea");
        var name = create("input");
        var honeypot = create("input");
        var cancel = create("button");
        var submit = create("button");
        var buttonHolder = create("div");

        textArea.id = REPLY_ID + id;
        name.id = NAME_INPUT_ID + id;
        honeypot.id = GOTCHA_ID + id;
        cancel.id = CANCEL_BTN_ID + id;
        submit.id = SUBMIT_BTN_ID + id;

        cancel.innerHTML = "Cancel";
        submit.innerHTML = "Reply";

        setData(cancel, COMMENT_ID_DATA, id);
        setData(submit, COMMENT_ID_DATA, id);

        addClass($(REPLY_BTN_ID + id), HIDDEN_CLASS);
        addClass(textArea, FORM_INPUT_CLASS);
        addClass(name, MARGIN_CLASS);
        addClass(name, FORM_INPUT_CLASS);
        addClass(honeypot, HIDDEN_CLASS);
        addClass(cancel, CANCEL_JS);
        addClass(cancel, MARGIN_CLASS);
        addClass(cancel, BUTTON_CLASS);
        addClass(submit, MARGIN_CLASS);
        addClass(submit, BUTTON_CLASS);
        addClass(submit, BUTTON_PRIMARY_CLASS);
        addClass(submit, SUBMIT_JS);
        addClass(buttonHolder, BUTTON_HOLDER_CLASS);

        setAttr(name, "placeholder", "Name");

        append($body, textArea);
        append(buttonHolder, name);
        if(_honeypot) {
            append(buttonHolder, honeypot);
        }
        append(buttonHolder, cancel);
        append(buttonHolder, submit);

        append($body, buttonHolder);
    };

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
        _spectreUrl = configuration.spectreUrl || (_serverUrl + _spectreUrl);
        _commentoCssUrl = configuration.commentoCssUrl || (_serverUrl + _commentoCssUrl);

        _api.get = _serverUrl + '/get';
        _api.create = _serverUrl + '/create';

        loadCSS(_spectreUrl);
        loadCSS(_commentoCssUrl);

        loadJS(_showdownUrl, function() {
            _showdownConverter = new showdown.Converter();

            var commento = $(COMMENTO_ID);
            var div = create("div");
            var textarea = create("textarea");
            var subArea  = create("div");
            var input = create("input");
            var button = create("button");
            var honeypot = create("input");
            var commentEl = create("div");

            textarea.id = ROOT_COMMENT_ID;
            input.id = ROOT_NAME_ID;
            commentEl.id = COMS_ID;
            honeypot.id = HONEYPOT_ID;

            button.innerHTML = "Post comment";

            addClass(div, COMMENTS_CLASS);
            addClass(textarea, FORM_INPUT_CLASS);
            addClass(subArea, SUBMIT_AREA_CLASS);
            addClass(input, FORM_INPUT_CLASS);
            addClass(input, ROOT_ELEMENT_CLASS);
            addClass(button, ROOT_ELEMENT_CLASS);
            addClass(button, BUTTON_CLASS);
            addClass(button, BUTTON_PRIMARY_CLASS);
            addClass(honeypot, HIDDEN_CLASS);

            setAttr(input, "placeholder", "Name");

            addEvent(button, 'click', _postRoot);
            addEvent(commento, 'click', makeEvent(SHOW_REPLY_JS, COMMENT_ID_DATA, _showReply));
            addEvent(commento, 'click', makeEvent(CANCEL_JS, COMMENT_ID_DATA, _cancelReply));
            addEvent(commento, 'click', makeEvent(SUBMIT_JS, COMMENT_ID_DATA, _submitReply));

            append(subArea, input);
            append(subArea, button);

            if(_honeypot){
                append(subArea, honeypot);
            }

            append(div, textarea);
            append(div, subArea);
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
