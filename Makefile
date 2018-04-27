SHELL = bash

BUILD_DIR                 = build
DEVEL_BUILD_DIR           = $(BUILD_DIR)/devel
PROD_BUILD_DIR            = $(BUILD_DIR)/prod

FRONTEND_BUILD_DIR        = frontend
FRONTEND_DEVEL_BUILD_DIR  = $(FRONTEND_BUILD_DIR)/$(DEVEL_BUILD_DIR)
FRONTEND_DEVEL_ROOT_DIR   = $(DEVEL_BUILD_DIR)
FRONTEND_PROD_BUILD_DIR   = $(FRONTEND_BUILD_DIR)/$(PROD_BUILD_DIR)
FRONTEND_PROD_ROOT_DIR    = $(PROD_BUILD_DIR)

API_BUILD_DIR             = api
API_DEVEL_BUILD_DIR       = $(API_BUILD_DIR)/$(DEVEL_BUILD_DIR)
API_DEVEL_ROOT_DIR        = $(DEVEL_BUILD_DIR)
API_PROD_BUILD_DIR        = $(API_BUILD_DIR)/$(PROD_BUILD_DIR)
API_PROD_ROOT_DIR         = $(PROD_BUILD_DIR)

TEMPLATES_BUILD_DIR       = templates
TEMPLATES_DEVEL_BUILD_DIR = $(TEMPLATES_BUILD_DIR)/$(DEVEL_BUILD_DIR)
TEMPLATES_DEVEL_ROOT_DIR  = $(DEVEL_BUILD_DIR)/templates
TEMPLATES_PROD_BUILD_DIR  = $(TEMPLATES_BUILD_DIR)/$(PROD_BUILD_DIR)
TEMPLATES_PROD_ROOT_DIR   = $(PROD_BUILD_DIR)/templates

devel: frontend api templates devel-copy

prod: frontend api templates prod-copy

# TODO: This can probably be written better: instead of explicitly defining
# each target subdirectory, define them at the top and automatically do stuff.

.PHONY: frontend
frontend:
	cd frontend && $(MAKE) $(MAKECMDGOALS)

.PHONY: api
api:
	cd api && $(MAKE) $(MAKECMDGOALS)

.PHONY: templates
templates:
	cd templates && $(MAKE) $(MAKECDMGOALS)

devel-copy: devel-copy-frontend devel-copy-api devel-copy-templates

prod-copy: prod-copy-frontend prod-copy-api prod-copy-templates

devel-copy-frontend:
	cp -r $(FRONTEND_DEVEL_BUILD_DIR)/* $(FRONTEND_DEVEL_ROOT_DIR)

devel-copy-api:
	cp -r $(API_DEVEL_BUILD_DIR)/* $(API_DEVEL_ROOT_DIR)

devel-copy-templates:
	cp -r $(TEMPLATES_DEVEL_BUILD_DIR)/* $(TEMPLATES_DEVEL_ROOT_DIR)

prod-copy-frontend:
	cp -r $(FRONTEND_PROD_BUILD_DIR)/* $(FRONTEND_PROD_ROOT_DIR)

prod-copy-api:
	cp -r $(API_PROD_BUILD_DIR)/* $(API_PROD_ROOT_DIR)

prod-copy-templates:
	cp -r $(TEMPLATES_PROD_BUILD_DIR)/* $(TEMPLATES_PROD_ROOT_DIR)

clean: clean-root clean-frontend clean-api clean-templates

clean-root:
	rm -rf build

clean-frontend:
	cd frontend && $(MAKE) $(MAKECMDGOALS)

clean-api:
	cd api && $(MAKE) $(MAKECMDGOALS)

clean-templates:
	cd templates && $(MAKE) $(MAKECMDGOALS)

$(shell mkdir -p $(FRONTEND_DEVEL_ROOT_DIR) $(API_DEVEL_ROOT_DIR) $(TEMPLATES_DEVEL_ROOT_DIR) $(FRONTEND_PROD_ROOT_DIR) $(API_PROD_ROOT_DIR) $(TEMPLATES_DEVEL_ROOT_DIR))
