import { useCallback, useEffect, useState } from "react";
import { Weekday } from "../time";
import { ModelContext } from "./ModelContext";
import { empty, load, ModelState } from "./model";
import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { Reader } from "../gen/reader_pb";

export interface ModelProviderProps {
  baseUrl: string;
  until: Date;
  numWeeks: number;
  weekday: Weekday;
  children: React.ReactNode;
}

export const ModelProvider = ({
  baseUrl = "/rpc",
  until,
  numWeeks,
  weekday,
  children,
}: ModelProviderProps) => {
  const [model, setModel] = useState<ModelState>(
    empty(createClient(Reader, createConnectTransport({ baseUrl }))),
  );

  const refresh = useCallback(async () => {
    load(model.client, until, numWeeks, weekday, setModel);
  }, [until, numWeeks, weekday, model.client]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  return (
    <ModelContext.Provider value={model}>{children}</ModelContext.Provider>
  );
};
