import React, { useState } from "react";
import { Link, Navigate, Route, Routes, useNavigate } from "react-router-dom";
import Login from "./pages/Login";
import Signup from "./pages/Signup";
import Dashboard from "./pages/Dashboard";
import CreatePoll from "./pages/CreatePoll";
import Poll from "./pages/Poll";

function getUser() {
  return JSON.parse(localStorage.getItem("user") || "null");
}

export default function App() {
  const [user, setUser] = useState(getUser());

  function saveAuth(data) {
    localStorage.setItem("token", data.token);
    localStorage.setItem("user", JSON.stringify(data.user));
    setUser(data.user);
  }

  function logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    setUser(null);
  }

  return (
      <div className="app">
        <header className="navbar">
          <Link to="/" className="brand">PulsePoll</Link>

          <nav>
            {user ? (
              <>
                <Link to="/dashboard">My Polls</Link>
                <Link to="/create">Create</Link>
                <button className="link-button" onClick={logout}>Logout</button>
              </>
            ) : (
              <>
                <Link to="/login">Login</Link>
                <Link to="/signup" className="nav-cta">Sign up</Link>
              </>
            )}
          </nav>
        </header>

        <main className="container">
          <Routes>
            <Route path="/" element={<Home user={user} />} />
            <Route path="/login" element={<Login onLogin={saveAuth} />} />
            <Route path="/signup" element={<Signup onLogin={saveAuth} />} />

            <Route
              path="/dashboard"
              element={user ? <Dashboard /> : <Navigate to="/login" />}
            />

            <Route
              path="/create"
              element={user ? <CreatePoll /> : <Navigate to="/login" />}
            />

            <Route path="/poll/:id" element={<Poll />} />
          </Routes>
        </main>
      </div>
  );
}

function Home({ user }) {
  const navigate = useNavigate();

  return (
    <section className="hero">
      <div>
        <span className="badge">LIVE POLLING</span>

        <h1>
          Ask a question.<br />
          See the vote live.
        </h1>

        <p>
          Create a poll, share the link, and watch the results change instantly
          as people vote.
        </p>

        <div className="hero-actions">
          <button onClick={() => navigate(user ? "/create" : "/signup")}>
            Create a poll
          </button>

          <button
            className="secondary"
            onClick={() => navigate(user ? "/dashboard" : "/login")}
          >
            {user ? "View my polls" : "Log in"}
          </button>
        </div>
      </div>

      <div className="hero-card">
        <div className="live-dot">
          <span /> Live results
        </div>

        <div className="fake-question">
          Which feature should we build next?
        </div>

        <div className="fake-row">
          <span>Dark mode</span>
          <b style={{ width: "72%" }} />
        </div>

        <div className="fake-row">
          <span>Poll analytics</span>
          <b style={{ width: "52%" }} />
        </div>

        <div className="fake-row">
          <span>Team sharing</span>
          <b style={{ width: "35%" }} />
        </div>

        <small>Updates without refreshing</small>
      </div>
    </section>
  );
}
