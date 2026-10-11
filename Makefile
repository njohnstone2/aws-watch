GOFLAGS ?= -ldflags=-X=github.com/njohnstone2/aws-watch/pkg/utils/project.Version=$(shell git describe --tags --always)
WITH_GOFLAGS = GOFLAGS="$(GOFLAGS)"

# Publish repo
AWS_REGION = "us-east-1"
AWS_ACCOUNT_ID ?= $(shell aws sts get-caller-identity --query Account --output text)
KO_DOCKER_REPO ?= ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/njohnstone2
KO_DEFAULTBASEIMAGE = "public.ecr.aws/lambda/go:1"
KO_TAGS ?= "latest"

TEST_SUITE ?= "..."

test:
	go test -v ./$(TEST_SUITE)

vulncheck:
	@govulncheck ./...

build:
	$(WITH_GOFLAGS) \
		KO_DEFAULTBASEIMAGE=$(KO_DEFAULTBASEIMAGE) \
		ko build --local -B github.com/njohnstone2/aws-watch/cmd/aws-watch

# Build and publish the container image.
publish:
	$(WITH_GOFLAGS) \
		KO_DOCKER_REPO="$(KO_DOCKER_REPO)" \
		KO_DEFAULTBASEIMAGE=$(KO_DEFAULTBASEIMAGE) \
		ko build -B github.com/njohnstone2/aws-watch/cmd/aws-watch -t $(KO_TAGS)

.PHONY: build test vulncheck
