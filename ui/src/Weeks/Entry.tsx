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

  const [status, setStatus] = useState<proto.Status>(entry.status);
  const [showSummary, setShowSummary] = useState(false);

  const {
    available: summaryAvailable,
    summarize,
    summary,
    loading: summaryLoading,
  } = useSummary(entry.id);

  const handleClick = useCallback(() => {
    setStatus(proto.Status.READ);
  }, []);

  const toggleSummary = useCallback(
    (e: React.MouseEvent<HTMLAnchorElement>) => {
      e.preventDefault();
      setShowSummary(!showSummary);
      if (!showSummary) {
        summarize();
      }
    },
    [showSummary, summarize]
  );

  // TODO(kellegous): we have to submit the rpc request for this.
  const toggleStatus = useCallback(
    (e: React.MouseEvent<HTMLAnchorElement>) => {
      e.preventDefault();
      setStatus(
        status === proto.Status.READ ? proto.Status.UNREAD : proto.Status.READ
      );
    },
    [status]
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
        <div>
          <a href="#" onClick={toggleStatus}>
            {status === proto.Status.READ ? "mark unread" : "mark read"}
          </a>
        </div>
        {summaryAvailable && (
          <div>
            <a href="#" onClick={toggleSummary}>
              {showSummary ? "hide summary" : "show summary"}
            </a>
          </div>
        )}
      </div>
      {showSummary && (
        <div className={styles.summary}>
          <div>{summary}</div>
          {!summaryLoading && (
            <a href="#" onClick={toggleSummary} className={styles.hide}>
              hide
            </a>
          )}
        </div>
      )}
    </div>
  );
};

const classForStatus = (status: proto.Status) => {
  switch (status) {
    case proto.Status.READ:
      return styles.read;
    case proto.Status.REMOVED:
      return styles.removed;
    default:
      return "";
  }
};
