package web

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/kellegous/glue/logging"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"miniflux.app/v2/client"

	"github.com/kellegous/reader"
	"github.com/kellegous/reader/internal"
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

	entries := make([]*reader.Entry, 0, len(res.Entries))
	for _, entry := range res.Entries {
		entries = append(entries, toEntry(entry))
	}
	return &reader.GetEntriesResponse{
		Entries: entries,
	}, nil
}

func toEntry(entry *client.Entry) *reader.Entry {
	feed := entry.Feed
	return &reader.Entry{
		Id:          entry.ID,
		PublishedAt: timestamppb.New(entry.Date),
		ChangedAt:   timestamppb.New(entry.ChangedAt),
		CreatedAt:   timestamppb.New(entry.CreatedAt),
		Feed: &reader.Feed{
			Id:      feed.ID,
			FeedUrl: feed.FeedURL,
			SiteUrl: feed.SiteURL,
			Title:   feed.Title,
		},
		Url:         entry.URL,
		Title:       entry.Title,
		Content:     entry.Content,
		ReadingTime: int32(entry.ReadingTime),
	}
}

func (r *rpc) GetEntriesByWeek(ctx context.Context, req *reader.GetEntriesByWeekRequest) (*reader.GetEntriesByWeekResponse, error) {
	user, err := r.client.MeContext(ctx)
	if err != nil {
		return nil, newBackendError(ctx, err)
	}

	loc, err := time.LoadLocation(user.Timezone)
	if err != nil {
		return nil, newBackendError(ctx, err)
	}

	var publishedBefore time.Time
	var publishedAfter time.Time

	switch t := req.Range.(type) {
	case *reader.GetEntriesByWeekRequest_FromWeekToWeek_:
		v := t.FromWeekToWeek
		publishedAfter = internal.WeekOf(v.FromWeekOf.AsTime(), time.Weekday(req.WeekStartsDay), loc).BeginsAt()
		publishedBefore = internal.WeekOf(v.ToWeekOf.AsTime(), time.Weekday(req.WeekStartsDay), loc).EndsAt()
	case *reader.GetEntriesByWeekRequest_NWeeksFromWeek_:
		v := t.NWeeksFromWeek
		after := internal.WeekOf(v.FromWeekOf.AsTime(), time.Weekday(req.WeekStartsDay), loc)
		before := after.Add(int(v.NWeeks))
		publishedAfter = after.BeginsAt()
		publishedBefore = before.EndsAt()
	default:
		return nil, twirp.InvalidArgumentError("range", fmt.Sprintf("invalid range: %T", req.Range))
	}

	res, err := r.client.EntriesContext(ctx, &client.Filter{
		PublishedAfter:  publishedAfter.Unix(),
		PublishedBefore: publishedBefore.Unix(),
		Order:           "published_at",
		Direction:       "desc",
	})
	if err != nil {
		return nil, newBackendError(ctx, err)
	}

	var weeksOfEntries []*reader.GetEntriesByWeekResponse_WeekOfEntries
	for _, entry := range res.Entries {
		feed := entry.Feed

		e := &reader.Entry{
			Id:          entry.ID,
			PublishedAt: timestamppb.New(entry.Date),
			ChangedAt:   timestamppb.New(entry.ChangedAt),
			CreatedAt:   timestamppb.New(entry.CreatedAt),
			Feed: &reader.Feed{
				Id:      feed.ID,
				FeedUrl: feed.FeedURL,
				SiteUrl: feed.SiteURL,
				Title:   feed.Title,
			},
			Url:         entry.URL,
			Title:       entry.Title,
			Content:     entry.Content,
			ReadingTime: int32(entry.ReadingTime),
		}

		week := internal.WeekOf(entry.Date, time.Weekday(req.WeekStartsDay), loc)
		if len(weeksOfEntries) == 0 {
			weeksOfEntries = append(weeksOfEntries, &reader.GetEntriesByWeekResponse_WeekOfEntries{
				Week: &reader.Week{
					BeginsAt: timestamppb.New(week.BeginsAt()),
					EndsAt:   timestamppb.New(week.EndsAt()),
				},
				Entries: []*reader.Entry{e},
			})
			continue
		}

		currentWeek := internal.WeekOf(weeksOfEntries[len(weeksOfEntries)-1].GetWeek().GetBeginsAt().AsTime(), time.Weekday(req.WeekStartsDay), loc)
		if !currentWeek.Equals(week) {
			weeksOfEntries = append(weeksOfEntries, &reader.GetEntriesByWeekResponse_WeekOfEntries{
				Week: &reader.Week{
					BeginsAt: timestamppb.New(week.BeginsAt()),
					EndsAt:   timestamppb.New(week.EndsAt()),
				},
				Entries: []*reader.Entry{e},
			})
			continue
		}

		weeksOfEntries[len(weeksOfEntries)-1].Entries = append(weeksOfEntries[len(weeksOfEntries)-1].Entries, e)
	}

	return &reader.GetEntriesByWeekResponse{
		WeeksOfEntries: weeksOfEntries,
	}, nil
}
