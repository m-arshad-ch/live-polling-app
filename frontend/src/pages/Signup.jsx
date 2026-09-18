import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { signup } from "../services/api";

export default function Signup({ onLogin }) {
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  async function handleSubmit(e) {
    e.preventDefault();
    setError("");
    try {
      const data = await signup(name, email, password);
      onLogin(data);
      navigate("/dashboard");
    } catch (err) {
      setError(err.message);
    }
  }

  return (
    <div className="form-card">
      <span className="badge">GET STARTED</span>
      <h2>Create your account</h2>
      <p className="muted">You need an account to create or manage polls.</p>
      <form onSubmit={handleSubmit}>
        <label>Name<input value={name} onChange={e => setName(e.target.value)} required /></label>
        <label>Email<input type="email" value={email} onChange={e => setEmail(e.target.value)} required /></label>
        <label>Password<input type="password" value={password} onChange={e => setPassword(e.target.value)} minLength="6" required /></label>
        {error && <div className="error">{error}</div>}
        <button type="submit">Create account</button>
      </form>
      <p className="form-footer">Already registered? <Link to="/login">Log in</Link></p>
    </div>
  );
}
