import { ExperimentsContext, ExperimentsState } from "./ExperimentsContext";
import { useState, useEffect } from "react";

export const ExperimentsProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [state, setState] = useState<ExperimentsState>({
    showHeader: false,
  });

  useEffect(() => {
    return bindKey({
      key: "h",
      ctrl: true,
      shift: true,
      action: () =>
        setState((state) => ({ ...state, showHeader: !state.showHeader })),
    });
  }, []);

  return (
    <ExperimentsContext.Provider value={state}>
      {children}
    </ExperimentsContext.Provider>
  );
};

interface KeyBinding {
  key: string;
  ctrl?: boolean;
  shift?: boolean;
  action: () => void;
}

const bindKey = ({ key, ctrl = false, shift = false, action }: KeyBinding) => {
  const onKeydown = (event: KeyboardEvent) => {
    if (
      event.key.toLowerCase() === key.toLowerCase() &&
      event.ctrlKey === ctrl &&
      event.shiftKey === shift
    ) {
      action();
    }
  };
  window.addEventListener("keydown", onKeydown);
  return () => {
    window.removeEventListener("keydown", onKeydown);
  };
};
