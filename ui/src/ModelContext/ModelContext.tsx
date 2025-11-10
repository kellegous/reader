import { createContext } from "react";
import { Model } from "./model";

export interface ModelState {
  model: Model | null;
  loading: boolean;
  refresh: () => Promise<void>;
}

export const ModelContext = createContext<ModelState>({
  model: null,
  loading: false,
  refresh: () => Promise.resolve(),
});
