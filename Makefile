# Docs
docs_fmt:
	terraform fmt -recursive ./examples/

docs: docs_fmt
	go get github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@v0.19.4
	go mod tidy
	make docs_fmt
	tfplugindocs generate --tf-version=1.9.7 --provider-name=rmon

build:
	go build -o bin/terraform-provider-rmon