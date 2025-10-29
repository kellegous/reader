import * as proto from "../gen/reader";
import { Entry } from "./Entry";

export interface FeedProps {
  feed: proto.Feed;
  entries: proto.Entry[];
}

export const Feed = ({ feed, entries }: FeedProps) => {
  return (
    <div>
      <div>{feed.title}</div>
      <div>
        {entries.map((entry) => (
          <Entry key={entry.id} entry={entry} />
        ))}
      </div>
    </div>
  );
};
