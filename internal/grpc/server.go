package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	_ "github.com/ponyjackal/go-microservice-boilerplate/docs"
	"github.com/ponyjackal/go-microservice-boilerplate/internal/domain/services"
	"github.com/ponyjackal/go-microservice-boilerplate/pkg/logger"

	// protobuf
	ServiceServer "github.com/ponyjackal/go-microservice-boilerplate/proto/service"
	pbTag "github.com/ponyjackal/go-microservice-boilerplate/proto/tag"

	"github.com/bufbuild/protovalidate-go"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/emptypb"
)

// server is used to implement service.ServiceServer.
type server struct {
	ServiceServer.UnimplementedServiceServer
	tagService *services.TagService
	validator  *protovalidate.Validator
}

// GetTags implements service.ServiceServer
func (s *server) GetTags(ctx context.Context, query *pbTag.GetTagsQuery) (*pbTag.GetTagsResponse, error) {
	response, err := s.tagService.GetTags(query.Name)
	if err != nil {
		logger.Errorf("failed to get tags: %s", err)
		return nil, err
	}
	return response, nil
}

// GetTagById implements service.ServiceServer
func (s *server) GetTagById(ctx context.Context, request *pbTag.TagId) (*pbTag.Tag, error) {
	tag, err := s.tagService.GetTagById(request)
	if err != nil {
		logger.Errorf("failed to get a tag by id: %s", err)
		return nil, err
	}

	return tag, err
}

// SaveTag implements service.ServiceServer
func (s *server) SaveTag(ctx context.Context, request *pbTag.SaveTagRequest) (*pbTag.Tag, error) {
	response, err := s.tagService.SaveTag(request)
	if err != nil {
		logger.Errorf("failed to save tag: %s", err)
		return nil, err
	}

	return response, nil
}

// UpdateTag implements service.ServiceServer
func (s *server) UpdateTag(ctx context.Context, request *pbTag.UpdateTagRequest) (*pbTag.Tag, error) {
	response, err := s.tagService.UpdateTag(request)
	if err != nil {
		logger.Errorf("failed to update tag: %s", err)
		return nil, err
	}

	return response, nil
}

// DeleteTag implements service.ServiceServer
func (s *server) DeleteTag(ctx context.Context, request *pbTag.TagId) (*emptypb.Empty, error) {
	err := s.tagService.DeleteTag(request)
	if err != nil {
		logger.Errorf("failed to delete a tag: %s", err)
		return nil, err
	}

	return &empty.Empty{}, err
}

func newServer(
	tagService *services.TagService,
) *server {
	validator, err := protovalidate.New()
	if err != nil {
		fmt.Println("failed to initialize validator:", err)
	}

	s := &server{
		tagService: tagService,
		validator:  validator,
	}
	return s
}

func newUnaryInterceptor() (grpc.UnaryServerInterceptor, error) {
	// Create a new validator
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %v", err)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Infof("Handling gRPC request for method %s", info.FullMethod)

		// Type assert req to a protoreflect.ProtoMessage
		if p, ok := req.(protoreflect.ProtoMessage); ok {
			// Validate the request
			if err := validator.Validate(p); err != nil {
				return nil, err
			}
		}

		resp, err := handler(ctx, req)
		logger.Infof("Completed gRPC request for method %s", info.FullMethod)
		return resp, err
	}, nil
}

func loggingHandler(h http.Handler, logMessage string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log before handling the request
		logger.Infof("Handling request for: %s with log message: %s", r.URL.Path, logMessage)

		// Pass the request to the original handler
		h.ServeHTTP(w, r)

		// Log after handling the request if you want
	})
}

// func serveSwaggerUI(gwmux *runtime.ServeMux, httpMux *http.ServeMux) {
// 	swaggerDir := "../docs"

// 	logger.Infof("grpc gateway swagger doc is ready")
// 	// Serving swagger UI
// 	httpMux.Handle("/doc", loggingHandler(http.StripPrefix("/doc", http.FileServer(http.Dir(swaggerDir))), "Serving Swagger UI"))
// 	// Serving swagger doc
// 	httpMux.HandleFunc("doc.json", func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(docs.SwaggerJSON)
// 	})
// }

func mergeHandlers(gwmux *runtime.ServeMux, httpMux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/doc" || r.URL.Path == "/doc.json" {
			httpMux.ServeHTTP(w, r)
			return
		}
		gwmux.ServeHTTP(w, r)
	})
}

func StartServer(
	tagService *services.TagService,
) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("GRPC_SERVER_PORT")))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	// creds, err := credentials.NewServerTLSFromFile(Path(constants.CERT_FILE), Path(constants.KEY_FILE))
	// if err != nil {
	// 	logger.Fatalf("Failed to generate credentials: %v", err)
	// }
	// opts = append(opts, grpc.Creds(creds))

	unaryInterceptor, err := newUnaryInterceptor()
	if err != nil {
		logger.Fatalf("failed to create interceptor: %v", err)
	}
	opts = append(opts, grpc.UnaryInterceptor(unaryInterceptor))
	s := grpc.NewServer(opts...)

	serverInstance := newServer(tagService)
	ServiceServer.RegisterServiceServer(s, serverInstance)
	logger.Infof("grpc server listening at %v", lis.Addr())
	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Fatalf("failed to serve: %v", err)
		}
	}()

	// grpc gateway
	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	// conn, err := grpc.DialContext(
	// 	context.Background(),
	// 	fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("GRPC_SERVER_PORT")),
	// 	grpc.WithBlock(),
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )
	// if err != nil {
	// 	logger.Errorf("Failed to dial server:", err)
	// }

	// gwmux := runtime.NewServeMux()
	// Register Service handler
	// err = ServiceServer.RegisterServiceHandler(context.Background(), gwmux, conn)
	// if err != nil {
	// 	logger.Errorf("Failed to register gateway:", err)
	// }

	// Serve Swagger UI and doc
	// httpMux := http.NewServeMux()
	// httpMux.Handle("/", gwmux)

	// swaggerFiles.ReadFile("../docs/service.swagger.json")
	// httpMux.Handle("/docs/", loggingHandler(swaggerFiles.Handler, "swagger doc"))

	// gwServer := &http.Server{
	// 	Addr:    fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
	// 	Handler: httpMux,
	// }

	// logger.Infof("Serving gRPC-Gateway on %s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))
	// if err := gwServer.ListenAndServe(); err != nil {
	// 	logger.Fatalf("failed to serve grpc gateway: %v", err)
	// }
}
