import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { doesSessionExist } from "../utils/sessions";

export function Signup() {
  const [email, setEmail] = useState("");
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  async function signup(e: React.MouseEvent<HTMLButtonElement, MouseEvent>) {
    e.preventDefault();
    setLoading(true);

    try {
      let res = await fetch("http://localhost:8080/api/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, name, password }),
      });

      if (res.ok) {
        navigate("/login");
      }
    } catch (e) {
      console.error(e);
    }

    setLoading(false);
  }

  useEffect(() => {
    async function checkSession() {
      const isSessionCreated = await doesSessionExist();
      if (isSessionCreated) {
        navigate("/");
      }
    }

    checkSession();
  }, []);

  return (
    <>
      <form>
        <div>
          <label>Enter your email</label>
          <input
            type="email"
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <label>Enter your name</label>
          <input
            type="text"
            onChange={(e) => setName(e.target.value)}
            required
          />
        </div>
        <div>
          <label>Enter your password</label>
          <input
            type="password"
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <div>
          <button type="submit" onClick={signup}>
            {loading ? "Loading..." : "Submit"}
          </button>
        </div>
      </form>
    </>
  );
}
