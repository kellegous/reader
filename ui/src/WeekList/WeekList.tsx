import { useReaderData } from "../ReaderDataContext";
import styles from "./WeekList.module.scss";

export const WeekList = () => {
  const { weeks } = useReaderData();

  return (
    <div className={styles.root}>
      {weeks.map(({ week, entries }) => (
        <div key={week.startsAt.getTime()}>
          <h2>
            {week.toString()} ({entries.length})
          </h2>
          <ul>
            {entries.map((entry) => (
              <li key={entry.id}>
                <a href={entry.url} target="_blank" rel="noopener noreferrer">
                  {entry.feed?.title} / {entry.title}
                </a>
              </li>
            ))}
          </ul>
        </div>
      ))}
    </div>
  );
};
