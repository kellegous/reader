import { useContext } from "react";
import { ReaderDataContext } from "./ReaderDataContext";

export const useReaderData = () => {
  const context = useContext(ReaderDataContext);
  if (!context) {
    throw new Error("useReaderData must be used within a ReaderDataProvider");
  }
  return context;
};
