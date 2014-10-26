check:
	golint *go
	go fmt
	go vet
	errcheck bitbucket.org/utils/config
