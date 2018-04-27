SHELL = bash

BUILD_DIR              = build
DEVEL_BUILD_DIR        = $(BUILD_DIR)/devel
PROD_BUILD_DIR         = $(BUILD_DIR)/prod

TMPL_SRC_DIR           = .
TMPL_SRC_FILES         = $(wildcard $(TMPL_SRC_DIR)/*.tmpl)
TMPL_DEVEL_BUILD_DIR   = $(DEVEL_BUILD_DIR)
TMPL_DEVEL_BUILD_FILES = $(patsubst $(TMPL_SRC_DIR)/%, $(TMPL_DEVEL_BUILD_DIR)/%, $(TMPL_SRC_FILES))
TMPL_PROD_BUILD_DIR    = $(PROD_BUILD_DIR)
TMPL_PROD_BUILD_FILES  = $(patsubst $(TMPL_SRC_DIR)/%, $(TMPL_PROD_BUILD_DIR)/%, $(TMPL_SRC_FILES))

devel: devel-templates

prod: prod-templates

clean:
	rm -rf $(BUILD_DIR)

devel-templates: $(TMPL_DEVEL_BUILD_FILES)

$(TMPL_DEVEL_BUILD_FILES): $(TMPL_DEVEL_BUILD_DIR)/%.tmpl: $(TMPL_SRC_DIR)/%.tmpl
	cp $^ $@;

prod-templates: $(TMPL_PROD_BUILD_FILES)

$(TMPL_PROD_BUILD_FILES): $(TMPL_PROD_BUILD_DIR)/%.tmpl: $(TMPL_SRC_DIR)/%.tmpl
	cp $^ $@;

$(shell mkdir -p $(TMPL_DEVEL_BUILD_DIR) $(TMPL_PROD_BUILD_DIR))
