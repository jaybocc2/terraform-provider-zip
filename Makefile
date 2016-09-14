.PHONY: all deps fmt .git/hooks/pre-commit terraform-provider-zipfile clean test package

all: deps fmt .git/hooks/pre-commit test terraform-provider-zipfile

fmt:
	go fmt ./...

clean:
	rm -f terraform-provider-zip

dev: terraform-provider-zipfile
	cp terraform-provider-zip $$(echo $$GOPATH|sed -e's/://')/bin

install: terraform-provider-zipfile
	cp terraform-provider-zip $$(dirname $$(which terraform))

terraform-provider-zipfile: fmt test
	go build

test:
	go test -v $$(glide novendor)

.git/hooks/pre-commit:
	if [ ! -f .git/hooks/pre-commit ]; then ln -s ../../git-hooks/pre-commit .git/hooks/pre-commit; fi

deps:
	go get -u github.com/Masterminds/glide
	glide install
	glide up
