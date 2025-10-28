const dayFormat = new Intl.DateTimeFormat("en-US", {
  year: "numeric",
  month: "2-digit",
  day: "2-digit",
});

export enum Weekday {
  Sunday = 0,
  Monday = 1,
  Tuesday = 2,
  Wednesday = 3,
  Thursday = 4,
  Friday = 5,
  Saturday = 6,
}

export class Day {
  private constructor(public readonly startsAt: Date) {}

  get endsAt(): Date {
    const { startsAt } = this;
    return new Date(
      new Date(
        startsAt.getFullYear(),
        startsAt.getMonth(),
        startsAt.getDate() + 1
      ).getTime() - 1
    );
  }

  add(days: number): Day {
    const { startsAt } = this;
    return new Day(
      new Date(
        startsAt.getFullYear(),
        startsAt.getMonth(),
        startsAt.getDate() + days
      )
    );
  }

  static of(date: Date): Day {
    return new Day(
      new Date(date.getFullYear(), date.getMonth(), date.getDate())
    );
  }

  toString(): string {
    return dayFormat.format(this.startsAt);
  }
}

export class Week {
  private constructor(public readonly startsAt: Date) {}

  add(weeks: number): Week {
    const { startsAt } = this;
    return new Week(
      new Date(
        startsAt.getFullYear(),
        startsAt.getMonth(),
        startsAt.getDate() + weeks * 7
      )
    );
  }

  get endsAt(): Date {
    const { startsAt } = this;
    return new Date(
      new Date(
        startsAt.getFullYear(),
        startsAt.getMonth(),
        startsAt.getDate() + 7
      ).getTime() - 1
    );
  }

  toString(): string {
    return `${dayFormat.format(this.startsAt)} - ${dayFormat.format(
      this.endsAt
    )}`;
  }

  static of(date: Date, weekday: Weekday): Week {
    const today = Day.of(date);
    const offset = weekday - today.startsAt.getDay();
    if (offset > 0) {
      return new Week(today.add(-7 + offset).startsAt);
    } else if (offset < 0) {
      return new Week(today.add(offset).startsAt);
    }
    return new Week(today.startsAt);
  }
}
