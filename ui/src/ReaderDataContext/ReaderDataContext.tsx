import { User, Entry } from "../gen/reader";
import { createContext } from "react";

export interface ReaderDataState {
  me: User | null;
  entries: Entry[];
  loading: boolean;
}

export const ReaderDataContext = createContext<ReaderDataState>({
  me: null,
  entries: [],
  loading: false,
});
