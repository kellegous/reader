import { timestampDate, timestampFromDate } from "@bufbuild/protobuf/wkt";
import { Client } from "@connectrpc/connect";
import {
  Config,
  Feed,
  GetEntriesRequest_Order,
  GetEntriesRequest_SortKey,
  Entry as ProtoEntry,
  Reader,
  Status,
  User,
} from "../gen/reader_pb";
import { Week, Weekday } from "../time";
import { Summarizer } from "./summarizer";

const defaultModel = "gemma4:e4b";

export interface ModelState {
  client: Client<typeof Reader>;
  until: Date;
  numWeeks: number;
  weekday: Weekday;
  user: User | null;
  config: Config | null;
  weeks: { week: Week; entries: Entry[] }[];
  summarizer: Summarizer | null;
  loading: boolean;
  refresh: () => Promise<void>;
  updateEntryStatus: (entryId: bigint, status: Status) => Promise<void>;
}

export interface Entry {
  id: bigint;
  publishedAt: Date;
  changedAt: Date;
  createdAt: Date;
  feed: Feed;
  url: string;
  title: string;
  content: string;
  readingTime: number;
  status: Status;
}

export const empty = (client: Client<typeof Reader>): ModelState => {
  const updateEntryStatus = async (
    entryId: bigint,
    status: Status,
  ): Promise<void> => {
    await client.setEntryStatus({ entryId, status });
  };

  return {
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
    updateEntryStatus,
  };
};

export const load = async (
  client: Client<typeof Reader>,
  until: Date,
  numWeeks: number,
  weekday: Weekday,
  setState: (fn: (model: ModelState) => ModelState) => void,
): Promise<void> => {
  setState((model) => ({ ...model, loading: true }));

  const [, config] = await Promise.all([
    getUser(client).then((user) => setState((model) => ({ ...model, user }))),
    getConfig(client).then((config) => {
      setState((model) => ({ ...model, config }));
      return config;
    }),
    getWeeks(client, until, numWeeks, weekday).then(({ weeks }) =>
      setState((model) => ({ ...model, weeks })),
    ),
  ]).finally(() => setState((model) => ({ ...model, loading: false })));

  const summarizer = await getSummarizer(config);
  setState((model) => ({ ...model, summarizer }));
};

const getUser = async (client: Client<typeof Reader>): Promise<User> => {
  const { user } = await client.getMe({});
  return user!;
};

export const getConfig = async (
  client: Client<typeof Reader>,
): Promise<Config> => {
  const { config } = await client.getConfig({});
  return config!;
};

const getWeeks = async (
  client: Client<typeof Reader>,
  until: Date,
  numWeeks: number,
  weekday: Weekday,
): Promise<{ weeks: { week: Week; entries: Entry[] }[] }> => {
  const latest = Week.of(until, weekday);
  const earliest = latest.add(-numWeeks);
  const { entries, feeds } = await client.getEntries({
    publishedAfter: timestampFromDate(earliest.startsAt),
    publishedBefore: timestampFromDate(latest.endsAt),
    sortKey: GetEntriesRequest_SortKey.PUBLISHED_AT,
    order: GetEntriesRequest_Order.DESC,
    includeContent: false,
  });

  const feedById = new Map(feeds.map((f) => [f.id, f]));

  return {
    weeks: Array.from(
      toWeeks(
        latest,
        earliest,
        weekday,
        entries.map((e) => toEntry(e, feedById)),
      ),
    ),
  };
};

const getSummarizer = async (config: Config): Promise<Summarizer | null> => {
  if (!config) {
    return null;
  }
  const { url, model } = config.ollama!;
  if (!url) {
    return null;
  }

  return await Summarizer.createIfAvailable(url, model || defaultModel);
};

const toEntry = (
  {
    id,
    publishedAt,
    changedAt,
    createdAt,
    feedId,
    url,
    title,
    content,
    readingTime,
    status,
  }: ProtoEntry,
  feedById: Map<bigint, Feed>,
): Entry => {
  const feed = feedById.get(feedId);
  if (!feed) {
    throw new Error(`Feed not found for entry ${id}`);
  }

  return {
    id,
    publishedAt: timestampDate(publishedAt!),
    changedAt: timestampDate(changedAt!),
    createdAt: timestampDate(createdAt!),
    feed,
    url: url,
    title: title,
    content: content,
    readingTime: readingTime,
    status: status,
  };
};

function* toWeeks(
  latest: Week,
  earliest: Week,
  weekday: Weekday,
  entries: Entry[],
) {
  const byWeek = new Map<number, Entry[]>();

  for (const entry of entries) {
    const key = Week.of(entry.publishedAt, weekday).startsAt.getTime();
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
