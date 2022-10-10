package structs

import "github.com/zklevsha/gophkeeper/internal/pb"

// Gclient represents collection of various GRPC clients
type Gclient struct {
	Auth pb.AuthClient
}
