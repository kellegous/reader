import { FetchRPC } from "twirp-ts";
import { ReaderClientJSON } from "../gen/reader.twirp";

export class Summarizer {
  private readonly store = new Map<bigint, string>();

  constructor(
    private readonly client: ReaderClientJSON,
    private readonly baseUrl: string
  ) {}

  async summarize(entryId: bigint): Promise<string> {
    const { client, store, baseUrl } = this;
    console.log(baseUrl);
    return store.get(entryId) ?? (await client.GetEntryText({ entryId })).text;
  }

  static async createIfAvailable(baseUrl: string): Promise<Summarizer | null> {
    try {
      await fetch(`${baseUrl}/api/ps`);
      return new Summarizer(
        new ReaderClientJSON(FetchRPC({ baseUrl })),
        baseUrl
      );
    } catch {
      return null;
    }
  }
}
