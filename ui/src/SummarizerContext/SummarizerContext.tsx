import { Summarizer } from "./summarizer";
import { createContext } from "react";

export interface SummarizerState {
  summarizer: Summarizer | null;
  available: boolean;
}

export const SummarizerContext = createContext<SummarizerState | null>(null);
