import React from "react";
import styles from "./Textarea.module.scss";

type Props = React.TextareaHTMLAttributes<HTMLTextAreaElement> & {
  label?: string;
  error?: string;
};

export const Textarea = React.forwardRef<HTMLTextAreaElement, Props>(({ label, error, className, ...props }, ref) => {
  return (
    <label className={styles.wrap}>
      {label ? <span className={styles.label}>{label}</span> : null}
      <textarea ref={ref} className={`${styles.textarea} ${className ?? ""}`} {...props} />
      {error ? <span className={styles.error}>{error}</span> : null}
    </label>
  );
});
