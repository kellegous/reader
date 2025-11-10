import { FetchRPC } from "twirp-ts";
import { Timestamp } from "../gen/google/protobuf/timestamp";
import {
  Config,
  Entry,
  GetEntriesRequest_Order,
  GetEntriesRequest_SortKey,
  User,
} from "../gen/reader";
import { ReaderClientJSON } from "../gen/reader.twirp";
import { Week, Weekday } from "../time";

export class Model {
  private constructor(
    public readonly client: ReaderClientJSON,
    public readonly until: Date,
    public readonly numWeeks: number,
    public readonly weekday: Weekday,
    public readonly user: User,
    public readonly config: Config,
    public readonly weeks: { week: Week; entries: Entry[] }[]
  ) {}

  static async load(
    baseUrl: string,
    until: Date,
    numWeeks: number,
    weekday: Weekday
  ): Promise<Model> {
    const client = new ReaderClientJSON(FetchRPC({ baseUrl }));

    const latest = Week.of(until, weekday);
    const earliest = latest.add(-numWeeks);

    const [user, config, entries] = await Promise.all([
      client.GetMe({}).then(({ user }) => user!),
      client.GetConfig({}).then(({ config }) => config!),
      client
        .GetEntries({
          publishedAfter: Timestamp.fromDate(earliest.startsAt),
          publishedBefore: Timestamp.fromDate(latest.endsAt),
          sortKey: GetEntriesRequest_SortKey.PUBLISHED_AT,
          order: GetEntriesRequest_Order.DESC,
          includeContent: false,
        })
        .then(({ entries }) => entries),
    ]);

    return new Model(
      client,
      until,
      numWeeks,
      weekday,
      user,
      config,
      Array.from(toWeeks(latest, earliest, weekday, entries))
    );
  }
}

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
