import { User, Entry, Config } from "../gen/reader";
import { createContext } from "react";
import { Week } from "../time";

export interface ReaderDataState {
  me: User | null;
  config: Config | null;
  weeks: { week: Week; entries: Entry[] }[];
  loading: boolean;
  until: Date;
  numWeeks: number;
  refresh: () => Promise<void>;
}

export const ReaderDataContext = createContext<ReaderDataState>({
  me: null,
  config: null,
  weeks: [],
  loading: false,
  until: new Date(),
  numWeeks: 5,
  refresh: () => Promise.resolve(),
});
