default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: lint
lint:
	golangci-lint run

.PHONY: docs
docs:
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	tfplugindocs generate
	@echo "Use this site to preview markdown rendering: https://registry.terraform.io/tools/doc-preview"
