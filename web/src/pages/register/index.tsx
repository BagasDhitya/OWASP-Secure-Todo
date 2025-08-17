// src/pages/Register.tsx
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import API from "../../api/axios";

interface RegisterForm {
    username: string;
    email: string;
    password: string;
}

export default function Register() {
    const [form, setForm] = useState<RegisterForm>({
        username: "",
        email: "",
        password: "",
    });
    const [errors, setErrors] = useState<string | null>(null);
    const navigate = useNavigate();

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setForm({ ...form, [e.target.name]: e.target.value });
        setErrors(null);
    };

    const validateForm = () => {
        const { username, email, password } = form;
        const alphanumRegex = /^[a-zA-Z0-9]+$/;

        if (!username) {
            return "Username is required";
        }
        if (username.length < 3 || username.length > 50) {
            return "Username must be between 3 and 50 characters";
        }
        if (!alphanumRegex.test(username)) {
            return "Username must only contain letters and numbers";
        }

        if (!email) {
            return "Email is required";
        }
        if (!/\S+@\S+\.\S+/.test(email)) {
            return "Invalid email address";
        }

        if (!password) {
            return "Password is required";
        }
        if (password.length < 8 || password.length > 72) {
            return "Password must be between 8 and 72 characters";
        }

        return null;
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        const validationError = validateForm();
        if (validationError) {
            setErrors(validationError);
            return;
        }

        try {
            const res = await API.post("auth/register", form, {
                headers: { "Content-Type": "application/json" },
            });
            if (res.status === 201) {
                navigate("/");
            }
        } catch (err: any) {
            if (err.response?.data?.error) {
                setErrors(err.response.data.error);
            } else {
                setErrors("Something went wrong");
            }
        }
    };

    return (
        <div className="min-h-screen w-screen flex items-center justify-center bg-gray-100 p-4">
            <div className="bg-white p-8 rounded shadow-md w-full max-w-md">
                <h2 className="text-2xl font-bold mb-6 text-center">Register</h2>
                {errors && (
                    <div className="mb-4 text-red-600 font-medium text-sm">{errors}</div>
                )}
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label className="block mb-1 font-medium">Username</label>
                        <input
                            type="text"
                            name="username"
                            value={form.username}
                            onChange={handleChange}
                            className="w-full border px-3 py-2 rounded focus:outline-none focus:ring focus:border-blue-300"
                            placeholder="Enter username"
                            required
                        />
                    </div>
                    <div>
                        <label className="block mb-1 font-medium">Email</label>
                        <input
                            type="email"
                            name="email"
                            value={form.email}
                            onChange={handleChange}
                            className="w-full border px-3 py-2 rounded focus:outline-none focus:ring focus:border-blue-300"
                            placeholder="Enter email"
                            required
                        />
                    </div>
                    <div>
                        <label className="block mb-1 font-medium">Password</label>
                        <input
                            type="password"
                            name="password"
                            value={form.password}
                            onChange={handleChange}
                            className="w-full border px-3 py-2 rounded focus:outline-none focus:ring focus:border-blue-300"
                            placeholder="Enter password"
                            required
                        />
                    </div>
                    <button
                        type="submit"
                        className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition"
                    >
                        Register
                    </button>
                </form>
            </div>
        </div>
    );
}
