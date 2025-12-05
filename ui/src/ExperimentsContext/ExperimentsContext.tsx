import { createContext } from "react";

export interface ExperimentsState {
  showHeader: boolean;
}

export const ExperimentsContext = createContext<ExperimentsState>({
  showHeader: false,
});
