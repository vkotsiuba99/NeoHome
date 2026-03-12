import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useAuth } from "@/app/providers";
import { useUpdateMe } from "@/features/auth/hooks";
import { UpdateProfileValues, updateProfileSchema } from "@/features/auth/schema";
import { Card } from "@/shared/ui/Card/Card";
import { Input } from "@/shared/ui/Input/Input";
import { Button } from "@/shared/ui/Button/Button";
import { getErrorMessage } from "@/shared/utils/error";
import styles from "./ProfilePage.module.scss";

export const ProfilePage = () => {
  const { user, setUser } = useAuth();
  const updateMe = useUpdateMe();
  const [message, setMessage] = useState<string>("");

  const form = useForm<UpdateProfileValues>({
    resolver: zodResolver(updateProfileSchema),
    defaultValues: {
      email: user?.email ?? "",
      password: "",
      login: user?.login ?? "",
      phone: user?.phone ?? "",
    },
  });

  useEffect(() => {
    if (!user) {
      return;
    }
    form.reset({
      email: user.email,
      password: "",
      login: user.login,
      phone: user.phone,
    });
  }, [form, user]);

  const onSubmit = async (values: UpdateProfileValues) => {
    setMessage("");
    try {
      const response = await updateMe.mutateAsync(values);
      setUser(response.user);
      setMessage("Profile updated");
      form.setValue("password", "");
    } catch (error) {
      setMessage(getErrorMessage(error));
    }
  };

  return (
    <div className="container page">
      <h2>Profile</h2>
      <Card>
        <form className={styles.form} onSubmit={form.handleSubmit(onSubmit)}>
          <Input label="Email" type="email" {...form.register("email")} error={form.formState.errors.email?.message} />
          <Input label="Login" {...form.register("login")} error={form.formState.errors.login?.message} />
          <Input label="Phone" {...form.register("phone")} error={form.formState.errors.phone?.message} />
          <Input
            label="New password"
            type="password"
            {...form.register("password")}
            error={form.formState.errors.password?.message}
          />
          <Button type="submit" loading={updateMe.isPending}>
            Save profile
          </Button>
        </form>
        {message ? <p className={styles.message}>{message}</p> : null}
      </Card>
    </div>
  );
};
