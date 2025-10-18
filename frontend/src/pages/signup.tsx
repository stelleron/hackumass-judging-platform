export default function Signup() {
  return (
    <div className="flex items-center justify-center h-screen">
        <fieldset className="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
            <legend className="fieldset-legend">Sign Up to ResumeGen</legend>

            <label className="label">Username</label>
            <input type="username" className="input" placeholder="Username"/>

            <label className="label">Password</label>
            <input type="password" className="input" placeholder="Password"/>

            <button className="btn btn-neutral">Sign Up</button>
        </fieldset>
    </div>
  );
}
