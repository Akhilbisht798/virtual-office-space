import { useState } from "react";
import { NavLink } from "react-router";
import axios from "axios";
import { SERVER } from "../../global";

const SignUp = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault(); 
    const signUpUrl = `${SERVER}/api/v1/signup`;

    try {
      const signUpData = {
        username,
        password,
      };
      const res = await axios.post(signUpUrl, signUpData);
      console.log("request succesfull.", res)

      setError("");
    } catch (err) {
      console.error("SignUp failed:", err);
      setError(err.message);
    }
  };
  return (
    <div className="flex justify-center items-center h-screen bg-gray-100">
      <div className="bg-white p-6 rounded shadow-md w-full max-w-sm">
        <h1 className="text-xl font-semibold mb-4">SignUp</h1>
        {error && <p className="text-red-500 text-sm mb-4">{error}</p>}
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label
              htmlFor="username"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Username
            </label>
            <input
              type="text"
              id="username"
              className="w-full border border-gray-300 rounded p-2 focus:outline-none focus:ring focus:ring-blue-300"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </div>
          <div className="mb-4">
            <label
              htmlFor="password"
              className="block text-sm font-medium text-gray-700 mb-1"
            >
              Password
            </label>
            <input
              type="password"
              id="password"
              className="w-full border border-gray-300 rounded p-2 focus:outline-none focus:ring focus:ring-blue-300"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <button
            type="submit"
            className="w-full bg-blue-500 text-white p-2 rounded hover:bg-blue-600"
          >
            SignUp
          </button>
        </form>
        <NavLink to="/login" end>Already have a account. Login</NavLink>
      </div>
    </div>
  );
}

export default SignUp;
