export type Styles = {
  root: string;
  item: string;
  icon: string;
  title: string;

  // icons
  unread: string;
  feeds: string;
  starred: string;
  settings: string;
  search: string;
};

export type ClassNames = keyof Styles;

declare const styles: Styles;

export default styles;
