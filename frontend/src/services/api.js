const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api";

async function request(path, options = {}) {
  const response = await fetch(`${API_URL}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    ...options,
  });

  const data = await response.json().catch(() => ({}));

  if (!response.ok) {
    throw new Error(data.error || "Something went wrong");
  }

  return data;
}

export function signup(name, email, password) {
  return request("/auth/signup", {
    method: "POST",
    body: JSON.stringify({ name, email, password }),
  });
}

export function login(email, password) {
  return request("/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

export function createPoll(question, options, token) {
  return request("/polls", {
    method: "POST",
    headers: { Authorization: `Bearer ${token}` },
    body: JSON.stringify({ question, options }),
  });
}

export function getMyPolls(token) {
  return request("/polls/my", {
    headers: { Authorization: `Bearer ${token}` },
  });
}

export function getPoll(id) {
  return request(`/polls/${id}`);
}

export function getResults(id) {
  return request(`/polls/${id}/results`);
}

export function vote(id, optionId, voterId) {
  return request(`/polls/${id}/vote`, {
    method: "POST",
    body: JSON.stringify({ optionId, voterId }),
  });
}

export function websocketUrl(id) {
  const httpUrl = API_URL.replace(/\/api$/, "");
  return httpUrl.replace(/^http/, "ws") + `/api/polls/${id}/live`;
}
