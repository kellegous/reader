import { FetchRPC } from "twirp-ts";
import { ReaderClientJSON } from "../gen/reader.twirp";
import { ReaderDataContext, ReaderDataState } from "./ReaderDataContext";
import { useState, useEffect } from "react";

export const ReaderDataProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [state, setState] = useState<ReaderDataState>({
    me: null,
    entries: [],
    loading: false,
  });

  useEffect(() => {
    const client = new ReaderClientJSON(FetchRPC({ baseUrl: "/twirp" }));
    loadState(client, setState);
  }, []);

  return (
    <ReaderDataContext.Provider value={state}>
      {children}
    </ReaderDataContext.Provider>
  );
};

const loadState = async (
  client: ReaderClientJSON,
  setState: (state: ReaderDataState) => void
) => {
  let state: ReaderDataState = {
    me: null,
    entries: [],
    loading: true,
  };

  setState(state);

  await Promise.all([
    client.GetMe({}).then(({ user }) => {
      state = { ...state, me: user ?? null };
      setState(state);
    }),
  ]);

  state = { ...state, loading: false };
  setState(state);
};
