import * as proto from "../gen/reader";
import { formatElapsedTime } from "../elapsed-time";
import { Timestamp } from "../gen/google/protobuf/timestamp";
import styles from "./Entry.module.scss";
import { useCallback, useState } from "react";
import { useSummary } from "../SummarizerContext";

export interface EntryProps {
  entry: proto.Entry;
}

export const Entry = ({ entry }: EntryProps) => {
  const url = `/unread/entry/${entry.id}`;

  const [status, setStatus] = useState<proto.Entry_Status>(entry.status);
  const [showSummary, setShowSummary] = useState(false);

  const {
    available: summaryAvailable,
    summarize,
    summary,
    loading: summaryLoading,
  } = useSummary(entry.id);

  const handleClick = useCallback(() => {
    setStatus(proto.Entry_Status.READ);
  }, []);

  const handleSummarize = useCallback(
    (e: React.MouseEvent<HTMLAnchorElement>) => {
      e.preventDefault();
      setShowSummary(true);
      summarize();
    },
    [summarize]
  );

  const handleHideSummary = useCallback(
    (e: React.MouseEvent<HTMLAnchorElement>) => {
      e.preventDefault();
      setShowSummary(false);
    },
    []
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
        {summaryAvailable && !showSummary && (
          <a href="#" onClick={handleSummarize} className={styles.summarize}>
            summarize
          </a>
        )}
      </div>
      {showSummary && (
        <div className={styles.summary}>
          <div>{summary}</div>
          {!summaryLoading && (
            <a href="#" onClick={handleHideSummary} className={styles.hide}>
              hide
            </a>
          )}
        </div>
      )}
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
