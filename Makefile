.PHONY: build clean test schema

BIN_DIR := bin
BINARIES := basanos assert_equals assert_contains assert_matches assert_gt assert_gte assert_lt assert_lte

build: $(addprefix $(BIN_DIR)/,$(BINARIES))

$(BIN_DIR)/basanos: main.go $(shell find internal -name '*.go')
	@mkdir -p $(BIN_DIR)
	go build -o $@ .

$(BIN_DIR)/assert_%: cmd/assert_%/main.go $(shell find internal -name '*.go')
	@mkdir -p $(BIN_DIR)
	go build -o $@ ./cmd/assert_$*

clean:
	rm -rf $(BIN_DIR)

test:
	go test ./... -count=1

schema:
	@mkdir -p schema
	go run ./cmd/gen-schema > schema/events.json
