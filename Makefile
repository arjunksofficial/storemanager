install-swagger:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

generate-swagger: install-swagger
	GO111MODULE=on swagger generate spec -w ./cmd/app -o ./api/swagger.json -m