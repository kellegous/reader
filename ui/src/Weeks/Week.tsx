import * as time from "../time";
import * as proto from "../gen/reader";
import { Feed } from "./Feed";

export interface WeekProps {
  week: time.Week;
  feeds: { feed: proto.Feed; entries: proto.Entry[] }[];
}

export const Week = ({ week, feeds }: WeekProps) => {
  return (
    <div>
      <h2>{week.toString()}</h2>
      <div>
        {feeds.map(({ feed, entries }) => (
          <Feed key={feed.id} feed={feed} entries={entries} />
        ))}
      </div>
    </div>
  );
};
