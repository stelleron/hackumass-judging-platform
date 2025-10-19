import Link from "next/link"

export default function Login() {
  return (
    <div className="flex items-center justify-center h-screen">
      <fieldset className="fieldset bg-base-200 border-base-300 rounded-box w-xs border p-4">
        <legend className="fieldset-legend text-xl font-semibold">
        HackUMass Judging Login
        </legend>

        <label className="label">Username</label>
        <input type="username" className="input" placeholder="Username"/>

        <label className="label">Password</label>
        <input type="password" className="input" placeholder="Password"/>

        <div className="">Don&apos;t have an account? <Link href="/signup" className="text-blue-300">Sign Up</Link></div>

        <button className="btn btn-neutral">Login</button>
      </fieldset>
    </div>
  );
}
