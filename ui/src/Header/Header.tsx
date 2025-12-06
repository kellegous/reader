import styles from "./Header.module.scss";

export const Header = () => {
  return (
    <div className={styles.root}>
      <div className={styles.item}>
        <a href="/unread/">
          <div className={`${styles.icon} ${styles.unread}`}></div>
          <div className={styles.title}>unread</div>
        </a>
      </div>

      <div className={styles.item}>
        <a href="/feeds/">
          <div className={`${styles.icon} ${styles.feeds}`}></div>
          <div className={styles.title}>feeds</div>
        </a>
      </div>

      <div className={styles.item}>
        <a href="/starred/">
          <div className={`${styles.icon} ${styles.starred}`}></div>
          <div className={styles.title}>starred</div>
        </a>
      </div>

      <div className={styles.item}>
        <a href="/search/">
          <div className={`${styles.icon} ${styles.search}`}></div>
          <div className={styles.title}>search</div>
        </a>
      </div>

      <div className={styles.item}>
        <a href="/settings/">
          <div className={`${styles.icon} ${styles.search}`}></div>
          <div className={styles.title}>settings</div>
        </a>
      </div>
    </div>
  );
};
