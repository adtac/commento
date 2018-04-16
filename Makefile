SHELL = bash

BUILD_DIR                = build
DEVEL_BUILD_DIR          = $(BUILD_DIR)/devel
PROD_BUILD_DIR           = $(BUILD_DIR)/prod

FRONTEND_BUILD_DIR       = frontend
FRONTEND_DEVEL_BUILD_DIR = $(FRONTEND_BUILD_DIR)/$(DEVEL_BUILD_DIR)
FRONTEND_PROD_BUILD_DIR  = $(FRONTEND_BUILD_DIR)/$(PROD_BUILD_DIR)

API_BUILD_DIR            = api
API_DEVEL_BUILD_DIR      = $(API_BUILD_DIR)/$(DEVEL_BUILD_DIR)
API_PROD_BUILD_DIR       = $(API_BUILD_DIR)/$(PROD_BUILD_DIR)

devel: frontend api devel-copy

prod: frontend api prod-copy

# TODO: This can probably be written better: instead of explicitly defining
# each target subdirectory, define them at the top and automatically do stuff.

.PHONY: frontend
frontend:
	cd frontend && $(MAKE) $(MAKECMDGOALS)

.PHONY: api
api:
	cd api && $(MAKE) $(MAKECMDGOALS)

devel-copy: devel-copy-frontend devel-copy-api

prod-copy: prod-copy-frontend prod-copy-api

devel-copy-frontend:
	cp -r $(FRONTEND_DEVEL_BUILD_DIR)/* $(DEVEL_BUILD_DIR)

devel-copy-api:
	cp -r $(API_DEVEL_BUILD_DIR)/* $(DEVEL_BUILD_DIR)

prod-copy-frontend:
	cp -r $(FRONTEND_PROD_BUILD_DIR)/* $(PROD_BUILD_DIR)

prod-copy-api:
	cp -r $(API_PROD_BUILD_DIR)/* $(PROD_BUILD_DIR)

clean: clean-root clean-frontend clean-api

clean-root:
	rm -rf build

clean-frontend:
	cd frontend && $(MAKE) $(MAKECMDGOALS)

clean-api:
	cd api && $(MAKE) $(MAKECMDGOALS)

$(shell mkdir -p $(DEVEL_BUILD_DIR) $(PROD_BUILD_DIR))
