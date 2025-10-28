import { useReaderData } from "../ReaderDataContext";

export const WeekList = () => {
  const { weeks } = useReaderData();

  return weeks.map(({ week, entries }) => (
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
  ));
};
