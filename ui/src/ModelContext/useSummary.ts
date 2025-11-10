import { useEffect, useState } from "react";
import { useModel } from "./useModel";

export const useSummary = (entryId: bigint) => {
  const [summary, setSummary] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const { model } = useModel();

  useEffect(() => {
    setLoading(true);
    if (!model || !model.canSummarize) {
      setLoading(false);
      return;
    }

    if (loading || summary !== "") {
      return;
    }

    model.summarize(entryId, setSummary).finally(() => setLoading(false));
  }, [model, entryId, loading, summary]);

  return {
    summary,
    loading,
    available: model?.canSummarize ?? false,
    summarize: () => model?.summarize(entryId, setSummary),
  };
};
