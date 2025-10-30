import * as time from "../time";
import * as proto from "../gen/reader";
import { Feed } from "./Feed";
import styles from "./Week.module.scss";

const dayFormat = new Intl.DateTimeFormat("en-US", {
  month: "short",
  day: "numeric",
});

const formatWeek = (week: time.Week) =>
  `${dayFormat.format(week.startsAt)} - ${dayFormat.format(week.endsAt)}`;

export interface WeekProps {
  week: time.Week;
  feeds: { feed: proto.Feed; entries: proto.Entry[] }[];
}

export const Week = ({ week, feeds }: WeekProps) => {
  return (
    <div className={styles.root}>
      <div className={styles.title}>{formatWeek(week)}</div>
      <div className={styles.feeds}>
        {feeds.map(({ feed, entries }) => (
          <Feed key={feed.id} feed={feed} entries={entries} />
        ))}
      </div>
    </div>
  );
};
