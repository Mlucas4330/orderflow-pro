// cmd/product-service/main.go
package main

import (
	"context"
	"log"
	"net"

	pb "github.com/mlucas4330/orderflow-pro/pkg/productpb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedProductServiceServer
}

func (s *server) GetProductDetails(ctx context.Context, req *pb.GetProductDetailsRequest) (*pb.GetProductDetailsResponse, error) {
	productID := req.GetProductId()
	log.Printf("Recebida requisição para buscar detalhes do produto: %s", productID)

	return &pb.GetProductDetailsResponse{
		Id:    productID,
		Name:  "Café Especial Grão Mágico 250g (v_grpc)",
		Price: "55.90",
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Falha ao escutar na porta: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterProductServiceServer(grpcServer, &server{})

	log.Printf("Servidor gRPC escutando em %v", listener.Addr())

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Falha ao servir gRPC: %v", err)
	}
}
