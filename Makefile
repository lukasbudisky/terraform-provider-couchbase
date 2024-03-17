####################
# Go configuration #
####################

TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=budisky.com
NAMESPACE=couchbase
NAME=couchbase
VERSION=1.1.1
BINARY="terraform-provider-${NAME}_${VERSION}"
OS_ARCH=linux_amd64
CGO_ENABLED=0

#########################
# Linters configuration #
#########################

WORKSPACE="$(shell pwd)"
GIT_BRANCH="$(shell git rev-parse --abbrev-ref HEAD)"
LOG_LEVEL="INFO"

##################
# Build provider #
##################

default: install

build:
	GOOS=linux GOARCH=amd64 go build -v -a -o ${BINARY}

release:
	GOOS=darwin GOARCH=amd64 go build -v -a -o "./bin/${BINARY}_darwin_amd64_${VERSION}"
	GOOS=freebsd GOARCH=386 go build -v -a -o "./bin/${BINARY}_freebsd_386_${VERSION}"
	GOOS=freebsd GOARCH=amd64 go build -v -a -o "./bin/${BINARY}_freebsd_amd64_${VERSION}"
	GOOS=freebsd GOARCH=arm go build -v -a -o "./bin/${BINARY}_freebsd_arm_${VERSION}"
	GOOS=linux GOARCH=386 go build -v -a -o "./bin/${BINARY}_linux_386_${VERSION}"
	GOOS=linux GOARCH=amd64 go build -v -a -o "./bin/${BINARY}_linux_amd64_${VERSION}"
	GOOS=linux GOARCH=arm go build -v -a -o "./bin/${BINARY}_linux_arm_${VERSION}"
	GOOS=openbsd GOARCH=386 go build -v -a -o "./bin/${BINARY}_openbsd_386_${VERSION}"
	GOOS=openbsd GOARCH=amd64 go build -v -a -o "./bin/${BINARY}_openbsd_amd64_${VERSION}"
	GOOS=solaris GOARCH=amd64 go build -v -a -o "./bin/${BINARY}_solaris_amd64_${VERSION}"
	GOOS=windows GOARCH=386 go build -v -a -o "./bin/${BINARY}_windows_386_${VERSION}"
	GOOS=windows GOARCH=amd64 go build -v -a -o "./bin/${BINARY}_windows_amd64_${VERSION}"

move:
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	
install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

#################
# Test provider #
#################

test:
	go clean -testcache
	go test -count=1 $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc:
	go clean -testcache
	CB_ADDRESS=127.0.0.1 CB_CLIENT_PORT=8091 CB_NODE_PORT=11210 CB_USERNAME=Administrator CB_PASSWORD=123456 TF_ACC=1 go test -count=1 $(TEST) -v $(TESTARGS) -timeout 120m 

########################
# Local infrastructure #
########################

cbnetup:
	docker network create couchbase

cbnetdown:
	docker network rm couchbase

cbup:
	docker-compose -f terraform_example/docker-compose.yml up -d --build

cbinit:
	./terraform_example/initialization.sh http://127.0.0.1 8091

cbdown:
	docker-compose -f terraform_example/docker-compose.yml down

###########
# Linters #
###########

lint:
	docker run --rm --platform=linux/amd64 \
		-e LOG_LEVEL=${LOG_LEVEL} \
		-e VALIDATE_ALL_CODEBASE=true \
		-e RUN_LOCAL=true \
		-e DEFAULT_BRANCH=${GIT_BRANCH} \
		-e VALIDATE_GO=false \
		-v ${WORKSPACE}:/tmp/lint \
		ghcr.io/super-linter/super-linter:latest

