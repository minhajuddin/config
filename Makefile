check:
	go get -t
	golint *go
	go fmt
	go vet
	errcheck bitbucket.org/utils/config
	go test
