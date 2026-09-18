import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { createPoll } from "../services/api";

export default function CreatePoll() {
  const navigate = useNavigate();
  const [question, setQuestion] = useState("");
  const [options, setOptions] = useState(["", ""]);
  const [error, setError] = useState("");

  function updateOption(index, value) {
    const copy = [...options];
    copy[index] = value;
    setOptions(copy);
  }

  function addOption() {
    if (options.length < 6) setOptions([...options, ""]);
  }

  function removeOption(index) {
    if (options.length <= 2) return;
    setOptions(options.filter((_, i) => i !== index));
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setError("");

    const cleaned = options.map(x => x.trim()).filter(Boolean);
    if (question.trim().length < 3) {
      setError("Please enter a question.");
      return;
    }
    if (cleaned.length < 2) {
      setError("Add at least two options.");
      return;
    }

    try {
      const data = await createPoll(
        question,
        cleaned,
        localStorage.getItem("token")
      );
      navigate(`/poll/${data.id}`);
    } catch (err) {
      setError(err.message);
    }
  }

  return (
    <div className="form-card wide">
      <span className="badge">NEW POLL</span>
      <h2>Ask your audience</h2>
      <p className="muted">Keep it simple. You can share the generated link after creating it.</p>

      <form onSubmit={handleSubmit}>
        <label>
          Question
          <input
            value={question}
            onChange={e => setQuestion(e.target.value)}
            placeholder="What should we have for lunch?"
            maxLength="200"
            required
          />
        </label>

        <div className="options-title">Options</div>
        {options.map((option, index) => (
          <div className="option-input" key={index}>
            <input
              value={option}
              onChange={e => updateOption(index, e.target.value)}
              placeholder={`Option ${index + 1}`}
              maxLength="100"
              required
            />
            {options.length > 2 && (
              <button type="button" className="remove" onClick={() => removeOption(index)}>×</button>
            )}
          </div>
        ))}

        {options.length < 6 && (
          <button type="button" className="secondary small" onClick={addOption}>+ Add option</button>
        )}

        {error && <div className="error">{error}</div>}
        <button type="submit">Create poll</button>
      </form>
    </div>
  );
}
