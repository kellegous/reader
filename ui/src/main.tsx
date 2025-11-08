import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./main.scss";
import { ReaderDataProvider } from "./ReaderDataContext";
import { Weekday } from "./time";
import { Weeks } from "./Weeks";
import { SummarizerProvider } from "./SummarizerContext";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ReaderDataProvider
      until={new Date()}
      numWeeks={5}
      weekday={Weekday.Monday}
    >
      <SummarizerProvider>
        <Weeks />
      </SummarizerProvider>
    </ReaderDataProvider>
  </StrictMode>
);
