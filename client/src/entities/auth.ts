import { User, UserWrapResponse } from "./user";

export type RegisterRequest = {
  email: string;
  password: string;
  login: string;
  phone: string;
};

export type LoginEmailRequest = {
  email: string;
  password: string;
};

export type AuthResponse = {
  accessToken: string;
  expiresAt: number;
  user: User;
};

export type RegisterResponse = UserWrapResponse;
