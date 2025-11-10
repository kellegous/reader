import { useCallback, useState } from "react";
import { Weekday } from "../time";
import { Model } from "./model";
import { ModelContext } from "./ModelContext";

export interface ModelProviderProps {
  baseUrl: string;
  until: Date;
  numWeeks: number;
  weekday: Weekday;
  children: React.ReactNode;
}

export const ModelProvider = ({
  baseUrl = "/twirp",
  until,
  numWeeks,
  weekday,
  children,
}: ModelProviderProps) => {
  const [model, setModel] = useState<Model | null>(null);
  const [loading, setLoading] = useState(true);

  const refresh = useCallback(async () => {
    setLoading(true);
    setModel(await Model.load(baseUrl, until, numWeeks, weekday));
    setLoading(false);
  }, [until, numWeeks, weekday, baseUrl]);

  return (
    <ModelContext.Provider value={{ model, loading, refresh }}>
      {children}
    </ModelContext.Provider>
  );
};
