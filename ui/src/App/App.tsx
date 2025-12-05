import { ModelProvider } from "../ModelContext";
import { Weekday } from "../time";
import { Weeks } from "../Weeks";
import { useSessionRefresh } from "../useSessionRefresh";
import { Header } from "../Header";
import { useExperiments } from "../ExperimentsContext";

export const App = () => {
  useSessionRefresh();

  const { showHeader } = useExperiments();

  return (
    <ModelProvider
      baseUrl="/twirp"
      until={new Date()}
      numWeeks={5}
      weekday={Weekday.Monday}
    >
      {showHeader && <Header />}
      <Weeks />
    </ModelProvider>
  );
};
