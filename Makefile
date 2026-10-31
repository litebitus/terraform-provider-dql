.PHONY: build docs install snapshot test

build: test
	go build -o bin/terraform-provider-dql .

docs:
	go tool tfplugindocs generate --provider-name dql

install: build
	@printf '%s\n' \
		'Add to ~/.terraformrc for plans to use this build:' \
		'provider_installation {' \
		'  dev_overrides {' \
		'    "litebitus/dql" = "$(CURDIR)/bin"' \
		'  }' \
		'  direct {}' \
		'}'

snapshot: test
	goreleaser release --snapshot --clean --skip=sign

test:
	go vet ./...
	@unformatted="$$(gofmt -l .)"; [ -z "$$unformatted" ] || { echo "not gofmt-formatted:"; echo "$$unformatted"; exit 1; }
	go test ./...
	go tool tfplugindocs validate --provider-name dql
