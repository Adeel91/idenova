
GO_SERVICES := services/document services/identity services/verification

.PHONY: fmt vet lint test check

fmt:
	@for service in $(GO_SERVICES); do \
		(cd $$service && go fmt ./...) || exit 1; \
	done

lint:
	@for service in $(GO_SERVICES); do \
		(cd $$service && golangci-lint run --config ../../.golangci.yml ./...) || exit 1; \
	done

vet:
	@for service in $(GO_SERVICES); do \
		(cd $$service && go vet ./...) || exit 1; \
	done

test:
	@for service in $(GO_SERVICES); do \
		(cd $$service && go test ./...) || exit 1; \
	done

check: vet lint test
