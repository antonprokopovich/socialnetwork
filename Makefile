# Include .env vars
include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N]' && read ans && [ $${ans:-N} = y ]


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/web: run the cmd/web application
.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${DB_URL}

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${DB_URL}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${DB_URL} up

.PHONY: generate
generate:
	mockgen -package=mock \
    		-source=internal/data/models.go \
    		-destination=internal/data/mock/models.go

.PHONY: lint
lint: ## Run lint
	golangci-lint run --allow-parallel-runners --out-format code-climate > ${QUALITY_REPORT} || EXIT_CODE=$$?
	@echo 'Code quality report:'
	@jq -r '.[] | "\(.location.path):\(.location.lines.begin) \(.description)"' ${QUALITY_REPORT}

.PHONY: test
test: ## Run tests
	go test -v -race -count=1 -coverpkg=${TARGET_PATH} -coverprofile=${COVERAGE_TXT}.tmp 2>&1 ${TARGET_PATH} | go-junit-report -iocopy -out ${REPORT_XML} -set-exit-code
	@echo "mode: atomic" > ${COVERAGE_TXT}
	@sed -e '/\/mock\//d' -e '/\.pb\.go:/d' -e '/main.go:/d' ${COVERAGE_TXT}.tmp \
	| awk '$$1!~/^mode:/{a[$$1]=$$2;b[$$1]+=$$3}END{for (k in a) {print k, a[k], b[k]}}' | sort >> ${COVERAGE_TXT}
	@go tool cover -html=${COVERAGE_TXT} -o ${COVERAGE_HTML}
	@gocover-cobertura -ignore-gen-files < ${COVERAGE_TXT} > ${COVERAGE_XML}
	@awk -F'"' '/<line number="/{if ($$4 > 0) yes=yes+1; else no=no+1} END{printf("total coverage: %.2f%% of statements\n", yes*100/(yes+no))}' ${COVERAGE_XML}

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...' go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor