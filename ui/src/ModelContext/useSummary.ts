import hljs from "highlight.js";
import { Marked } from "marked";
import { markedHighlight } from "marked-highlight";
import { useCallback, useState } from "react";
import { useModel } from "./useModel";

const marked = new Marked(
  markedHighlight({
    emptyLangClass: "hljs",
    langPrefix: "hljs language-",
    highlight: (code, lang) => {
      const language = hljs.getLanguage(lang) ? lang : "plaintext";
      return hljs.highlight(code, { language }).value;
    },
  })
);

const renderMarkdown = (text: string) => {
  return marked.parse(text) as string;
};

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
    await summarizer.summarize(entryId, (summary) =>
      setSummary(renderMarkdown(summary))
    );
    setLoading(false);
  }, [summarizer, loading, summary, entryId]);

  return {
    summary,
    loading,
    available: summarizer !== null,
    summarize,
  };
};
