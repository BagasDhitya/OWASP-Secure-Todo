import axios from "axios";

// helper ambil cookie
function getCookie(name: string): string | null {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop()!.split(";").shift() || null;
    return null;
}

const API = axios.create({
    baseURL: "http://localhost:8080/api", // ganti sesuai backend
    withCredentials: true, // supaya cookie ikut dikirim (access_token, refresh_token, csrf)
});

// interceptor untuk auto-attach CSRF token
API.interceptors.request.use((config) => {
    const csrfToken = getCookie("csrf");
    if (csrfToken) {
        config.headers["X-CSRF-Token"] = csrfToken;
    }
    return config;
});

export default API;
