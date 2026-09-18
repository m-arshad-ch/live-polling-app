import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { getMyPolls } from "../services/api";

export default function Dashboard() {
  const [polls, setPolls] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    getMyPolls(localStorage.getItem("token"))
      .then(setPolls)
      .catch(err => setError(err.message));
  }, []);

  return (
    <section>
      <div className="page-heading">
        <div>
          <span className="badge">DASHBOARD</span>
          <h2>My polls</h2>
          <p className="muted">Create a poll and send its link to your audience.</p>
        </div>
        <Link to="/create" className="button-link">+ Create poll</Link>
      </div>

      {error && <div className="error">{error}</div>}

      <div className="poll-list">
        {polls.length === 0 ? (
          <div className="empty">
            <h3>No polls yet</h3>
            <p>Create your first poll to get started.</p>
            <Link to="/create" className="button-link">Create poll</Link>
          </div>
        ) : polls.map(poll => (
          <Link className="poll-item" to={`/poll/${poll.id}`} key={poll.id}>
            <div>
              <strong>{poll.question}</strong>
              <small>{poll.options.length} options · {poll.isClosed ? "Closed" : "Live"}</small>
            </div>
            <span>→</span>
          </Link>
        ))}
      </div>
    </section>
  );
}
