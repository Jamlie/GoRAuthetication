import { createContext, useContext, useState } from "react";

type Session = {
  name: string;
};

type SessionContextType = {
  session: Session;
  setName: (name: string) => void;
};

const SessionContext = createContext<SessionContextType | undefined>(undefined);

export function useSession() {
  const context = useContext(SessionContext);
  if (!context) {
    throw new Error("useSession must be used within a SessionProvider");
  }
  return context;
}

export function SessionProvider({ children }: { children: React.ReactNode }) {
  const [session, setSession] = useState<Session>({ name: "" });

  function setName(name: string) {
    setSession({ name: name });
  }

  return (
    <SessionContext.Provider value={{ session, setName }}>
      {children}
    </SessionContext.Provider>
  );
}
