/// <reference types="bun" />
import { describe, expect, test } from "bun:test";
import type { Lesson, TimeSlot } from "$lib/types/schedule";
import { currentPeriodId, formatPeriods, periodIds, periodStatus, slotsFor, timeToMinutes } from "./timetable";

process.env.TZ = "Europe/Zagreb";

const at = (time: string) => `2026-09-28T${time}:00+02:00`;
const lesson = (short: string, groups: string[] = []): Lesson => ({
  subject: { name: short.toLowerCase(), short },
  teacher: { firstName: "", lastName: "", short: "T" },
  room: "1",
  groups,
});
const slot = (start: string, end: string, ...lessons: Lesson[]): TimeSlot => ({
  start: at(start),
  end: at(end),
  lessons,
});

describe("timeToMinutes", () => {
  test("counts minutes from midnight", () => {
    expect(timeToMinutes("08:50")).toBe(530);
  });
});

describe("periodIds", () => {
  test("one period", () => expect(periodIds(at("08:00"), at("08:45"))).toEqual([1]));
  test("a double period", () => expect(periodIds(at("08:00"), at("09:35"))).toEqual([1, 2]));
  test("the whole day", () => expect(periodIds(at("08:00"), at("14:45"))).toEqual([1, 2, 3, 4, 5, 6, 7, 8]));
  test("a range covering no whole period", () => expect(periodIds(at("08:10"), at("08:40"))).toEqual([]));
});

describe("formatPeriods", () => {
  test("one period", () => expect(formatPeriods([3])).toBe("3"));
  test("consecutive periods read as a range", () => expect(formatPeriods([1, 2, 3])).toBe("1–3"));
  test("seven periods read as the whole day", () => expect(formatPeriods([1, 2, 3, 4, 5, 6, 7])).toBe("Svi sati"));
  test("all eight periods are a range, so the column never overflows", () =>
    expect(formatPeriods([1, 2, 3, 4, 5, 6, 7, 8])).toBe("1–8"));
  test("no periods", () => expect(formatPeriods([])).toBe(""));
});

describe("currentPeriodId", () => {
  const minutes = (time: string) => timeToMinutes(time);
  test("before school it is the first period", () => expect(currentPeriodId(minutes("07:00"))).toBe(1));
  test("a period includes its start", () => expect(currentPeriodId(minutes("08:00"))).toBe(1));
  test("a period includes its end", () => expect(currentPeriodId(minutes("08:45"))).toBe(1));
  test("a break shows the coming period", () => expect(currentPeriodId(minutes("08:47"))).toBe(2));
  test("the long break shows the fourth period", () => expect(currentPeriodId(minutes("10:30"))).toBe(4));
  test("the last period", () => expect(currentPeriodId(minutes("14:45"))).toBe(8));
  test("after school there is none", () => expect(currentPeriodId(minutes("14:46"))).toBeNull());
});

describe("slotsFor", () => {
  const day = [
    slot("08:00", "08:45", lesson("MAT")),
    slot("08:50", "09:35", lesson("FIZ")),
    slot("09:40", "10:25", lesson("OE", ["1. grupa"]), lesson("RAČ", ["2. grupa"])),
    slot("10:40", "12:15", lesson("PN")),
  ];

  test("nothing without a schedule or a period", () => {
    expect(slotsFor([], 1, 0)).toBeNull();
    expect(slotsFor(day, null, 0)).toBeNull();
  });

  test("nothing when no lesson falls in the period", () => {
    expect(slotsFor(day, 7, 0)).toBeNull();
  });

  test("the current lesson and the next one", () => {
    const row = slotsFor(day, 1, 0)!;
    expect(row.current.lesson.subject.short).toBe("MAT");
    expect(row.current.key).toBe("MAT-");
    expect(row.next?.lesson.subject.short).toBe("FIZ");
  });

  test("groups taught at the same time take turns by tick", () => {
    expect(slotsFor(day, 3, 0)!.current.key).toBe("OE-1. grupa");
    expect(slotsFor(day, 3, 1)!.current.key).toBe("RAČ-2. grupa");
    expect(slotsFor(day, 3, 2)!.current.key).toBe("OE-1. grupa");
  });

  test("groups starting next take turns by tick too", () => {
    expect(slotsFor(day, 2, 0)!.next?.lesson.subject.short).toBe("OE");
    expect(slotsFor(day, 2, 1)!.next?.lesson.subject.short).toBe("RAČ");
  });

  test("a double period is current in both of its periods", () => {
    expect(slotsFor(day, 4, 0)!.current.lesson.subject.short).toBe("PN");
    expect(slotsFor(day, 5, 0)!.current.lesson.subject.short).toBe("PN");
    expect(slotsFor(day, 5, 0)!.next).toBeNull();
  });

  test("a whole-day block has no next lesson and is hidden in the eighth period", () => {
    const practice = [slot("08:00", "13:55", lesson("SP"))];
    const row = slotsFor(practice, 3, 0)!;
    expect(row.current.key).toBe("SP");
    expect(row.next).toBeNull();
    expect(slotsFor(practice, 8, 0)).toBeNull();
  });
});

describe("periodStatus", () => {
  // 2026-09-28 is a Monday.
  const monday = (time: string) => new Date(`2026-09-28T${time}+02:00`);

  test("during a lesson it counts down and shows progress", () => {
    const status = periodStatus(monday("10:00:00"), true);
    expect(status).toEqual({ title: "3. sat", detail: "još 25 min", progress: 20 / 45 });
  });

  test("the countdown rounds up, so the last minute reads 1 min", () => {
    expect(periodStatus(monday("10:24:30"), true).detail).toBe("još 1 min");
  });

  test("a break names the coming period and its start", () => {
    const status = periodStatus(monday("10:30:00"), true);
    expect(status.title).toBe("Odmor");
    expect(status.detail).toBe("4. sat u 10:40");
    expect(status.progress).toBeCloseTo(5 / 15);
  });

  test("a lesson ends exactly at its end time", () => {
    expect(periodStatus(monday("10:25:00"), true).title).toBe("Odmor");
  });

  test("before school", () => {
    expect(periodStatus(monday("07:45:00"), true)).toEqual({ title: "Nastava počinje", detail: "u 08:00", progress: null });
  });

  test("after school", () => {
    expect(periodStatus(monday("15:10:00"), true)).toEqual({ title: "Nastava je završila", detail: "", progress: null });
  });

  test("no timetable on a weekend means no lessons", () => {
    const saturday = new Date("2026-10-03T10:00:00+02:00");
    expect(periodStatus(saturday, false)).toEqual({ title: "Nema nastave danas", detail: "", progress: null });
  });

  // A holiday and a failed refresh look the same from here; this wording is true for both.
  test("no timetable on a weekday says it isn't available", () => {
    expect(periodStatus(monday("10:00:00"), false)).toEqual({ title: "Raspored nije dostupan", detail: "", progress: null });
  });
});
