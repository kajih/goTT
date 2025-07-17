OPENAPI_FILE := api/openapi_mcu.yaml
OUTPUT_FILE := web/mcu/openapi.gen.go
PACKAGE_NAME := mcu

.PHONY: generate clean

generate:
	@echo "Generating Go types from $(OPENAPI_FILE)..."
	oapi-codegen \
	  -generate types,spec \
	  -package $(PACKAGE_NAME) \
	  -o $(OUTPUT_FILE) \
	  $(OPENAPI_FILE)

clean:
	@echo "Cleaning generated files..."
	rm -f $(OUTPUT_FILE)