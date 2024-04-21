import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { doesSessionExist, getSession } from "../utils/sessions";
import { useSession } from "../context/session";

export function First() {
  const navigate = useNavigate();
  const { session, setName } = useSession();

  useEffect(() => {
    async function checkSession() {
      const isSessionCreated = await doesSessionExist();
      if (!isSessionCreated) {
        navigate("/");
      }
    }

    async function setSessionName() {
      setName((await getSession()).name);
    }

    if (session.name === "") {
      setSessionName();
    }

    checkSession();
  }, []);
  return <h1>Welcome {session.name} To The First Page</h1>;
}
