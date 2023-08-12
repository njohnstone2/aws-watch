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
	$(eval BUILD_IMG=$(shell $(WITH_GOFLAGS) KO_DOCKER_REPO="$(KO_DOCKER_REPO)" KO_DEFAULTBASEIMAGE=$(KO_DEFAULTBASEIMAGE) ko build -B github.com/njohnstone2/aws-watch/cmd/aws-watch -t $(KO_TAGS)))
	$(eval IMG_REPOSITORY=$(shell echo $(BUILD_IMG) | cut -d "@" -f 1 | cut -d ":" -f 1))
	$(eval IMG_TAG=$(shell echo $(BUILD_IMG) | cut -d "@" -f 1 | cut -d ":" -f 2 -s))
	$(eval IMG_DIGEST=$(shell echo $(BUILD_IMG) | cut -d "@" -f 2))

.PHONY: build test vulncheck
