import { http } from "@/shared/api";
import { AuthResponse, LoginEmailRequest, RegisterRequest, RegisterResponse } from "@/entities/auth";
import { UserWrapResponse, UpdateUserRequest } from "@/entities/user";

export const authApi = {
  register: async (payload: RegisterRequest) => {
    const response = await http.post<RegisterResponse>("/auth/register", payload, { authMode: "none" });
    return response.data;
  },
  loginEmail: async (payload: LoginEmailRequest) => {
    const response = await http.post<AuthResponse>("/auth/login/email", payload, { authMode: "none" });
    return response.data;
  },
  me: async () => {
    const response = await http.get<UserWrapResponse>("/users/me", { authMode: "required" });
    return response.data;
  },
  updateMe: async (payload: UpdateUserRequest) => {
    const response = await http.put<UserWrapResponse>("/users/me", payload, { authMode: "required" });
    return response.data;
  },
};
