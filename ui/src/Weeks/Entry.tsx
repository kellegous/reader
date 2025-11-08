import * as proto from "../gen/reader";
import { formatElapsedTime } from "../elapsed-time";
import { Timestamp } from "../gen/google/protobuf/timestamp";
import styles from "./Entry.module.scss";
import { useCallback, useState } from "react";
import { useSummary } from "../SummarizerContext";
const enableSummaries = true;
export interface EntryProps {
  entry: proto.Entry;
}

export const Entry = ({ entry }: EntryProps) => {
  const url = `/unread/entry/${entry.id}`;

  const [status, setStatus] = useState<proto.Entry_Status>(entry.status);

  const {
    available: summaryAvailable,
    loading: summaryLoading,
    summarize,
    summary,
  } = useSummary(entry.id);

  const handleClick = useCallback(() => {
    setStatus(proto.Entry_Status.READ);
  }, []);

  const handleSummarize = useCallback(
    (e: React.MouseEvent<HTMLAnchorElement>) => {
      e.preventDefault();
      summarize();
    },
    [summarize]
  );

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
        {summaryAvailable && enableSummaries && (
          <a href="#" onClick={handleSummarize}>
            summarize
          </a>
        )}
      </div>
      <div
        className={styles.summary}
        style={{ display: summary || summaryLoading ? "block" : "none" }}
      >
        {summary}
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
