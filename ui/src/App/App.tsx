import { ModelProvider } from "../ModelContext";
import { Weekday } from "../time";
import { Weeks } from "../Weeks";
import { useSessionRefresh } from "../useSessionRefresh";

export const App = () => {
  useSessionRefresh();

  return (
    <ModelProvider
      baseUrl="/twirp"
      until={new Date()}
      numWeeks={5}
      weekday={Weekday.Monday}
    >
      <Weeks />
    </ModelProvider>
  );
};
