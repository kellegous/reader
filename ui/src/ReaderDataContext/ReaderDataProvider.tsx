import { FetchRPC } from "twirp-ts";
import { ReaderClientJSON } from "../gen/reader.twirp";
import { ReaderDataContext, ReaderDataState } from "./ReaderDataContext";
import { useState, useEffect } from "react";
import { Week, Weekday } from "../time";
import { Timestamp } from "../gen/google/protobuf/timestamp";
import {
  GetEntriesRequest_SortKey,
  GetEntriesRequest_Order,
  Entry,
} from "../gen/reader";

export interface ReaderDataProviderProps {
  until: Date;
  numWeeks: number;
  weekday: Weekday;
  children: React.ReactNode;
}

export const ReaderDataProvider = ({
  until,
  numWeeks,
  weekday,
  children,
}: ReaderDataProviderProps) => {
  const [state, setState] = useState<ReaderDataState>({
    me: null,
    weeks: [],
    loading: false,
    until,
    numWeeks,
  });

  useEffect(() => {
    const client = new ReaderClientJSON(FetchRPC({ baseUrl: "/twirp" }));
    loadState(client, until, numWeeks, weekday, setState);
  }, [until, numWeeks, weekday]);

  return (
    <ReaderDataContext.Provider value={state}>
      {children}
    </ReaderDataContext.Provider>
  );
};

const loadState = async (
  client: ReaderClientJSON,
  until: Date,
  numWeeks: number,
  weekday: Weekday,
  setState: (state: ReaderDataState) => void
) => {
  let state: ReaderDataState = {
    me: null,
    weeks: [],
    loading: true,
    until,
    numWeeks,
  };

  setState(state);

  const latest = Week.of(until, weekday);
  const earliest = latest.add(-numWeeks);

  await Promise.all([
    client.GetMe({}).then(({ user }) => {
      state = { ...state, me: user ?? null };
      setState(state);
    }),
    client
      .GetEntries({
        publishedAfter: Timestamp.fromDate(earliest.startsAt),
        publishedBefore: Timestamp.fromDate(latest.endsAt),
        sortKey: GetEntriesRequest_SortKey.PUBLISHED_AT,
        order: GetEntriesRequest_Order.DESC,
        includeContent: false,
      })
      .then(({ entries }) => {
        state = {
          ...state,
          weeks: Array.from(toWeeks(latest, earliest, weekday, entries)),
        };
        setState(state);
      }),
  ]);

  state = { ...state, loading: false };
  setState(state);
};

function* toWeeks(
  latest: Week,
  earliest: Week,
  weekday: Weekday,
  entries: Entry[]
) {
  const byWeek = new Map<number, Entry[]>();

  for (const entry of entries) {
    const key = Week.of(
      Timestamp.toDate(entry.publishedAt!),
      weekday
    ).startsAt.getTime();
    const entries = byWeek.get(key) ?? [];
    entries.push(entry);
    byWeek.set(key, entries);
  }

  for (
    let week = latest;
    week.startsAt >= earliest.startsAt;
    week = week.add(-1)
  ) {
    yield { week, entries: byWeek.get(week.startsAt.getTime()) ?? [] };
  }
}
