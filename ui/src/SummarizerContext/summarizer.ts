import { FetchRPC } from "twirp-ts";
import { ReaderClientJSON } from "../gen/reader.twirp";

export class Summarizer {
  constructor(
    private readonly client: ReaderClientJSON,
    private readonly baseUrl: string
  ) {}

  async summarize(entryId: bigint): Promise<string> {
    const { client, baseUrl } = this;
    console.log(baseUrl);
    const { text } = await client.GetEntryText({ entryId });
    // TODO(kellegous): Catch ollama for an LLM summary.
    return text;
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
