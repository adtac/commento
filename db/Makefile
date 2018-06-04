SHELL = bash

BUILD_DIR                   = build
DEVEL_BUILD_DIR             = $(BUILD_DIR)/devel
PROD_BUILD_DIR              = $(BUILD_DIR)/prod

DB_SRC_DIR                  = .
DB_SRC_FILES                = $(wildcard $(DB_SRC_DIR)/*.sql)
DB_DEVEL_BUILD_DIR          = $(DEVEL_BUILD_DIR)/db
DB_DEVEL_BUILD_FILES        = $(patsubst $(DB_SRC_DIR)/%, $(DB_DEVEL_BUILD_DIR)/%, $(DB_SRC_FILES))
DB_PROD_BUILD_DIR           = $(PROD_BUILD_DIR)/db
DB_PROD_BUILD_FILES         = $(patsubst $(DB_SRC_DIR)/%, $(DB_PROD_BUILD_DIR)/%, $(DB_SRC_FILES))

devel: devel-db

prod: prod-db

clean:
	rm -rf $(BUILD_DIR)

devel-db: $(DB_DEVEL_BUILD_FILES)

$(DB_DEVEL_BUILD_FILES): $(DB_DEVEL_BUILD_DIR)/%.sql: $(DB_SRC_DIR)/%.sql
	cp $^ $@;

prod-db: $(DB_PROD_BUILD_FILES)

$(DB_PROD_BUILD_FILES): $(DB_PROD_BUILD_DIR)/%.sql: $(DB_SRC_DIR)/%.sql
	cp $^ $@;

$(shell mkdir -p $(DB_DEVEL_BUILD_DIR) $(DB_PROD_BUILD_DIR))
