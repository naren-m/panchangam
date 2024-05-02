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

.PHONY: start
start:
	docker compose up --force-recreate --remove-orphans --detach
	@echo ""
	@echo "OpenTelemetry Demo is running."
	@echo "Go to http://localhost:8080 for the demo UI."
	@echo "Go to http://localhost:8080/jaeger/ui for the Jaeger UI."
	@echo "Go to http://localhost:8080/grafana/ for the Grafana UI."
	@echo "Go to http://localhost:8080/loadgen/ for the Load Generator UI."
	@echo "Go to https://opentelemetry.io/docs/demo/feature-flags/ to learn how to change feature flags."