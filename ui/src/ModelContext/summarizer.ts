import { Client, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { withSource } from "../withSource";
import { Reader } from "../gen/reader_pb";

export type Role = "user" | "assistant" | "system";

export interface Message {
  role: Role;
  content: string;
}

interface ChatCompletionChunk {
  choices: {
    index: number;
    delta: { content?: string | null };
  }[];
}

const apiUrl = (baseUrl: string, path: string) => {
  const base = baseUrl.replace(/\/+$/, "");
  return `${base.endsWith("/v1") ? base : `${base}/v1`}/${path}`;
};

export class Summarizer {
  constructor(
    private readonly client: Client<typeof Reader>,
    private readonly baseUrl: string,
    private readonly model: string,
  ) {}

  async summarize(
    entryId: bigint,
    setSummary: (summary: string) => void,
  ): Promise<string> {
    const { client, baseUrl, model } = this;
    const { text } = await client.getEntryText({ entryId });
    return streamSummary(
      await requestSummary(baseUrl, model, text),
      setSummary,
    );
  }

  static async createIfAvailable(
    baseUrl: string,
    model: string,
  ): Promise<Summarizer | null> {
    try {
      const response = await fetch(apiUrl(baseUrl, "models"));
      if (!response.ok) {
        return null;
      }
      return new Summarizer(
        createClient(
          Reader,
          createConnectTransport({
            baseUrl: "/rpc",
            interceptors: [withSource("reader")],
          }),
        ),
        baseUrl,
        model,
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

  return fetch(apiUrl(baseUrl, "chat/completions"), {
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
  setSummary: (summary: string) => void,
) => {
  if (!res.ok) {
    throw new Error(`Summary request failed: ${res.status} ${res.statusText}`);
  }

  const reader = res.body?.getReader();
  if (!reader) {
    return "";
  }

  const decoder = new TextDecoder();
  let buffer = "";
  let summary = "";

  const processLine = (line: string) => {
    if (!line.startsWith("data:")) {
      return false;
    }
    const data = line.slice(5).trim();
    if (data === "[DONE]") {
      return true;
    }
    if (!data) {
      return false;
    }

    const event: ChatCompletionChunk = JSON.parse(data);
    const content = event.choices.find((choice) => choice.index === 0)?.delta.content;
    if (content) {
      summary += content;
      setSummary(summary);
    }
    return false;
  };

  try {
    while (true) {
      const { value, done } = await reader.read();
      buffer += decoder.decode(value, { stream: !done });
      let newline: number;
      while ((newline = buffer.indexOf("\n")) !== -1) {
        const line = buffer.slice(0, newline);
        buffer = buffer.slice(newline + 1);
        if (processLine(line)) {
          return summary;
        }
      }
      if (done) {
        processLine(buffer);
        return summary;
      }
    }
  } finally {
    reader.releaseLock();
  }
};
