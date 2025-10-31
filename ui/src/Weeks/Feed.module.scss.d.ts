export type Styles = {
  root: string;
  title: string;
  entries: string;
  more: string;
};

export type ClassNames = keyof Styles;

declare const styles: Styles;

export default styles;
