import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { doesSessionExist } from "../utils/sessions";

export function Second() {
  const navigate = useNavigate();
  useEffect(() => {
    async function checkSession() {
      const isSessionCreated = await doesSessionExist();
      if (isSessionCreated) {
        navigate("/");
      }
    }

    checkSession();
  }, []);
  return <h1>Welcome To The Second Page</h1>;
}
