import { useEffect } from "react";

const defaultEndpoint = "/refresh-session";

export const useSessionRefresh = (endpoint: string = defaultEndpoint) => {
  useEffect(() => {
    fetch(endpoint);
  }, [endpoint]);
};
