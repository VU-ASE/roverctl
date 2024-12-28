# Makefile in accordance with the docs on git management (to use in combination with meta)
.PHONY: build start clean test

BUILD_DIR=bin/
BINARY_NAME=roverctl
LDFLAGS := -X 'github.com/VU-ASE/roverctl/src/views.version=$(VERSION)'

build-open-api:
	# Check if the spec/apispec.yaml file exists
	@if [ ! -f spec/apispec.yaml ]; then \
		echo "spec/apispec.yaml file not found. Download it from the roverd repository."; \
		exit 1; \
	fi
	@echo "generating openapi client"
	@openapi-generator-cli generate -i spec/apispec.yaml -g go -o src/openapi --additional-properties=withGoMod=false

build: #build-open-api
	@if [ -z ${VERSION} ]; then \
		echo "VERSION is not set. Please set it to the semantic version you want to build."; \
		exit 1; \
	fi
	@echo "building ${BINARY_NAME}"
	@cd src/ && go build -buildvcs=false -ldflags="$(LDFLAGS)" -o "../$(BUILD_DIR)${BINARY_NAME}" ${buildargs}


#
# You can specify run arguments and build arguments using runargs and buildargs, like this:
# make start runargs="-debug"
# make start runargs="-debug" buildargs="-verbose"
# make build buildargs="-verbose"
#
start: build
	@echo "starting ${BINARY_NAME}"
	./${BUILD_DIR}${BINARY_NAME} ${runargs}

clean:
	@echo "Cleaning all targets for ${BINARY_NAME}"
	rm -rf $(BUILD_DIR)
	rm -rf src/openapi
	go mod tidy

test:
	go test ./src -v -count=1 -timeout 0
