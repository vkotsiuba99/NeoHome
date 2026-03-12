import clsx from "clsx";
import React from "react";
import styles from "./Badge.module.scss";

type Props = {
  children: React.ReactNode;
  tone?: "neutral" | "success" | "danger";
  className?: string;
};

export const Badge: React.FC<Props> = ({ children, tone = "neutral", className }) => {
  return <span className={clsx(styles.badge, styles[tone], className)}>{children}</span>;
};
