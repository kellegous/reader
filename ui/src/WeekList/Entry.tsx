import * as proto from "../gen/reader";

export interface EntryProps {
  entry: proto.Entry;
}

export const Entry = ({ entry }: EntryProps) => {
  return (
    <div>
      <a href={entry.url} target="_blank" rel="noopener noreferrer">
        {entry.title}
      </a>
    </div>
  );
};
