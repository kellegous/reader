package web

import (
	"context"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"miniflux.app/v2/client"

	"github.com/kellegous/reader"
)

type rpc struct {
	client *client.Client
}

var _ reader.Reader = (*rpc)(nil)

func (r *rpc) CheckHealth(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	if err := r.client.HealthcheckContext(ctx); err != nil {
		return nil, twirp.NewError(twirp.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
