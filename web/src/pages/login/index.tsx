// src/pages/Login.tsx
import { useState } from "react";
import API from "../../api/axios";
import { useNavigate } from "react-router-dom";

export default function Login() {
    const [form, setForm] = useState({ email: "", password: "" });
    const [message, setMessage] = useState("");
    const navigate = useNavigate();

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setForm({ ...form, [e.target.name]: e.target.value });
        setMessage("");
    };

    const validateForm = () => {
        if (!form.email) {
            return "❌ Email is required.";
        }
        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
            return "❌ Please enter a valid email.";
        }
        if (!form.password) {
            return "❌ Password is required.";
        }
        return null;
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        const error = validateForm();
        if (error) {
            setMessage(error);
            return;
        }

        try {
            const res = await API.post("/auth/login", form, {
                headers: { "Content-Type": "application/json" }
            });
            // simpan access token / refresh token sesuai backend
            localStorage.setItem("token", res.data.token);
            navigate("/todos");
        } catch (err: any) {
            console.error(err);
            setMessage(err.response?.data?.error || "Invalid credentials");
        }
    };

    return (
        <div className="flex items-center justify-center w-screen min-h-screen bg-gradient-to-br from-green-400 via-green-500 to-green-700 p-4">
            <div className="bg-white shadow-lg rounded-2xl p-8 w-full max-w-sm">
                <h1 className="text-2xl font-bold text-center mb-2 text-gray-800">Welcome Back</h1>
                <p className="text-center text-gray-500 mb-6 text-sm">Please sign in to your account</p>

                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
                            Email
                        </label>
                        <input
                            id="email"
                            className="border border-gray-300 rounded-lg px-3 py-2 w-full focus:outline-none focus:ring-2 focus:ring-green-400"
                            type="email"
                            name="email"
                            placeholder="you@example.com"
                            value={form.email}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div>
                        <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1">
                            Password
                        </label>
                        <input
                            id="password"
                            className="border border-gray-300 rounded-lg px-3 py-2 w-full focus:outline-none focus:ring-2 focus:ring-green-400"
                            type="password"
                            name="password"
                            placeholder="••••••••"
                            value={form.password}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    {message && (
                        <p className="text-red-500 text-sm mt-2">{message}</p>
                    )}

                    <button
                        type="submit"
                        className="bg-green-500 hover:bg-green-600 text-white w-full py-2 rounded-lg font-semibold shadow transition-colors"
                    >
                        Login
                    </button>
                </form>

                <div className="mt-6 text-center">
                    <p className="text-sm text-gray-600">
                        Don't have an account?{" "}
                        <a href="/register" className="text-green-600 hover:underline font-medium">
                            Sign Up
                        </a>
                    </p>
                </div>
            </div>
        </div>
    );
}
