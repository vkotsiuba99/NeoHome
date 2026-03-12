import clsx from "clsx";
import React from "react";
import styles from "./Button.module.scss";

type Props = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "ghost" | "danger";
  loading?: boolean;
};

export const Button: React.FC<Props> = ({
  variant = "primary",
  loading = false,
  className,
  children,
  disabled,
  ...props
}) => {
  return (
    <button
      className={clsx(styles.button, styles[variant], className)}
      disabled={disabled || loading}
      {...props}
    >
      {loading ? "Please wait..." : children}
    </button>
  );
};
