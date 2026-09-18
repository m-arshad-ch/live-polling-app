import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { login } from "../services/api";

export default function Login({ onLogin }) {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function handleSubmit(e) {
    e.preventDefault();
    setError("");
    try {
      const data = await login(email, password);
      onLogin(data);
      navigate("/dashboard");
    } catch (err) {
      setError(err.message);
    }
  }

  return (
    <div className="form-card">
      <span className="badge">WELCOME BACK</span>
      <h2>Log in</h2>
      <p className="muted">Create and manage your polls.</p>
      <form onSubmit={handleSubmit}>
        <label>Email<input type="email" value={email} onChange={e => setEmail(e.target.value)} required /></label>
        <label>Password<input type="password" value={password} onChange={e => setPassword(e.target.value)} required /></label>
        {error && <div className="error">{error}</div>}
        <button type="submit">Log in</button>
      </form>
      <p className="form-footer">No account? <Link to="/signup">Sign up</Link></p>
    </div>
  );
}
