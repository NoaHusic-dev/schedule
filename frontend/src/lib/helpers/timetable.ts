import type { Lesson, TimeSlot } from "$lib/types/schedule";

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
 * The periods a time range covers, as "1,2", or ALL_DAY for a seven-period block.
 * @param startIso ISO date string, e.g. 2025-12-22T08:00:00+01:00
 */
export function getPeriodsFromRange(startIso: string, endIso: string): string {
  const startDate = new Date(startIso);
  const endDate = new Date(endIso);

  const startMins = startDate.getHours() * 60 + startDate.getMinutes();
  const endMins = endDate.getHours() * 60 + endDate.getMinutes();

  // A period counts only when it lies entirely inside the range.
  const activePeriods = TIMETABLE.filter((period) => {
    const periodStart = timeToMinutes(period.start);
    const periodEnd = timeToMinutes(period.end);
    return periodStart >= startMins && periodEnd <= endMins;
  });

  if (activePeriods.length === 7) {
    return ALL_DAY;
  }

  return activePeriods.map((p) => p.id).join(",");
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
    const periods = getPeriodsFromRange(item.start, item.end);
    if (periods === ALL_DAY) return periodId !== 8;
    return periods.split(",").map(Number).includes(periodId);
  });
  if (activeEntries.length === 0) return null;

  const currentPool = activeEntries.flatMap((slot) => slot.lessons.map((lesson) => ({ slot, lesson })));
  const currentItem = currentPool[tick % currentPool.length];

  if (getPeriodsFromRange(currentItem.slot.start, currentItem.slot.end) === ALL_DAY) {
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
