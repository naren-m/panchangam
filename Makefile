# Define variables
PROTO_DIR := proto
OUT_DIR := $(PROTO_DIR)/panchangam

# Define targets
.PHONY: all clean

# Default target
all: gen

# Rule to generate protobuf files
gen: $(PROTO_FILES)
	protoc --go_out=$(PROTO_DIR) --go-grpc_out=$(PROTO_DIR) $(PROTO_DIR)/panchangam.proto

# Clean target
clean:
	rm -rf $(OUT_DIR)

