package dto

import "google.golang.org/grpc/status"

type ErrorResponse struct {
	Error string `json:"error"`
}

func ErrorResponseFromGRPCError(err error) *ErrorResponse {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return &ErrorResponse{Error: "unknown error occurred"}
	}
	return &ErrorResponse{Error: st.Message()}
}
