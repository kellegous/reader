package web

import (
	"context"
	"errors"
	"sort"

	"github.com/kellegous/poop"
	"golang.org/x/sync/errgroup"
	"miniflux.app/v2/client"

	"github.com/kellegous/reader"
)

type feedSet struct {
	feeds  map[int64]*reader.Feed
	client *client.Client
}

func newFeedSet(client *client.Client) *feedSet {
	return &feedSet{
		feeds:  make(map[int64]*reader.Feed),
		client: client,
	}
}

func (s feedSet) toFeed(feed *client.Feed) *reader.Feed {
	if f, ok := s.feeds[feed.ID]; ok {
		return f
	}
	f := toFeed(feed)
	s.feeds[feed.ID] = f
	return f
}

func (s feedSet) toFeeds(ctx context.Context) ([]*reader.Feed, error) {
	if err := s.resolveIcons(ctx); err != nil {
		return nil, err
	}
	feeds := make([]*reader.Feed, 0, len(s.feeds))
	for _, feed := range s.feeds {
		feeds = append(feeds, feed)
	}
	sort.Slice(feeds, func(i, j int) bool {
		return feeds[i].Id < feeds[j].Id
	})
	return feeds, nil
}

func (s feedSet) resolveIcons(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, feed := range s.feeds {
		g.Go(func() error {
			icon, err := s.client.FeedIconContext(ctx, feed.Id)
			if errors.Is(err, client.ErrNotFound) {
				return nil
			} else if err != nil {
				return poop.Chain(err)
			}
			feed.IconDataUrl = "data://" + icon.Data
			return nil
		})
	}

	return g.Wait()
}
