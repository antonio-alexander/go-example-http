## ----------------------------------------------------------------------
## This makefile can be used to execute common functions to interact with
## the source code, these functions ease local development and can also be
## used in CI/CD pipelines.
## ----------------------------------------------------------------------

golangcilint_version=v1.44.2
rsa_bits=4096
ssl_config_file=./config/openssl.conf

# REFERENCE: https://stackoverflow.com/questions/16931770/makefile4-missing-separator-stop
help: ## - Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

check-lint: ## - validate/install golangci-lint installation
	@which golangci-lint || (go install github.com/golangci/golangci-lint/cmd/golangci-lint@${golangcilint_version})

lint: check-lint ## - lint the source with verbose output
	@golangci-lint run --verbose

build: ## - build the source (latest)
	@docker compose build --build-arg GIT_COMMIT=`git rev-parse HEAD` \
	--build-arg GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
	@docker image prune -f

run: ## - run the service and its dependencies (docker) detached
	@docker container rm -f go-example-http
	@docker image prune -f
	@docker compose up -d

check-openssl: ## Check if openssl is installed
	@which openssl || echo "openssl not found"

gen-certs: check-openssl ## Generate public/private SSL certificates
	@openssl req -x509 -newkey rsa:${rsa_bits} -sha256 -utf8 -days 1 -nodes \
	-config ${ssl_config_file} -keyout ./certs/ssl.key -out ./certs/ssl.crt 

stop:
	@docker compose down --volumes