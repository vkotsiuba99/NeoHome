import React from "react";
import styles from "./Input.module.scss";

type Props = React.InputHTMLAttributes<HTMLInputElement> & {
  label?: string;
  error?: string;
};

export const Input = React.forwardRef<HTMLInputElement, Props>(({ label, error, className, ...props }, ref) => {
  return (
    <label className={styles.wrap}>
      {label ? <span className={styles.label}>{label}</span> : null}
      <input ref={ref} className={`${styles.input} ${className ?? ""}`} {...props} />
      {error ? <span className={styles.error}>{error}</span> : null}
    </label>
  );
});
