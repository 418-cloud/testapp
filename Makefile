IMG=testapp:testing
verify:
	golint ./...
	go fmt -n ./...

test: verify
	go test -v ./...

build: test
	go build -v .

docker-build: test
	ko publish -BL .

run:
	HOSTNAME=local.machine go run main.go

docker-run: docker-build