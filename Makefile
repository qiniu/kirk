SOURCEDIR = kirksdk
GLIDENOVENDOR = $(shell glide novendor)

test:
	go test $(GLIDENOVENDOR)

fmt:
	gofmt -w $(SOURCEDIR)

style:
	gofmt -l $(SOURCEDIR)
