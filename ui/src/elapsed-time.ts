const numFormat = new Intl.NumberFormat("en-US");

const pluralize = (n: number, singular: string, plural: string) => {
  return n === 1
    ? `${numFormat.format(n)} ${singular} ago`
    : `${numFormat.format(n)} ${plural} ago`;
};

export const formatElapsedTime = (t: Date, now: Date = new Date()) => {
  const diffMs = now.getTime() - t.getTime();
  const diffSec = diffMs / 1000;
  const diffDays = diffSec / 86400;

  if (diffMs < 0) {
    return "in the future";
  }

  if (diffSec < 60) {
    return "now";
  } else if (diffSec < 3600) {
    return pluralize(Math.floor(diffSec / 60), "minute", "minutes");
  } else if (diffSec < 86400) {
    return pluralize(Math.floor(diffSec / 3600), "hour", "hours");
  } else if (Math.floor(diffDays) === 1) {
    return "yesterday";
  } else if (diffDays < 21) {
    return pluralize(Math.floor(diffDays), "day", "days");
  } else if (diffDays < 31) {
    return pluralize(Math.floor(diffDays / 7), "week", "weeks");
  } else if (diffDays < 365) {
    return pluralize(Math.floor(diffDays / 30), "month", "months");
  }
  return pluralize(Math.floor(diffDays / 365), "year", "years");
};
