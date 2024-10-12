help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## Run the CLI version of the tool
cli/run:
	go run ./cmd/relay365cli/relay365cli.go

## Run the Service version of the tool
svc/run:
	go run ./cmd/relay365svc/relay365svc.go

