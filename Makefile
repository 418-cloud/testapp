IMG=testapp:testing
verify:
	go get golang.org/x/lint/golint
	golint ./...
	go fmt -n ./...

test: verify
	go test -v ./...

build: test
	go build -v .

docker-build: test
	ko publish -BL .

run:
	HOST=local.machine go run main.go -templateRoot=kodata

docker-run: docker-build