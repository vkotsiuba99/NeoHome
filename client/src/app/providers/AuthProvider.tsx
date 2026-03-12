import React, { createContext, useContext, useEffect, useMemo, useState } from "react";
import { AuthResponse } from "@/entities/auth";
import { User } from "@/entities/user";
import { authEvents, authStore } from "@/shared/auth/token";
import { authApi } from "@/features/auth/api";

type AuthContextValue = {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (session: AuthResponse) => void;
  logout: () => void;
  setUser: (user: User) => void;
};

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export const AuthProvider: React.FC<React.PropsWithChildren> = ({ children }) => {
  const [token, setToken] = useState<string | null>(authStore.getToken());
  const [user, setUserState] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    return authEvents.onLogout(() => {
      authStore.clearToken();
      setToken(null);
      setUserState(null);
    });
  }, []);

  useEffect(() => {
    const initialize = async () => {
      const existingToken = authStore.getToken();
      if (!existingToken) {
        setIsLoading(false);
        return;
      }

      try {
        const profile = await authApi.me();
        setUserState(profile.user);
        setToken(existingToken);
      } catch {
        authStore.clearToken();
        setToken(null);
        setUserState(null);
      } finally {
        setIsLoading(false);
      }
    };

    void initialize();
  }, []);

  const login = (session: AuthResponse) => {
    authStore.setToken(session.accessToken);
    setToken(session.accessToken);
    setUserState(session.user);
  };

  const logout = () => {
    authStore.clearToken();
    setToken(null);
    setUserState(null);
  };

  const setUser = (nextUser: User) => {
    setUserState(nextUser);
  };

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      token,
      isLoading,
      isAuthenticated: Boolean(token && user),
      login,
      logout,
      setUser,
    }),
    [user, token, isLoading],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
};
