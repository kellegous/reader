import * as proto from "../gen/reader";
import { formatElapsedTime } from "../elapsed-time";
import { Timestamp } from "../gen/google/protobuf/timestamp";
import styles from "./Entry.module.scss";
import { useCallback, useState } from "react";
import { useSummarizer } from "../SummarizerContext";

const enableSummaries = false;
export interface EntryProps {
  entry: proto.Entry;
}

export const Entry = ({ entry }: EntryProps) => {
  const url = `/unread/entry/${entry.id}`;

  const [status, setStatus] = useState<proto.Entry_Status>(entry.status);

  const { available: summarizerAvailable } = useSummarizer();

  const handleClick = useCallback(() => {
    setStatus(proto.Entry_Status.READ);
  }, []);

  return (
    <div className={`${styles.root} ${classForStatus(status)}`}>
      <div className={styles.title}>
        <a
          href={url}
          target="_blank"
          rel="noopener noreferrer"
          onClick={handleClick}
        >
          {entry.title}
        </a>
      </div>
      <div className={styles.info}>
        <div>{formatElapsedTime(Timestamp.toDate(entry.publishedAt!))}</div>
        <div>{`${entry.readingTime} min`}</div>
        {summarizerAvailable && enableSummaries && <div>summary</div>}
      </div>
    </div>
  );
};

const classForStatus = (status: proto.Entry_Status) => {
  switch (status) {
    case proto.Entry_Status.READ:
      return styles.read;
    case proto.Entry_Status.REMOVED:
      return styles.removed;
    default:
      return "";
  }
};
