.PHONY: all
all: generate manifests

.PHONY: manifests
manifests: ## Generate CRD manifests
	controller-gen crd paths="./api/..." output:crd:artifacts:config=CRD

.PHONY: generate
generate: ## Generate code (DeepCopy, etc.)
	controller-gen object paths="./api/..."

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
