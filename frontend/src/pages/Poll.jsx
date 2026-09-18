import React, { useEffect, useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import { getPoll, getResults, vote, websocketUrl } from "../services/api";

function getVoterId() {
  let id = localStorage.getItem("poll_voter_id");
  if (!id) {
    id = crypto.randomUUID();
    localStorage.setItem("poll_voter_id", id);
  }
  return id;
}

export default function Poll() {
  const { id } = useParams();
  const [poll, setPoll] = useState(null);
  const [results, setResults] = useState([]);
  const [selected, setSelected] = useState("");
  const [voted, setVoted] = useState(false);
  const [error, setError] = useState("");
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    Promise.all([getPoll(id), getResults(id)])
      .then(([pollData, resultData]) => {
        setPoll(pollData);
        setResults(resultData.results);
      })
      .catch(err => setError(err.message));
  }, [id]);

  useEffect(() => {
    if (!id) return;
    const socket = new WebSocket(websocketUrl(id));

    socket.onmessage = event => {
      const data = JSON.parse(event.data);
      setResults(data.results || []);
    };

    socket.onerror = () => {
      console.log("Live connection error");
    };

    return () => socket.close();
  }, [id]);

  const totalVotes = useMemo(
    () => results.reduce((sum, item) => sum + Number(item.votes), 0),
    [results]
  );

  async function handleVote(e) {
    e.preventDefault();
    if (!selected) {
      setError("Choose an option first.");
      return;
    }

    setError("");
    try {
      const data = await vote(id, selected, getVoterId());
      setResults(data.results);
      setVoted(true);
    } catch (err) {
      setError(err.message);
    }
  }

  async function copyLink() {
    await navigator.clipboard.writeText(window.location.href);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  }

  if (error && !poll) return <div className="error">{error}</div>;
  if (!poll) return <div className="loading">Loading poll...</div>;

  return (
    <section className="poll-page">
      <div className="poll-header">
        <span className="live-label"><span /> LIVE</span>
        <button className="secondary small" onClick={copyLink}>
          {copied ? "Copied!" : "Copy share link"}
        </button>
      </div>

      <h2>{poll.question}</h2>
      <p className="muted">{totalVotes} total vote{totalVotes === 1 ? "" : "s"} · results update live</p>

      {!voted && !poll.isClosed && (
        <form className="vote-box" onSubmit={handleVote}>
          {poll.options.map(option => (
            <label className={`choice ${selected === option.id ? "selected" : ""}`} key={option.id}>
              <input
                type="radio"
                name="option"
                value={option.id}
                checked={selected === option.id}
                onChange={e => setSelected(e.target.value)}
              />
              <span>{option.text}</span>
            </label>
          ))}
          <button type="submit">Vote</button>
        </form>
      )}

      {poll.isClosed && <div className="notice">This poll is closed.</div>}
      {voted && <div className="success">Vote recorded. Watch the results update live.</div>}
      {error && <div className="error">{error}</div>}

      <div className="results-card">
        <div className="results-title">
          <h3>Live results</h3>
          <span>{totalVotes} votes</span>
        </div>

        {results.map(item => {
          const percent = totalVotes ? Math.round((Number(item.votes) / totalVotes) * 100) : 0;
          return (
            <div className="result-row" key={item.optionId}>
              <div className="result-label">
                <span>{item.text}</span>
                <strong>{percent}% · {item.votes}</strong>
              </div>
              <div className="bar"><i style={{ width: `${percent}%` }} /></div>
            </div>
          );
        })}
      </div>
    </section>
  );
}
