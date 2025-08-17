import { useEffect, useState } from "react";
import API from "../../api/axios";

interface Todo {
    id: number;
    title: string;
    description: string;
    status: "pending" | "completed";
}

export default function TodoList() {
    const [todos, setTodos] = useState<Todo[]>([]);
    const [form, setForm] = useState({ title: "", description: "", status: "pending" as "pending" | "completed" });
    const [editId, setEditId] = useState<number | null>(null);
    const [error, setError] = useState("");

    const token = localStorage.getItem("token");

    const fetchTodos = async () => {
        const res = await API.get("/tasks", {
            headers: { Authorization: `Bearer ${token}` },
        });
        setTodos(Array.isArray(res.data) ? res.data : []);
    };

    useEffect(() => {
        fetchTodos();
    }, []);

    const validateInput = (): boolean => {
        setError("");
        if (!form.title.trim()) {
            setError("Title is required.");
            return false;
        }
        if (form.title.length > 255) {
            setError("Title must not exceed 255 characters.");
            return false;
        }
        if (form.description.length > 5000) {
            setError("Description must not exceed 5000 characters.");
            return false;
        }
        return true;
    };

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
        setForm({ ...form, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!validateInput()) return;

        if (editId) {
            // Update task
            await API.put(
                `/tasks/${editId}`,
                form,
                { headers: { Authorization: `Bearer ${token}` } }
            );
            setEditId(null);
        } else {
            // Add task
            await API.post(
                "/tasks",
                form,
                { headers: { Authorization: `Bearer ${token}` } }
            );
        }

        setForm({ title: "", description: "", status: "pending" });
        fetchTodos();
    };

    const startEdit = (todo: Todo) => {
        setEditId(todo.id);
        setForm({
            title: todo.title,
            description: todo.description,
            status: todo.status,
        });
    };

    const deleteTodo = async (id: number) => {
        await API.delete(`/tasks/${id}`, {
            headers: { Authorization: `Bearer ${token}` },
        });
        fetchTodos();
    };

    return (
        <div className="min-h-screen w-screen bg-gray-100 flex items-center justify-center p-6">
            <div className="bg-white shadow-xl rounded-2xl w-full max-w-2xl p-6">
                <h2 className="text-2xl font-bold text-center text-gray-800 mb-6">
                    ✅ My Task Manager
                </h2>

                {/* Form Input */}
                <form onSubmit={handleSubmit} className="space-y-3 mb-6">
                    {error && (
                        <p className="bg-red-100 text-red-700 px-3 py-2 rounded-lg text-sm">
                            {error}
                        </p>
                    )}
                    <input
                        name="title"
                        className="border border-gray-300 rounded-lg p-3 w-full focus:outline-none focus:ring-2 focus:ring-blue-400"
                        value={form.title}
                        onChange={handleChange}
                        placeholder="Task title (required, max 255)"
                    />
                    <textarea
                        name="description"
                        className="border border-gray-300 rounded-lg p-3 w-full focus:outline-none focus:ring-2 focus:ring-blue-400"
                        value={form.description}
                        onChange={handleChange}
                        placeholder="Task description (optional, max 5000)"
                        rows={3}
                    />
                    <select
                        name="status"
                        className="border border-gray-300 rounded-lg p-3 w-full focus:outline-none focus:ring-2 focus:ring-blue-400"
                        value={form.status}
                        onChange={handleChange}
                    >
                        <option value="pending">⏳ Pending</option>
                        <option value="completed">✅ Completed</option>
                    </select>
                    <button
                        type="submit"
                        className="bg-blue-500 text-white w-full py-3 rounded-lg hover:bg-blue-600 transition"
                    >
                        {editId ? "Update Task" : "Add Task"}
                    </button>
                </form>

                {/* Task List */}
                <ul className="space-y-4">
                    {todos.length === 0 && (
                        <p className="text-gray-500 text-center">No tasks yet. 🎉</p>
                    )}
                    {todos.map((todo) => (
                        <li
                            key={todo.id}
                            className="bg-gray-50 p-4 rounded-lg shadow-sm hover:shadow-md transition"
                        >
                            <div className="flex justify-between items-center mb-2">
                                <h2
                                    className={`text-lg font-semibold ${todo.status === "completed"
                                        ? "line-through text-gray-400"
                                        : "text-gray-800"
                                        }`}
                                >
                                    {todo.title}
                                </h2>
                                <div className="flex gap-2">
                                    <button
                                        onClick={() => startEdit(todo)}
                                        className="px-3 py-1 rounded-lg text-sm font-medium bg-blue-100 text-blue-700 hover:bg-blue-200"
                                    >
                                        ✏ Edit
                                    </button>
                                    <button
                                        onClick={() => deleteTodo(todo.id)}
                                        className="px-3 py-1 rounded-lg text-sm font-medium bg-red-100 text-red-700 hover:bg-red-200"
                                    >
                                        ✖ Delete
                                    </button>
                                </div>
                            </div>
                            {todo.description && (
                                <p className="text-gray-600 mb-2">{todo.description}</p>
                            )}
                            <span
                                className={`inline-block px-2 py-1 rounded text-xs ${todo.status === "pending"
                                    ? "bg-yellow-100 text-yellow-700"
                                    : "bg-green-100 text-green-700"
                                    }`}
                            >
                                {todo.status}
                            </span>
                        </li>
                    ))}
                </ul>
            </div>
        </div>
    );
}
