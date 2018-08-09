.PHONY: clean build all
clean:
	rm -f secondary.so \
	rm -f ../preprocessor/preprocessor

../preprocessor/preprocessor: ../preprocessor/main.go
	(cd ../preprocessor; go build .)


secondary.so: secondary ../preprocessor/preprocessor
	../preprocessor/preprocessor ./secondary/main.go > secondary/main_generated.go && \
	GOPATH=`go env GOPATH`:`pwd`/secondary go build -buildmode=plugin ./secondary



build:
	(cd ../preprocessor; go build .)

all: build