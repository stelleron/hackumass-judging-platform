import Link from "next/link"
import { useState } from "react";

export default function Login() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = async () => {
    const res = await fetch("http://localhost:8000/api/login", {
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
        <legend className="fieldset-legend text-xl font-semibold">
        HackUMass Judging Login
        </legend>

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

        <div className="">Don&apos;t have an account? <Link href="/signup" className="text-blue-300">Sign Up</Link></div>

        <button onClick={handleLogin} className="btn btn-neutral">Login</button>
      </fieldset>
    </div>
  );
}
