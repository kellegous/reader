import { useReaderData } from "../ReaderDataContext";
import styles from "./Weeks.module.scss";
import { Entry, Feed } from "../gen/reader";
import { useMemo } from "react";
import { Week } from "./Week";

export const Weeks = () => {
  const { weeks } = useReaderData();

  const groupedWeeks = useMemo(() => {
    return weeks.map(({ week, entries }) => ({
      week,
      feeds: groupByFeed(entries),
    }));
  }, [weeks]);

  return (
    <div className={styles.root}>
      {groupedWeeks.map(({ week, feeds }) => (
        <Week key={week.startsAt.getTime()} week={week} feeds={feeds} />
      ))}
    </div>
  );
};

const groupByFeed = (entries: Entry[]) => {
  const byFeedId = new Map<bigint, { feed: Feed; entries: Entry[] }>();
  for (const entry of entries) {
    const f = entry.feed!;
    const { feed, entries } = byFeedId.get(f.id) ?? { feed: f, entries: [] };
    entries.push(entry);
    byFeedId.set(f.id, { feed, entries });
  }
  return Array.from(byFeedId.values()).sort((a, b) =>
    a.feed.title.localeCompare(b.feed.title)
  );
};
