SHELL = /bin/bash

.PHONY: test
test:
	go test -race -coverprofile=coverage.txt -covermode=atomic  ./...

mssh.linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./cli/mssh.linux-amd64  cmd/ssh/main.go

mssh.darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o ./cli/mssh.darwin-amd64 cmd/ssh/main.go

mssh.darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -o ./cli/mssh.darwin-arm64 cmd/ssh/main.go

mssh.windows-amd64.exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o ./cli/mssh.windows-amd64.exe cmd/ssh/main.go



kubectl-mexec.linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./cli/kubectl-mexec.linux-amd64  cmd/k8s/main.go

kubectl-mexec.darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o ./cli/kubectl-mexec.darwin-amd64 cmd/k8s/main.go

kubectl-mexec.darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -o ./cli/kubectl-mexec.darwin-arm64 cmd/k8s/main.go

kubectl-mexec.windows-amd64.exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o ./cli/kubectl-mexec.windows-amd64.exe cmd/k8s/main.go



docker-mexec.linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./cli/docker-mexec.linux-amd64  cmd/docker/main.go

docker-mexec.darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -o ./cli/docker-mexec.darwin-amd64 cmd/docker/main.go

docker-mexec.darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -o ./cli/docker-mexec.darwin-arm64 cmd/docker/main.go

docker-mexec.windows-amd64.exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o ./cli/docker-mexec.windows-amd64.exe cmd/docker/main.go



.PHONY: build-mssh
build-mssh: clean mssh.linux-amd64 mssh.darwin-amd64 mssh.darwin-arm64 mssh.windows-amd64.exe
	cd ./cli && find . -name 'mssh*' | xargs -I{} tar czf {}.tar.gz {}
	cd ./cli && shasum -a 256 mssh* > mssh_sha256sum.txt
	cat ./cli/mssh_sha256sum.txt

.PHONY: build-k8s-mexec
build-k8s-mexec: clean kubectl-mexec.linux-amd64 kubectl-mexec.darwin-amd64 kubectl-mexec.darwin-arm64 kubectl-mexec.windows-amd64.exe
	cd ./cli && find . -name 'kubectl-mexec*' | xargs -I{} tar czf {}.tar.gz {}
	cd ./cli && shasum -a 256 kubectl-mexec* > k8s-mexec_sha256sum.txt
	cat ./cli/k8s-mexec_sha256sum.txt

.PHONY: build-docker-mexec
build-docker-mexec: clean docker-mexec.linux-amd64 docker-mexec.darwin-amd64 docker-mexec.darwin-arm64 docker-mexec.windows-amd64.exe
	cd ./cli && find . -name 'docker-mexec*' | xargs -I{} tar czf {}.tar.gz {}
	cd ./cli && shasum -a 256 docker-mexec* > docker-mexec_sha256sum.txt
	cat ./cli/docker-mexec_sha256sum.txt


.PHONY: build
build: build-mssh build-k8s-mexec build-docker-mexec

.PHONY: release
release:
	git tag v$(V)
	git push origin v$(V)

.PHONY: clean
clean:
	-rm -r ./cli/

