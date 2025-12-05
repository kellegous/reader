import { useContext } from "react";
import { ExperimentsContext } from "./ExperimentsContext";

export const useExperiments = () => {
  const context = useContext(ExperimentsContext);
  if (!context) {
    throw new Error("useExperiments must be used within a ExperimentsProvider");
  }
  return context;
};
