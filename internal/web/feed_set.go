package web

import (
	"context"
	"errors"

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
