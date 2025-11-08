import { useCallback, useState } from "react";
import { useSummarizer } from "./useSummarizer";

export const useSummary = (entryId: bigint) => {
  const [summary, setSummary] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const { summarizer, available } = useSummarizer();

  const summarize = useCallback(() => {
    if (!available || !summarizer) {
      return;
    }

    setLoading(true);

    summarizer.summarize(entryId, setSummary).finally(() => setLoading(false));
  }, [available, summarizer, entryId]);

  return { summary, loading, available, summarize };
};
