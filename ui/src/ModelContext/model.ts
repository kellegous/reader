import { FetchRPC } from "twirp-ts";
import { Timestamp } from "../gen/google/protobuf/timestamp";
import {
  Config,
  Entry,
  GetEntriesRequest_Order,
  GetEntriesRequest_SortKey,
  User,
} from "../gen/reader";
import { ReaderClient, ReaderClientJSON } from "../gen/reader.twirp";
import { Week, Weekday } from "../time";
import { Summarizer } from "./summarizer";

export interface ModelState {
  client: ReaderClient;
  until: Date;
  numWeeks: number;
  weekday: Weekday;
  user: User | null;
  config: Config | null;
  weeks: { week: Week; entries: Entry[] }[];
  summarizer: Summarizer | null;
  loading: boolean;
  refresh: () => Promise<void>;
}

export const empty = (client: ReaderClient): ModelState => ({
  client,
  until: new Date(),
  numWeeks: 5,
  weekday: Weekday.Monday,
  user: null,
  config: null,
  weeks: [],
  summarizer: null,
  loading: false,
  refresh: () => Promise.resolve(),
});

export const load = async (
  baseUrl: string,
  until: Date,
  numWeeks: number,
  weekday: Weekday,
  setState: (fn: (model: ModelState) => ModelState) => void
): Promise<void> => {
  const client = new ReaderClientJSON(FetchRPC({ baseUrl }));
  setState((model) => ({ ...model, loading: true }));

  const [, config] = await Promise.all([
    getUser(client).then((user) => setState((model) => ({ ...model, user }))),
    getConfig(client).then((config) => {
      setState((model) => ({ ...model, config }));
      return config;
    }),
    getWeeks(client, until, numWeeks, weekday).then(({ weeks }) =>
      setState((model) => ({ ...model, weeks }))
    ),
  ]).finally(() => setState((model) => ({ ...model, loading: false })));

  const summarizer = await getSummarizer(config);
  setState((model) => ({ ...model, summarizer }));
};

const getUser = async (client: ReaderClient): Promise<User> => {
  const { user } = await client.GetMe({});
  return user!;
};

export const getConfig = async (client: ReaderClient): Promise<Config> => {
  const { config } = await client.GetConfig({});
  return config!;
};

const getWeeks = async (
  client: ReaderClient,
  until: Date,
  numWeeks: number,
  weekday: Weekday
): Promise<{ weeks: { week: Week; entries: Entry[] }[] }> => {
  const latest = Week.of(until, weekday);
  const earliest = latest.add(-numWeeks);
  const entries = await client.GetEntries({
    publishedAfter: Timestamp.fromDate(earliest.startsAt),
    publishedBefore: Timestamp.fromDate(latest.endsAt),
    sortKey: GetEntriesRequest_SortKey.PUBLISHED_AT,
    order: GetEntriesRequest_Order.DESC,
    includeContent: false,
  });
  return {
    weeks: Array.from(toWeeks(latest, earliest, weekday, entries.entries)),
  };
};

const getSummarizer = async (config: Config): Promise<Summarizer | null> => {
  if (!config || !config.ollamaUrl) {
    return null;
  }
  return await Summarizer.createIfAvailable(config.ollamaUrl);
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
