package domain

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"todo/pkg/http/responses"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrInvalidID  = errors.New("invalid id was provided")
	ErrBadRequest = errors.New("bad request was sent")
	ErrEmptyName  = errors.New("name of task must be set")
)

func HandleError(err error, r any) responses.Response {
	switch {
	case err == nil:
		return responses.OK(r)
	case errors.Is(err, ErrNotFound):
		return responses.NotFound(err)
	case errors.Is(err, ErrInvalidID):
		return responses.BadRequest(err)
	case errors.Is(err, ErrBadRequest):
		return responses.BadRequest(err)
	case errors.Is(err, ErrEmptyName):
		return responses.BadRequest(err)
	default:
		return responses.Unknown(err)
	}
}

func HandleGRPCError(err error) error {
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.OK:
			return nil
		case codes.NotFound:
			return ErrNotFound
		case codes.InvalidArgument:
			return ErrEmptyName
		default:
			return errors.New(st.String())
		}
	}
	return err
}
