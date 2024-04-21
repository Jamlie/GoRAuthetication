import { Link, Route, Routes } from "react-router-dom";
import { First } from "./components/First";
import { Second } from "./components/Second";
import { Home } from "./components/Home";
import { SessionProvider } from "./context/session";
import { Signup } from "./components/Signup";
import { Login } from "./components/Login";

function App() {
  return (
    <SessionProvider>
      <nav>
        <ul>
          <li>
            <Link to="/">Home Page</Link>
          </li>
          <li>
            <Link to="/first">First Page</Link>
          </li>
          <li>
            <Link to="/Second">Second Page</Link>
          </li>
        </ul>
      </nav>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/signup" element={<Signup />} />
        <Route path="/login" element={<Login />} />
        <Route path="/first" element={<First />} />
        <Route path="/second" element={<Second />} />
      </Routes>
    </SessionProvider>
  );
}

export default App;
