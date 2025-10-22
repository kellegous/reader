package web

import (
	"context"

	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/kellegous/reader"
)

type rpc struct {
}

var _ reader.Reader = (*rpc)(nil)

func (r *rpc) CheckHealth(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, twirp.NewError(twirp.Unimplemented, "not implemented")
}
