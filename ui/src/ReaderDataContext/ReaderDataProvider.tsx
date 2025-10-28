import { FetchRPC } from "twirp-ts";
import { ReaderClientJSON } from "../gen/reader.twirp";
import { ReaderDataContext, ReaderDataState } from "./ReaderDataContext";
import { useState, useEffect } from "react";
import { Week, Weekday } from "../time";
import { Timestamp } from "../gen/google/protobuf/timestamp";
import {
  GetEntriesRequest_SortKey,
  GetEntriesRequest_Order,
} from "../gen/reader";

export const ReaderDataProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [state, setState] = useState<ReaderDataState>({
    me: null,
    entries: [],
    loading: false,
  });

  useEffect(() => {
    const client = new ReaderClientJSON(FetchRPC({ baseUrl: "/twirp" }));
    loadState(client, setState);
  }, []);

  return (
    <ReaderDataContext.Provider value={state}>
      {children}
    </ReaderDataContext.Provider>
  );
};

const loadState = async (
  client: ReaderClientJSON,
  setState: (state: ReaderDataState) => void
) => {
  let state: ReaderDataState = {
    me: null,
    entries: [],
    loading: true,
  };

  setState(state);

  // TODO(kellegous): Let's start by just getting 5 weeks
  const end = Week.of(new Date(), Weekday.Monday);
  const start = end.add(-5);

  await Promise.all([
    client.GetMe({}).then(({ user }) => {
      state = { ...state, me: user ?? null };
      setState(state);
    }),
    client
      .GetEntries({
        publishedAfter: Timestamp.fromDate(start.startsAt),
        publishedBefore: Timestamp.fromDate(end.startsAt),
        sortKey: GetEntriesRequest_SortKey.PUBLISHED_AT,
        order: GetEntriesRequest_Order.DESC,
      })
      .then(({ entries }) => {
        state = { ...state, entries };
        setState(state);
      }),
  ]);

  state = { ...state, loading: false };
  setState(state);
};
