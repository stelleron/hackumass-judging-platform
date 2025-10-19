import { useState } from "react";

export default function Signup() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleSignup = async () => {
    const res = await fetch("http://localhost:8000/api/signup", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include", // important for cookies
      body: JSON.stringify({ username, password }),
    });

    const data = await res.json();
    alert(data.message);
  };


  return (
    <div className="flex items-center justify-center h-screen">
        <fieldset className="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
            <legend className="fieldset-legend">Sign Up to ResumeGen</legend>

            <label className="label">Username</label>
            <input
              type="username"
              className="input"
              placeholder="Username"
              value={username}
              onChange={e => setUsername(e.target.value)}
            />

            <label className="label">Password</label>
            <input
              type="password"
              className="input"
              placeholder="Password"
              value={password}
              onChange={e => setPassword(e.target.value)}
            />

            <button onClick={handleSignup} className="btn btn-neutral">Sign Up</button>
        </fieldset>
    </div>
  );
}
