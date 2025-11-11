import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./main.scss";
import { Weekday } from "./time";
import { Weeks } from "./Weeks";
import { ModelProvider } from "./ModelContext";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ModelProvider
      baseUrl="/twirp"
      until={new Date()}
      numWeeks={5}
      weekday={Weekday.Monday}
    >
      <Weeks />
    </ModelProvider>
  </StrictMode>
);
