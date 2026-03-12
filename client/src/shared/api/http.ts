import axios, { AxiosHeaders, AxiosRequestHeaders } from "axios";
import { env } from "@/shared/config";
import { authEvents, authStore } from "@/shared/auth/token";

export const http = axios.create({
  baseURL: env.apiBaseUrl,
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

http.interceptors.request.use((config) => {
  const authMode = config.authMode ?? "optional";
  const token = authStore.getToken();

  if (token && authMode !== "none") {
    const headers = AxiosHeaders.from((config.headers ?? {}) as AxiosRequestHeaders);
    headers.set("Authorization", `Bearer ${token}`);
    config.headers = headers;
  }

  return config;
});

http.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error?.response?.status === 401) {
      authStore.clearToken();
      authEvents.emitLogout();
    }

    return Promise.reject(error);
  },
);
