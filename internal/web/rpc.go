package web

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"miniflux.app/v2/client"

	"github.com/kellegous/glue/logging"
	"github.com/kellegous/reader"
)

type rpc struct {
	client *client.Client
}

var _ reader.Reader = (*rpc)(nil)

func idFrom(err error) string {
	hash := sha1.Sum([]byte(err.Error()))
	return hex.EncodeToString(hash[:8])
}

func newBackendError(ctx context.Context, err error) error {
	logging.L(ctx).Error("backend error", zap.Error(err))
	return twirp.NewError(twirp.Internal, fmt.Sprintf("backend error: %s", idFrom(err)))
}

func (r *rpc) CheckHealth(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	if err := r.client.HealthcheckContext(ctx); err != nil {
		return nil, newBackendError(ctx, err)
	}
	return &emptypb.Empty{}, nil
}

func (r *rpc) GetEntriesByWeek(ctx context.Context, req *reader.GetEntriesByWeekRequest) (*reader.GetEntriesByWeekResponse, error) {
	user, err := r.client.MeContext(ctx)
	if err != nil {
		return nil, newBackendError(ctx, err)
	}

	logging.L(ctx).Info("user",
		zap.String("username", user.Username),
		zap.String("timezone", user.Timezone),
	)

	return &reader.GetEntriesByWeekResponse{}, nil
}
