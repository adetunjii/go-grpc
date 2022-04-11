package service

import (
	"context"
	"errors"
	"github.com/Adetunjii/go-grpc/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type LaptopServer struct {
	Store LaptopStore
	*pb.UnimplementedLaptopServiceServer
}

func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{
		Store: store,
	}
}

// CreateLaptop
// Unary RPC to create a new laptop
func (server *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreatelaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()
	log.Printf("recieved a create laptop request wit id: %s", laptop.Id)

	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID: %v", err)
		}
	} else {
		// generate new id
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop id: %v", err)
		}
		laptop.Id = id.String()
	}

	time.Sleep(6 * time.Second)

	if ctx.Err() == context.Canceled {
		log.Print("cancelled")
		return nil, status.Error(codes.Canceled, "cancelled")
	}

	if ctx.Err() == context.DeadlineExceeded {
		log.Print("error timeout")
		return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	}

	err := server.Store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, DuplicateException) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to the store %v", err)
	}
	log.Printf("saved laptop with id: %s", laptop.Id)
	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}

	return res, nil
}
