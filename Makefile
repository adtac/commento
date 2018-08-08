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
TEMPLATES_DEVEL_ROOT_DIR  = $(DEVEL_BUILD_DIR)
TEMPLATES_PROD_BUILD_DIR  = $(TEMPLATES_BUILD_DIR)/$(PROD_BUILD_DIR)
TEMPLATES_PROD_ROOT_DIR   = $(PROD_BUILD_DIR)

DB_BUILD_DIR              = db
DB_DEVEL_BUILD_DIR        = $(DB_BUILD_DIR)/$(DEVEL_BUILD_DIR)
DB_DEVEL_ROOT_DIR         = $(DEVEL_BUILD_DIR)/db
DB_PROD_BUILD_DIR         = $(DB_BUILD_DIR)/$(PROD_BUILD_DIR)
DB_PROD_ROOT_DIR          = $(PROD_BUILD_DIR)/db

devel: devel-frontend devel-api devel-templates devel-db

prod: prod-frontend prod-api prod-templates prod-db

test: api

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
	cd templates && $(MAKE) $(MAKECMDGOALS)

.PHONY: db
db:
	cd db && $(MAKE) $(MAKECMDGOALS)

devel-frontend: frontend
	cp -r $(FRONTEND_DEVEL_BUILD_DIR)/* $(FRONTEND_DEVEL_ROOT_DIR)

devel-api: api
	cp -r $(API_DEVEL_BUILD_DIR)/* $(API_DEVEL_ROOT_DIR)

devel-templates: templates
	cp -r $(TEMPLATES_DEVEL_BUILD_DIR)/* $(TEMPLATES_DEVEL_ROOT_DIR)

devel-db: db
	cp -r $(DB_DEVEL_BUILD_DIR)/* $(DB_DEVEL_ROOT_DIR)

prod-frontend: frontend
	cp -r $(FRONTEND_PROD_BUILD_DIR)/* $(FRONTEND_PROD_ROOT_DIR)

prod-api: api
	cp -r $(API_PROD_BUILD_DIR)/* $(API_PROD_ROOT_DIR)

prod-templates: templates
	cp -r $(TEMPLATES_PROD_BUILD_DIR)/* $(TEMPLATES_PROD_ROOT_DIR)

prod-db: db
	cp -r $(DB_PROD_BUILD_DIR)/* $(DB_PROD_ROOT_DIR)

clean: clean-root clean-frontend clean-api clean-templates clean-db

clean-root:
	rm -rf build

clean-frontend:
	cd frontend && $(MAKE) $(MAKECMDGOALS)

clean-api:
	cd api && $(MAKE) $(MAKECMDGOALS)

clean-templates:
	cd templates && $(MAKE) $(MAKECMDGOALS)

clean-db:
	cd db && $(MAKE) $(MAKECMDGOALS)

$(shell mkdir -p $(FRONTEND_DEVEL_ROOT_DIR) $(API_DEVEL_ROOT_DIR) $(TEMPLATES_DEVEL_ROOT_DIR) $(FRONTEND_PROD_ROOT_DIR) $(API_PROD_ROOT_DIR) $(TEMPLATES_DEVEL_ROOT_DIR))
