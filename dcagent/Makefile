all: main

main:
	@mkdir -p dist/bin dist/conf dist/log
	@echo "Release version"
	@${GOROOT}/bin/go build -v
	@mv dcagent dist/bin/dcagent
	@cp conf/dcagent.conf dist/conf/
	@echo "Build done: binary in dist dir"
debug:
	@mkdir -p dist/bin dist/conf dist/log
	@echo "Debug version"
	@${GOROOT}/bin/go build -o dist/bin/dcagent -ldflags '-s -w' main.go
	@cp conf/dcagent.conf dist/conf/
	@echo "Build done: binary in dist dir"

#test:
#	@sh -c "'$(CURDIR)/scripts/test.sh'"
#cover:
#	@sh -c "'$(CURDIR)/scripts/test.sh' cover"


clean:
	@rm -rf dist

.PHONY: all main clean debug
