.PHONY: clean
clean:
	rm -f secondary.so

secondary.so: secondary
	GOPATH=`pwd`/secondary go build -buildmode=plugin ./secondary