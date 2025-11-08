export type Styles = {
  root: string;
  title: string;
  info: string;
  read: string;
  removed: string;
  summary: string;
  summarize: string;
  hidden: string;
};

export type ClassNames = keyof Styles;

declare const styles: Styles;

export default styles;
