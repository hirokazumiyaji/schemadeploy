VERSION=0.1.1

installdeps:
	go get -d
	go mod tidy

fmt:
	go fmt ./...

build: ./* cmd/schemadeploy/main.go
	go build -o bin/schemadeploy cmd/schemadeploy/main.go

build-cross: ./* cmd/schemadeploy/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/schemadeploy-${VERSION}/schemadeploy cmd/schemadeploy/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/amd64/schemadeploy-${VERSION}/schemadeploy cmd/schemadeploy/main.go

dist: build-cross
	cd bin/linux/amd64/schemadeploy-${VERSION}/ && tar zcvf schemadeploy-linux-amd64-${VERSION}.tar.gz schemadeploy
	cd bin/darwin/amd64/schemadeploy-${VERSION}/ && tar zcvf schemadeploy-darwin-amd64-${VERSION}.tar.gz schemadeploy

clean:
	rm -rf bin/*
