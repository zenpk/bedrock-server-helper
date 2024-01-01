import {
  createContext,
  Dispatch,
  ReactNode,
  SetStateAction,
  useState,
} from "react";
import { WorldPage } from "@/util/types.ts";

export const WorldsContext = createContext<
  [WorldPage[], Dispatch<SetStateAction<WorldPage[]>>] | null
>(null);

export function WorldsContextProvider({ children }: { children: ReactNode }) {
  const defaultValue = useState<WorldPage[]>([]);
  return (
    <WorldsContext.Provider value={defaultValue}>
      {children}
    </WorldsContext.Provider>
  );
}
