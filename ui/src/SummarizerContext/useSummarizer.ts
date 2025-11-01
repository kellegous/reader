import { useContext } from "react";
import { SummarizerContext } from "./SummarizerContext";

export const useSummarizer = () => {
  const context = useContext(SummarizerContext);
  if (!context) {
    throw new Error("useSummarizer must be used within a SummarizerProvider");
  }
  return context.summarizer;
};
