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
	@if validated="$$(go tool tfplugindocs validate --provider-name dql 2>&1)"; then \
		echo "tfplugindocs validate: the Registry pages match the provider"; \
	else \
		printf '%s\n' "$$validated" | awk '/\[DEBUG\]/ { next } /^Error/ { failed = 1 } \
			{ if (failed) printf "\033[31m%s\033[0m\n", $$0; else print }'; \
		printf '\033[31m%s\033[0m\n' "tfplugindocs validate: failed, the Registry pages do not match the provider"; \
		exit 1; \
	fi
