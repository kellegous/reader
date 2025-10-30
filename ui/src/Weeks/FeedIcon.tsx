import styles from "./FeedIcon.module.scss";
export interface FeedIconProps {
  url: string;
  title: string;
}

export const FeedIcon = ({ url, title }: FeedIconProps) => {
  return (
    <div
      className={styles.root}
      style={{ backgroundImage: `url(${url})` }}
      title={title}
    ></div>
  );
};
