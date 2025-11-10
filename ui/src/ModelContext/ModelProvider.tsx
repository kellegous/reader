import { useCallback, useEffect, useState } from "react";
import { Weekday } from "../time";
import { ModelContext } from "./ModelContext";
import { empty, load, ModelState } from "./model";
import { ReaderClientJSON } from "../gen/reader.twirp";
import { FetchRPC } from "twirp-ts";

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
  const [model, setModel] = useState<ModelState>(
    empty(new ReaderClientJSON(FetchRPC({ baseUrl })))
  );

  const refresh = useCallback(async () => {
    load(baseUrl, until, numWeeks, weekday, setModel);
  }, [until, numWeeks, weekday, baseUrl]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  return (
    <ModelContext.Provider value={model}>{children}</ModelContext.Provider>
  );
};
