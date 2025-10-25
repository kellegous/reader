import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./main.scss";
import App from "./App.tsx";
import { ReaderDataProvider } from "./ReaderDataContext";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ReaderDataProvider>
      <App />
    </ReaderDataProvider>
  </StrictMode>
);
