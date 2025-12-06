import { ExperimentsContext, ExperimentsState } from "./ExperimentsContext";
import { useState, useEffect } from "react";

const experimentsKey = "expermiments";

export const ExperimentsProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [state, setState] = useState<ExperimentsState>(getState());

  useEffect(() => {
    return bindKey({
      key: "h",
      ctrl: true,
      shift: true,
      action: () =>
        setState((state) =>
          putState({ ...state, showHeader: !state.showHeader })
        ),
    });
  }, []);

  return (
    <ExperimentsContext.Provider value={state}>
      {children}
    </ExperimentsContext.Provider>
  );
};

const emptyState = {
  showHeader: false,
};

const getState = (): ExperimentsState => {
  try {
    const value = localStorage.getItem(experimentsKey);
    if (value) {
      return JSON.parse(value);
    }
  } catch {
    // fall-through
  }
  return emptyState;
};

const putState = (state: ExperimentsState): ExperimentsState => {
  localStorage.setItem(experimentsKey, JSON.stringify(state));
  return state;
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
