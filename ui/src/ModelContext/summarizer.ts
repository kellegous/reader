import { FetchRPC } from "twirp-ts";
import { ReaderClientJSON } from "../gen/reader.twirp";

export type Role = "user" | "assistant" | "system";

export interface Message {
  role: Role;
  content: string;
}

interface Event {
  created_at: string;
  done: boolean;
  message: Message;
  done_reason: string;
  model: string;
}

export class Summarizer {
  constructor(
    private readonly client: ReaderClientJSON,
    private readonly baseUrl: string,
    private readonly model: string
  ) {}

  async summarize(
    entryId: bigint,
    setSummary: (summary: string) => void
  ): Promise<string> {
    const { client, baseUrl, model } = this;
    const { text } = await client.GetEntryText({ entryId });
    return streamSummary(
      await requestSummary(baseUrl, model, text),
      setSummary
    );
  }

  static async createIfAvailable(
    baseUrl: string,
    model: string
  ): Promise<Summarizer | null> {
    try {
      await fetch(`${baseUrl}/api/ps`);
      return new Summarizer(
        new ReaderClientJSON(FetchRPC({ baseUrl: "/twirp" })),
        baseUrl,
        model
      );
    } catch {
      return null;
    }
  }
}

const requestSummary = (baseUrl: string, model: string, content: string) => {
  const messages = [
    {
      role: "system",
      content: "You are a helpful assistant that summarizes text.",
    },
    {
      role: "user",
      content: `Please give a single sentence summary of the following text:\n\n${content}`,
    },
  ];

  return fetch(`${baseUrl}/api/chat`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      model: model,
      messages,
      stream: true,
    }),
  });
};

const streamSummary = async (
  res: Response,
  setSummary: (summary: string) => void
) => {
  const reader = res.body?.getReader();
  if (!reader) {
    return "";
  }

  try {
    let summary = "";
    while (true) {
      const chunk = await reader.read();
      try {
        const event: Event = JSON.parse(new TextDecoder().decode(chunk.value));
        summary += event.message.content;
        setSummary(summary);

        if (event.done) {
          return summary;
        }
      } catch (e) {
        console.error(new TextDecoder().decode(chunk.value));
        throw e;
      }
    }
  } finally {
    reader.releaseLock();
  }
};
