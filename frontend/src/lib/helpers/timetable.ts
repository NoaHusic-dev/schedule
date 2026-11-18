import type { Lesson, TimeSlot } from "#lib/types/schedule.ts";

/** The school's bell schedule. */
export const TIMETABLE = [
  { id: 1, start: "08:00", end: "08:45" },
  { id: 2, start: "08:50", end: "09:35" },
  { id: 3, start: "09:40", end: "10:25" },
  { id: 4, start: "10:40", end: "11:25" },
  { id: 5, start: "11:30", end: "12:15" },
  { id: 6, start: "12:20", end: "13:05" },
  { id: 7, start: "13:10", end: "13:55" },
  { id: 8, start: "14:00", end: "14:45" },
] as const;

/** Shown instead of the period list for a block that covers the whole day. */
export const ALL_DAY = "Svi sati";

/** Converts HH:mm to minutes since midnight. */
export const timeToMinutes = (timeStr: string): number => {
  const [hours, minutes] = timeStr.split(":").map(Number);
  return hours * 60 + minutes;
};

/**
 * The periods a time range covers. A period counts only when it lies entirely inside the range.
 * @param startIso ISO date string, e.g. 2025-12-22T08:00:00+01:00
 */
export function periodIds(startIso: string, endIso: string): number[] {
  const startDate = new Date(startIso);
  const endDate = new Date(endIso);

  const startMins = startDate.getHours() * 60 + startDate.getMinutes();
  const endMins = endDate.getHours() * 60 + endDate.getMinutes();

  return TIMETABLE.filter((p) => timeToMinutes(p.start) >= startMins && timeToMinutes(p.end) <= endMins).map((p) => p.id);
}

/** "3", or "1–3" for consecutive periods (they always are), or ALL_DAY for a seven-period block. */
export function formatPeriods(ids: number[]): string {
  if (ids.length === 7) return ALL_DAY;
  if (ids.length > 1) return `${ids[0]}–${ids[ids.length - 1]}`;
  return ids.join("");
}

/** The period in progress, or during a break the one about to start; null after the last period. */
export function currentPeriodId(minutes: number): number | null {
  const active = TIMETABLE.find((p) => minutes >= timeToMinutes(p.start) && minutes <= timeToMinutes(p.end));
  if (active) return active.id;

  const nextPeriod = TIMETABLE.find((p) => timeToMinutes(p.start) > minutes);
  return nextPeriod ? nextPeriod.id : null;
}

/** One lesson to show in a row half, with the key its fade transition restarts on. */
export interface SlotView {
  slot: TimeSlot;
  lesson: Lesson;
  key: string;
}

export interface RowView {
  current: SlotView;
  next: SlotView | null;
}

const lessonKey = (lesson: Lesson) => `${lesson.subject.short}-${lesson.groups?.join("")}`;

/**
 * What a class's row shows in a period: the current lesson and the next one. Lessons held at the
 * same time (groups) take turns, one per tick.
 */
export function slotsFor(schedule: TimeSlot[], periodId: number | null, tick: number): RowView | null {
  if (!schedule.length) return null;
  if (!periodId) return null;

  const activeEntries = schedule.filter((item) => {
    const ids = periodIds(item.start, item.end);
    if (ids.length === 7) return periodId !== 8;
    return ids.includes(periodId);
  });
  if (activeEntries.length === 0) return null;

  const currentPool = activeEntries.flatMap((slot) => slot.lessons.map((lesson) => ({ slot, lesson })));
  const currentItem = currentPool[tick % currentPool.length];

  if (periodIds(currentItem.slot.start, currentItem.slot.end).length === 7) {
    return { current: { ...currentItem, key: currentItem.lesson.subject.short }, next: null };
  }

  const nextPool = schedule
    .filter((item) => item.start >= currentItem.slot.end)
    .sort((a, b) => a.start.localeCompare(b.start));

  let nextItem: { slot: TimeSlot; lesson: Lesson } | null = null;
  if (nextPool.length > 0) {
    const nextStartTime = nextPool[0].start;
    const flattenedNext = nextPool
      .filter((item) => item.start === nextStartTime)
      .flatMap((slot) => slot.lessons.map((lesson) => ({ slot, lesson })));
    if (flattenedNext.length > 0) {
      nextItem = flattenedNext[tick % flattenedNext.length];
    }
  }

  return {
    current: { ...currentItem, key: lessonKey(currentItem.lesson) },
    next: nextItem ? { ...nextItem, key: lessonKey(nextItem.lesson) } : null,
  };
}

/** The top bar's line: what is happening now, with a progress fraction for lessons and breaks. */
export interface PeriodStatus {
  title: string;
  detail: string;
  progress: number | null;
  bigBreak: boolean;
}

/**
 * @param hasToday whether the server's timetable is for today's date. Without it a weekend has no
 *   lessons, while on a weekday a holiday and a failed refresh look the same, so say only that the
 *   timetable isn't available.
 */
export function periodStatus(now: Date, hasToday: boolean): PeriodStatus {
  if (!hasToday) {
    const weekend = now.getDay() === 0 || now.getDay() === 6;
    return { title: weekend ? "Nema nastave danas" : "Raspored nije dostupan", detail: "", progress: null, bigBreak: false };
  }

  const minutes = now.getHours() * 60 + now.getMinutes() + now.getSeconds() / 60;
  const active = TIMETABLE.find((p) => minutes >= timeToMinutes(p.start) && minutes < timeToMinutes(p.end));
  if (active) {
    const start = timeToMinutes(active.start);
    const end = timeToMinutes(active.end);
    return {
      title: `${active.id}. sat`,
      detail: `još ${Math.ceil(end - minutes)} min`,
      progress: (minutes - start) / (end - start),
      bigBreak: false,
    };
  }

  const next = TIMETABLE.find((p) => timeToMinutes(p.start) > minutes);
  if (!next) return { title: "Nastava je završila", detail: "", progress: null, bigBreak: false };
  if (next.id === 1) return { title: "Nastava počinje", detail: `u ${next.start}`, progress: null, bigBreak: false };

  const breakStart = timeToMinutes(TIMETABLE[next.id - 2].end);
  const breakEnd = timeToMinutes(next.start);
  const bigBreak = next.id === 4;
  const secondsSinceMidnight =
    now.getHours() * 3600 + now.getMinutes() * 60 + now.getSeconds() + now.getMilliseconds() / 1000;
  const secondsRemaining = Math.ceil(breakEnd * 60 - secondsSinceMidnight);
  return {
    title: bigBreak ? "Veliki odmor" : "Odmor",
    detail: bigBreak
      ? `${Math.floor(secondsRemaining / 60)}:${String(secondsRemaining % 60).padStart(2, "0")}`
      : `${next.id}. sat u ${next.start}`,
    progress: (minutes - breakStart) / (breakEnd - breakStart),
    bigBreak,
  };
}
