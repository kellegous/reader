package web

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/kellegous/glue/logging"
	"github.com/kellegous/poop"
	"github.com/kellegous/reader/internal/plaintext"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"miniflux.app/v2/client"

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

func (r *rpc) GetEntries(ctx context.Context, req *reader.GetEntriesRequest) (*reader.GetEntriesResponse, error) {
	res, err := r.client.EntriesContext(ctx, &client.Filter{
		PublishedAfter:  req.PublishedAfter.AsTime().Unix(),
		PublishedBefore: req.PublishedBefore.AsTime().Unix(),
		Order:           strings.ToLower(req.SortKey.String()),
		Direction:       strings.ToLower(req.Order.String()),
	})
	if err != nil {
		return nil, newBackendError(ctx, err)
	}

	fs := newFeedSet(r.client)

	entries := make([]*reader.Entry, 0, len(res.Entries))
	for _, entry := range res.Entries {
		e, err := toEntry(entry, fs.toFeed, req.GetIncludeContent())
		if err != nil {
			return nil, newBackendError(ctx, err)
		}
		entries = append(entries, e)
	}

	if err := fs.resolveIcons(ctx); err != nil {
		return nil, newBackendError(ctx, err)
	}

	return &reader.GetEntriesResponse{
		Entries: entries,
	}, nil
}

func (r *rpc) GetMe(ctx context.Context, req *emptypb.Empty) (*reader.GetMeResponse, error) {
	user, err := r.client.MeContext(ctx)
	if err != nil {
		return nil, newBackendError(ctx, err)
	}
	return &reader.GetMeResponse{
		User: toUser(user),
	}, nil
}

func (r *rpc) GetEntryText(ctx context.Context, req *reader.GetEntryTextRequest) (*reader.GetEntryTextResponse, error) {
	entry, err := r.client.EntryContext(ctx, req.EntryId)
	if err != nil {
		return nil, newBackendError(ctx, err)
	}

	return &reader.GetEntryTextResponse{
		Text: plaintext.From(entry.Content),
	}, nil
}

func toUser(user *client.User) *reader.User {
	return &reader.User{
		Id:       user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		Theme:    user.Theme,
		Language: user.Language,
		Timezone: user.Timezone,
	}
}

func toFeed(feed *client.Feed) *reader.Feed {
	return &reader.Feed{
		Id:      feed.ID,
		FeedUrl: feed.FeedURL,
		SiteUrl: feed.SiteURL,
		Title:   feed.Title,
	}
}

func toStatus(status string) (reader.Entry_Status, error) {
	s, ok := reader.Entry_Status_value[strings.ToUpper(status)]
	if !ok {
		return 0, fmt.Errorf("invalid status: %s", status)
	}
	return reader.Entry_Status(s), nil
}

func toEntry(
	entry *client.Entry,
	toFeed func(*client.Feed) *reader.Feed,
	includeContent bool,
) (*reader.Entry, error) {
	var content string
	if includeContent {
		content = entry.Content
	}

	status, err := toStatus(entry.Status)
	if err != nil {
		return nil, poop.Chain(err)
	}

	return &reader.Entry{
		Id:          entry.ID,
		PublishedAt: timestamppb.New(entry.Date),
		ChangedAt:   timestamppb.New(entry.ChangedAt),
		CreatedAt:   timestamppb.New(entry.CreatedAt),
		Feed:        toFeed(entry.Feed),
		Url:         entry.URL,
		Title:       entry.Title,
		Content:     content,
		ReadingTime: int32(entry.ReadingTime),
		Status:      status,
	}, nil
}
