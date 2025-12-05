import styles from "./Header.module.scss";

export const Header = () => {
  return (
    <div className={styles.root}>
      <a href="/ui/" className={styles.logo}>
        <div className={styles.icon}></div>
        <div className={styles.title}>reader</div>
      </a>
    </div>
  );
};
