# karma
Job test case

### Build
1. Install Go, protoc, plugins: https://grpc.io/docs/languages/go/quickstart/
2. Install go-swagger: https://github.com/go-swagger/go-swagger
3. Clone repo: `git clone https://github.com/sevings/karma.git`
4. Generate Server API: `protoc --go_out=gen/server --go_opt=paths=source_relative --go-grpc_out=gen/server --go-grpc_opt=paths=source_relative server.proto`
5. Generate Storage API: `protoc --go_out=gen/storage --go_opt=paths=source_relative --go-grpc_out=gen/storage --go-grpc_opt=paths=source_relative storage.proto`
6. Generate HTTP API: `cd server && swagger generate server -qf ../swagger.yaml`
7. Generate HTTP SDK: `cd ../client && swagger generate client -qf ../swagger.yaml`

### Run
1. Server: `go run ./server/cmd/karma-server/ --port 8000`
2. Storages: `go run ./storage/cmd/ --port 37000`
3. `go run ./storage/cmd/ --port 37001`
4. `go run ./storage/cmd/ --port 37002`
5. `go run ./storage/cmd/ --port 37003`
6. `go run ./storage/cmd/ --port 37004`
7. ...
8. Test client: `go run ./client/cmd/`

### Ways to improve
- Use the real file system
- Use a persistent DB
- Send file content in streams
- Access storages in parallel
- Allow concurrent requests
- Use buffer pools
- Send error messages
- Add unit tests
- Write comments
- ...
