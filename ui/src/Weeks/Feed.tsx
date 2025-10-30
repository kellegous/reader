import * as proto from "../gen/reader";
import { Entry } from "./Entry";
import { useState } from "react";
import styles from "./Feed.module.scss";
import { FeedIcon } from "./FeedIcon";

const DEFAULT_LIMIT = 5;

export interface FeedProps {
  feed: proto.Feed;
  entries: proto.Entry[];
  limit?: number;
}

export const Feed = ({ feed, entries, limit = DEFAULT_LIMIT }: FeedProps) => {
  const [showMore, setShowMore] = useState(false);
  const toDisplay = showMore ? entries : entries.slice(0, limit);

  const toggleShowMore = (e: React.MouseEvent<HTMLAnchorElement>) => {
    e.preventDefault();
    setShowMore((showMore) => !showMore);
  };

  return (
    <div className={styles.root}>
      <div className={styles.title}>
        <FeedIcon url={feed.iconDataUrl} title={feed.title} />
        <a href={feed.siteUrl} target="_blank" rel="noopener noreferrer">
          {feed.title}
        </a>
      </div>
      <div className={styles.entries}>
        {toDisplay.map((entry) => (
          <Entry key={entry.id} entry={entry} />
        ))}
        {entries.length > limit && (
          <a href="#" onClick={toggleShowMore}>
            {showMore ? "Show less" : "Show more"}
          </a>
        )}
      </div>
    </div>
  );
};
