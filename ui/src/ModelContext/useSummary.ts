import { useCallback, useState } from "react";
import { useModel } from "./useModel";

export const useSummary = (entryId: bigint) => {
  const [summary, setSummary] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const { summarizer } = useModel();

  const summarize = useCallback(async () => {
    if (!summarizer) {
      return;
    }

    if (loading || summary !== "") {
      return;
    }

    setLoading(true);
    await summarizer.summarize(entryId, setSummary);
    setLoading(false);
  }, [summarizer, loading, summary, entryId]);

  return {
    summary,
    loading,
    available: summarizer !== null,
    summarize,
  };
};
