import { createContext } from "react";

export interface ExperimentsState {
  showHeader: boolean;
}

export const emptyState = {
  showHeader: true,
};

export const ExperimentsContext = createContext<ExperimentsState>(emptyState);
