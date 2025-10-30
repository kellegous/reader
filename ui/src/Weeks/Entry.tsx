import * as proto from "../gen/reader";
import { formatElapsedTime } from "../elapsed-time";
import { Timestamp } from "../gen/google/protobuf/timestamp";
import styles from "./Entry.module.scss";

export interface EntryProps {
  entry: proto.Entry;
}

export const Entry = ({ entry }: EntryProps) => {
  return (
    <div className={styles.root}>
      <div className={styles.title}>
        <a href={entry.url} target="_blank" rel="noopener noreferrer">
          {entry.title}
        </a>
      </div>
      <div className={styles.info}>
        <div>{formatElapsedTime(Timestamp.toDate(entry.publishedAt!))}</div>
        <div>{`${entry.readingTime} min`}</div>
      </div>
    </div>
  );
};
