import { z } from "zod";

export const loginEmailSchema = z.object({
  email: z.string().email("Valid email is required"),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

export type LoginEmailFormValues = z.infer<typeof loginEmailSchema>;

export const registerSchema = z.object({
  email: z.string().email("Valid email is required"),
  password: z.string().min(8, "Password must be at least 8 characters"),
  login: z.string().min(3, "Login must be at least 3 characters"),
  phone: z.string().min(8, "Phone is required"),
});

export type RegisterFormValues = z.infer<typeof registerSchema>;

export const updateProfileSchema = z.object({
  email: z.string().email("Valid email is required"),
  password: z
    .string()
    .refine((value) => value.length === 0 || value.length >= 8, "Password must be empty or at least 8 characters"),
  login: z.string().min(3, "Login must be at least 3 characters"),
  phone: z.string().min(8, "Phone is required"),
});

export type UpdateProfileValues = z.infer<typeof updateProfileSchema>;
