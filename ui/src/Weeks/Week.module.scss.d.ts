export type Styles = {
  root: string;
  title: string;
  feeds: string;
};

export type ClassNames = keyof Styles;

declare const styles: Styles;

export default styles;
