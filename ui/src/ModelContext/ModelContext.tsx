import { createContext } from "react";
import { empty, ModelState } from "./model";
import { ReaderClientJSON } from "../gen/reader.twirp";
import { FetchRPC } from "twirp-ts";

export const ModelContext = createContext<ModelState>(
  empty(new ReaderClientJSON(FetchRPC({ baseUrl: "/twirp" })))
);
