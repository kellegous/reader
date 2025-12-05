import "./main.scss";

import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { App } from "./App";
import { ExperimentsProvider } from "./ExperimentsContext";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ExperimentsProvider>
      <App />
    </ExperimentsProvider>
  </StrictMode>
);
