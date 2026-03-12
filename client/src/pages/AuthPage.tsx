import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useAuth } from "@/app/providers";
import { useLoginEmail, useRegister } from "@/features/auth/hooks";
import { LoginEmailFormValues, loginEmailSchema, RegisterFormValues, registerSchema } from "@/features/auth/schema";
import { Card } from "@/shared/ui/Card/Card";
import { Input } from "@/shared/ui/Input/Input";
import { Button } from "@/shared/ui/Button/Button";
import { getErrorMessage } from "@/shared/utils/error";
import styles from "./AuthPage.module.scss";

type Mode = "login" | "register";

export const AuthPage = () => {
  const navigate = useNavigate();
  const { isAuthenticated, login } = useAuth();
  const [mode, setMode] = useState<Mode>("login");
  const [formError, setFormError] = useState<string>("");

  const loginMutation = useLoginEmail();
  const registerMutation = useRegister();

  const isBusy = loginMutation.isPending || registerMutation.isPending;

  const loginForm = useForm<LoginEmailFormValues>({
    resolver: zodResolver(loginEmailSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const registerForm = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      email: "",
      password: "",
      login: "",
      phone: "",
    },
  });

  useEffect(() => {
    if (isAuthenticated) {
      navigate("/dashboard", { replace: true });
    }
  }, [isAuthenticated, navigate]);

  const submitLogin = async (values: LoginEmailFormValues) => {
    setFormError("");
    try {
      const session = await loginMutation.mutateAsync(values);
      login(session);
      navigate("/dashboard", { replace: true });
    } catch (error) {
      setFormError(getErrorMessage(error));
    }
  };

  const submitRegister = async (values: RegisterFormValues) => {
    setFormError("");
    try {
      await registerMutation.mutateAsync(values);
      const session = await loginMutation.mutateAsync({
        email: values.email,
        password: values.password,
      });
      login(session);
      navigate("/dashboard", { replace: true });
    } catch (error) {
      setFormError(getErrorMessage(error));
    }
  };

  return (
    <div className="container page">
      <Card className={styles.card}>
        <div className={styles.tabs}>
          <Button variant={mode === "login" ? "primary" : "secondary"} onClick={() => setMode("login")}>
            Login
          </Button>
          <Button variant={mode === "register" ? "primary" : "secondary"} onClick={() => setMode("register")}>
            Register
          </Button>
        </div>

        {formError ? <p className={styles.error}>{formError}</p> : null}

        {mode === "login" ? (
          <form className={styles.form} onSubmit={loginForm.handleSubmit(submitLogin)}>
            <Input
              label="Email"
              type="email"
              {...loginForm.register("email")}
              error={loginForm.formState.errors.email?.message}
            />
            <Input
              label="Password"
              type="password"
              {...loginForm.register("password")}
              error={loginForm.formState.errors.password?.message}
            />
            <Button type="submit" loading={isBusy}>
              Sign in
            </Button>
          </form>
        ) : (
          <form className={styles.form} onSubmit={registerForm.handleSubmit(submitRegister)}>
            <Input
              label="Email"
              type="email"
              {...registerForm.register("email")}
              error={registerForm.formState.errors.email?.message}
            />
            <Input
              label="Password"
              type="password"
              {...registerForm.register("password")}
              error={registerForm.formState.errors.password?.message}
            />
            <Input
              label="Login"
              {...registerForm.register("login")}
              error={registerForm.formState.errors.login?.message}
            />
            <Input
              label="Phone"
              {...registerForm.register("phone")}
              error={registerForm.formState.errors.phone?.message}
            />
            <Button type="submit" loading={isBusy}>
              Create account
            </Button>
          </form>
        )}
      </Card>
    </div>
  );
};
