package server

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

type healthRes struct {
	Body struct {
		Detail string `json:"detail"`
	}
}

func registerHealthRoutes(s *server) {
	huma.Get(s.huma, "/health", func(ctx context.Context, req *struct{}) (*healthRes, error) {
		err := s.db.PingContext(ctx)
		if err != nil {
			return nil, err
		}
		return &healthRes{
			Body: struct {
				Detail string `json:"detail"`
			}{
				Detail: "OK",
			},
		}, nil
	})
}
