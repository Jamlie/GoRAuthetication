import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { doesSessionExist } from "../utils/sessions";

export function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  async function login(e: React.MouseEvent<HTMLButtonElement, MouseEvent>) {
    e.preventDefault();
    setLoading(true);

    try {
      let res = await fetch("http://localhost:8080/api/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      });

      if (res.ok) {
        navigate("/");
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
          <label>Enter your password</label>
          <input
            type="password"
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <div>
          <button type="submit" onClick={login}>
            {loading ? "Loading..." : "Submit"}
          </button>
        </div>
      </form>
    </>
  );
}
