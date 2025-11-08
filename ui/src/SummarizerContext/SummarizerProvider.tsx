import { useReaderData } from "../ReaderDataContext";
import { SummarizerContext } from "./SummarizerContext";
import { Summarizer } from "./summarizer";
import { useEffect, useState } from "react";

const defaultOllamaBaseUrl = "http://localhost:11434";
export interface SummarizerProviderProps {
  children: React.ReactNode;
}

export const SummarizerProvider = ({ children }: SummarizerProviderProps) => {
  const [summarizer, setSummarizer] = useState<Summarizer | null>(null);

  const { config } = useReaderData();

  const ollamaBaseUrl = config?.ollamaUrl ?? defaultOllamaBaseUrl;

  useEffect(() => {
    Summarizer.createIfAvailable(ollamaBaseUrl).then(setSummarizer);
  }, [ollamaBaseUrl]);

  return (
    <SummarizerContext.Provider
      value={{ summarizer, available: summarizer !== null }}
    >
      {children}
    </SummarizerContext.Provider>
  );
};
