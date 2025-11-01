import { SummarizerContext } from "./SummarizerContext";
import { Summarizer } from "./summarizer";
import { useEffect, useState } from "react";

export interface SummarizerProviderProps {
  ollamaBaseUrl: string;
  children: React.ReactNode;
}

export const SummarizerProvider = ({
  ollamaBaseUrl,
  children,
}: SummarizerProviderProps) => {
  const [summarizer, setSummarizer] = useState<Summarizer | null>(null);

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
