import { createContext } from "react";
import { empty, ModelState } from "./model";
import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { withSource } from "../lib/withSource";
import { Reader } from "../gen/reader_pb";

export const ModelContext = createContext<ModelState>(
  empty(
    createClient(
      Reader,
      createConnectTransport({
        baseUrl: "/rpc",
        interceptors: [withSource("reader")],
      }),
    ),
  ),
);
