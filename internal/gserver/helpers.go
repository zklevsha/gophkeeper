package gserver

import (
	"github.com/zklevsha/gophkeeper/internal/structs"
	"google.golang.org/grpc/codes"
)

// getCode returns gRCP error code based on error type
func getCode(e error) codes.Code {
	switch e {
	case structs.ErrUserAlreadyExists, structs.ErrPdataAlreatyEsists:
		return codes.AlreadyExists
	case structs.ErrUserAuth, structs.ErrInvalidToken, structs.ErrNoToken:
		return codes.Unauthenticated
	case structs.ErrPdataNotFound:
		return codes.NotFound
	default:
		return codes.Internal
	}

}
