import styles from "./Spinner.module.scss";

type Props = {
  label?: string;
};

export const Spinner = ({ label = "Loading..." }: Props) => {
  return (
    <div className={styles.wrap}>
      <span className={styles.spinner} />
      <span>{label}</span>
    </div>
  );
};
