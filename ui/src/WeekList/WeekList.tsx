import { useReaderData } from "../ReaderDataContext";

export const WeekList = () => {
  const { weeks } = useReaderData();

  return weeks.map(({ week }) => (
    <div key={week.startsAt.getTime()}>{week.toString()}</div>
  ));
};
