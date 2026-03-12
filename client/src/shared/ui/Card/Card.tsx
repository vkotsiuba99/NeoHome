import React from "react";
import clsx from "clsx";
import styles from "./Card.module.scss";

type Props = React.HTMLAttributes<HTMLDivElement>;

export const Card: React.FC<Props> = ({ className, children, ...props }) => {
  return (
    <section className={clsx(styles.card, className)} {...props}>
      {children}
    </section>
  );
};
