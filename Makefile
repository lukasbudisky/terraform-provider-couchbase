TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=budisky.com
NAMESPACE=couchbase
NAME=couchbase
VERSION=0.0.1
BINARY="terraform-provider-${NAME}_${VERSION}"
OS_ARCH=linux_amd64
CGO_ENABLED=0

default: install

build:
	go build -v -a -o ${BINARY}

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

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	CB_ADDRESS=couchbase1.couchbase CB_PORT=8091 CB_USERNAME=Administrator CB_PASSWORD=123456 TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m 