import React from "react";
import styles from "./Input.module.scss";

type Props = React.InputHTMLAttributes<HTMLInputElement> & {
  label?: string;
  error?: string;
};

export const Input: React.FC<Props> = ({ label, error, className, ...props }) => {
  return (
    <label className={styles.wrap}>
      {label ? <span className={styles.label}>{label}</span> : null}
      <input className={`${styles.input} ${className ?? ""}`} {...props} />
      {error ? <span className={styles.error}>{error}</span> : null}
    </label>
  );
};
