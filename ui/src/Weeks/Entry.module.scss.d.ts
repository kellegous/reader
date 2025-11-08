export type Styles = {
  root: string;
  title: string;
  info: string;
  read: string;
  removed: string;
  summary: string;
};

export type ClassNames = keyof Styles;

declare const styles: Styles;

export default styles;
