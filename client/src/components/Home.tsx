import { useEffect, useState } from "react";
import { doesSessionExist, getSession } from "../utils/sessions";
import { useSession } from "../context/session";
import { useNavigate } from "react-router-dom";

export function Home() {
  const [loading, setLoading] = useState(false);
  const { session, setName } = useSession();
  const navigate = useNavigate();

  async function setSessionName() {
    setName((await getSession()).name);
    setLoading(false);
  }

  async function handleAddClick() {
    try {
      const response = await fetch("http://localhost:8080/add-session", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (!response.ok) {
        console.log("Failed");
      } else {
        await setSessionName();
      }
    } catch (e) {
      console.error(e);
    }
  }

  useEffect(() => {
    if (session.name === "" || session.name === null) {
      setLoading(true);
      setSessionName();
    }

    async function checkSession() {
      const isSessionCreated = await doesSessionExist();
      if (!isSessionCreated) {
        navigate("/login");
      }
    }

    checkSession();
  }, []);

  return (
    <>
      {loading ? (
        <h1>Loading</h1>
      ) : (
        <h1>
          Welcome{" "}
          {session.name === "" || session.name === null ? (
            <span>User</span>
          ) : (
            session.name
          )}{" "}
          To The Home Page
        </h1>
      )}

      <button type="button" onClick={handleAddClick}>
        Add
      </button>
    </>
  );
}
