// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"google.golang.org/grpc"
	pb "karma/gen/server"
	"karma/server/service"
	"log"
	"net"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"karma/server/restapi/operations"
)

//go:generate swagger generate server --target ../../server --name Karma --spec ../../swagger.yaml --principal interface{}

func startGrpcService(server *service.Service) {
	lis, err := net.Listen("tcp", ":37700")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterServerServer(s, server)

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func configureFlags(api *operations.KarmaAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.KarmaAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.MultipartformConsumer = runtime.DiscardConsumer

	api.JSONProducer = runtime.JSONProducer()

	paths := service.NewStore()
	storages := service.NewMegaStorage(paths)
	server := service.NewService(storages)

	go startGrpcService(server)

	api.GetFileHandler = operations.GetFileHandlerFunc(service.NewDownloadHandler(storages))
	api.PutFileHandler = operations.PutFileHandlerFunc(service.NewUploadHandler(storages))

	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// operations.PutFileMaxParseMemory = 32 << 20

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
