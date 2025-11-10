import { useEffect, useState } from "react";
import { useModel } from "./useModel";

export const useSummary = (entryId: bigint) => {
  const [summary, setSummary] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const { summarizer } = useModel();

  useEffect(() => {
    if (!summarizer) {
      return;
    }

    if (loading || summary !== "") {
      return;
    }

    setLoading(true);
    summarizer.summarize(entryId, setSummary).finally(() => setLoading(false));
  }, [summarizer, entryId, loading, summary]);

  return {
    summary,
    loading,
    available: summarizer !== null,
    summarize: () =>
      summarizer?.summarize(entryId, setSummary) ?? Promise.resolve(""),
  };
};
