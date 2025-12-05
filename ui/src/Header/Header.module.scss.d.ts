export type Styles = {
  root: string;
  logo: string;
  icon: string;
  title: string;
};

export type ClassNames = keyof Styles;

declare const styles: Styles;

export default styles;
