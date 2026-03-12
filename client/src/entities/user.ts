export type User = {
  userId: number;
  email: string;
  phone: string;
  login: string;
  role: string;
  createdAt: number;
  updatedAt: number;
};

export type UserWrapResponse = {
  user: User;
};

export type UpdateUserRequest = {
  email: string;
  password: string;
  login: string;
  phone: string;
};
