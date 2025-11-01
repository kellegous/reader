import { useEffect, useState } from "react";
import { Summarizer } from "./summarizer";

export const useSummary = (summarizer: Summarizer, entryId: bigint) => {
  const [summary, setSummary] = useState<string>("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true);
    summarizer
      .summarize(entryId)
      .then(setSummary)
      .finally(() => setLoading(false));
  }, [summarizer, entryId]);

  return { summary, loading };
};
