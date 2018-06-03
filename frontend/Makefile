SHELL = bash

# list of JS files to be built
JS_BUILD = jquery.js vue.js highlight.js chartist.js login.js signup.js dashboard.js

jquery.js = jquery.js
vue.js = vue.js
highlight.js = highlight.js
chartist.js = chartist.js
login.js = utils.js http.js auth-common.js login.js
signup.js = utils.js http.js auth-common.js signup.js
dashboard.js = utils.js http.js errors.js self.js dashboard.js dashboard-setting.js dashboard-domain.js dashboard-installation.js dashboard-general.js dashboard-moderation.js dashboard-statistics.js dashboard-import.js dashboard-danger.js

# for each file in $(JS_BUILD), list its composition

BUILD_DIR              = build
DEVEL_BUILD_DIR        = $(BUILD_DIR)/devel
PROD_BUILD_DIR         = $(BUILD_DIR)/prod

HTML_SRC_DIR           = .
HTML_SRC_FILES         = $(wildcard $(HTML_SRC_DIR)/*.html)
HTML_DEVEL_BUILD_DIR   = $(DEVEL_BUILD_DIR)
HTML_DEVEL_BUILD_FILES = $(patsubst $(HTML_SRC_DIR)/%, $(HTML_DEVEL_BUILD_DIR)/%, $(HTML_SRC_FILES))
HTML_PROD_BUILD_DIR    = $(PROD_BUILD_DIR)
HTML_PROD_BUILD_FILES  = $(patsubst $(HTML_SRC_DIR)/%, $(HTML_PROD_BUILD_DIR)/%, $(HTML_SRC_FILES))

HTML_MINIFIER          = html-minifier
HTML_MINIFIER_FLAGS    = --collapse-whitespace --remove-comments

JS_SRC_DIR             = js
JS_SRC_FILES           = $(wildcard $(JS_SRC_DIR)/*.js)
JS_DEVEL_BUILD_DIR     = $(DEVEL_BUILD_DIR)/js
JS_DEVEL_BUILD_FILES   = $(addprefix $(JS_DEVEL_BUILD_DIR)/, $(JS_BUILD))
JS_PROD_BUILD_DIR      = $(PROD_BUILD_DIR)/js
JS_PROD_BUILD_FILES    = $(addprefix $(JS_PROD_BUILD_DIR)/, $(JS_BUILD))

JS_MINIFIER            = uglifyjs
JS_MINIFIER_FLAGS      = --compress --mangle

SASS_SRC_DIR           = sass
SASS_SRC_FILES         = $(wildcard $(SASS_SRC_DIR)/*.scss)
CSS_DEVEL_BUILD_DIR    = $(DEVEL_BUILD_DIR)/css
CSS_DEVEL_BUILD_FILES  = $(patsubst $(SASS_SRC_DIR)/%.scss, $(CSS_DEVEL_BUILD_DIR)/%.css, $(SASS_SRC_FILES))
CSS_PROD_BUILD_DIR     = $(PROD_BUILD_DIR)/css
CSS_PROD_BUILD_FILES   = $(patsubst $(SASS_SRC_DIR)/%.scss, $(CSS_PROD_BUILD_DIR)/%.css, $(SASS_SRC_FILES))

CSS                    = sass
CSS_DEVEL_FLAGS        = --sourcemap=none
CSS_PROD_FLAGS         = $(CSS_DEVEL_FLAGS) --style compressed

IMGS_SRC_DIR           = images
IMGS_SRC_FILES         = $(wildcard $(IMGS_SRC_DIR)/*)
IMGS_DEVEL_BUILD_DIR   = $(DEVEL_BUILD_DIR)/images
IMGS_DEVEL_BUILD_FILES = $(patsubst $(IMGS_SRC_DIR)/%, $(IMGS_DEVEL_BUILD_DIR)/%, $(IMGS_SRC_FILES))
IMGS_PROD_BUILD_DIR    = $(PROD_BUILD_DIR)/images
IMGS_PROD_BUILD_FILES  = $(patsubst $(IMGS_SRC_DIR)/%, $(IMGS_PROD_BUILD_DIR)/%, $(IMGS_SRC_FILES))

devel: devel-html devel-js devel-css devel-imgs

prod: devel prod-html prod-js prod-css prod-imgs

clean:
	-rm -rf $(BUILD_DIR);

devel-html: $(HTML_DEVEL_BUILD_FILES)

$(HTML_DEVEL_BUILD_FILES): $(HTML_DEVEL_BUILD_DIR)/%.html: $(HTML_SRC_DIR)/%.html
	cp $^ $@;

prod-html: $(HTML_PROD_BUILD_FILES)

$(HTML_PROD_BUILD_FILES): $(HTML_PROD_BUILD_DIR)/%.html: $(HTML_DEVEL_BUILD_DIR)/%.html
	$(HTML_MINIFIER) $(HTML_MINIFIER_FLAGS) -o $@ $^;

devel-js: $(JS_DEVEL_BUILD_FILES)

.SECONDEXPANSION:
$(JS_DEVEL_BUILD_FILES): $(JS_DEVEL_BUILD_DIR)/%.js: $$(addprefix $$(JS_SRC_DIR)/, $$(%.js))
	>$@; \
	for f in $^; do \
		printf "// %s\n" "$$f" >>$@; \
		cat $$f >>$@; \
		printf "\n" >>$@; \
	done;

prod-js: $(JS_PROD_BUILD_FILES)

$(JS_PROD_BUILD_FILES): $(JS_PROD_BUILD_DIR)/%.js: $(JS_DEVEL_BUILD_DIR)/%.js
	$(JS_MINIFIER) $(JS_MINIFIER_FLAGS) -o $@ $^;

devel-css: $(CSS_DEVEL_BUILD_FILES)

$(CSS_DEVEL_BUILD_FILES): $(CSS_DEVEL_BUILD_DIR)/%.css: $(SASS_SRC_DIR)/%.scss $(SASS_SRC_FILES)
	$(CSS) $(CSS_DEVEL_FLAGS) $< $@;

prod-css: $(CSS_PROD_BUILD_FILES)

$(CSS_PROD_BUILD_FILES): $(CSS_PROD_BUILD_DIR)/%.css: $(SASS_SRC_DIR)/%.scss
	$(CSS) $(CSS_PROD_FLAGS) $^ $@;

$(shell mkdir -p $(HTML_DEVEL_BUILD_DIR) $(JS_DEVEL_BUILD_DIR) $(CSS_DEVEL_BUILD_DIR) $(HTML_PROD_BUILD_DIR) $(JS_PROD_BUILD_DIR) $(CSS_PROD_BUILD_DIR))

devel-imgs: $(IMGS_DEVEL_BUILD_FILES)

$(IMGS_DEVEL_BUILD_FILES): $(IMGS_DEVEL_BUILD_DIR)/%: $(IMGS_SRC_DIR)/%
	cp $^ $@;

prod-imgs: $(IMGS_PROD_BUILD_FILES)

$(IMGS_PROD_BUILD_FILES): $(IMGS_PROD_BUILD_DIR)/%: $(IMGS_SRC_DIR)/%
	cp $^ $@

$(shell mkdir -p $(HTML_DEVEL_BUILD_DIR) $(JS_DEVEL_BUILD_DIR) $(CSS_DEVEL_BUILD_DIR) $(IMGS_DEVEL_BUILD_DIR) $(HTML_PROD_BUILD_DIR) $(JS_PROD_BUILD_DIR) $(CSS_PROD_BUILD_DIR) $(IMGS_PROD_BUILD_DIR))
