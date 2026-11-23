import { isoDate } from "#lib/helpers/format.ts";
import { fetchCurrentSchedule } from "#lib/services/schedules.ts";
import type { ClassSchedule } from "#lib/types/schedule.ts";

const POLL_TODAY = 60 * 60 * 1000; // the server refreshes once a day, at 05:00
const POLL_STALE = 5 * 60 * 1000; // waiting for today's timetable: a weekend, or a refresh that failed
const POLL_ERROR = 60 * 1000; // the server is unreachable, e.g. still booting

function createScheduleStore() {
  let classes = $state.raw<Record<string, ClassSchedule>>({});
  let date = $state("");
  let loaded = $state(false);

  let pending: Promise<void> | null = null;
  let failed = false;
  let timer: ReturnType<typeof setTimeout> | null = null;
  let consumers = 0;

  // Single flight: overlapping callers share one request. A failure keeps the last good schedule
  // on screen; nobody interacts with the TV, so the next poll is the only retry that matters.
  function load(): Promise<void> {
    pending ??= fetchCurrentSchedule()
      .then((schedule) => {
        classes = Object.fromEntries(schedule.classes.map((c) => [c.class, c]));
        date = schedule.date;
        failed = false;
      })
      .catch((e) => {
        failed = true;
        console.warn("Schedule refresh failed, keeping the last one", e);
      })
      .finally(() => {
        loaded = true;
        pending = null;
      });
    return pending;
  }

  function stop() {
    if (timer !== null) clearTimeout(timer);
    timer = null;
  }

  // A chained timeout, so a slow response never stacks requests. It runs only while someone
  // shows the board and the tab is visible.
  function schedule() {
    stop();
    if (consumers === 0 || document.visibilityState === "hidden") return;
    const delay = failed ? POLL_ERROR : date === isoDate(new Date()) ? POLL_TODAY : POLL_STALE;
    timer = setTimeout(() => load().then(schedule), delay);
  }

  function onVisibilityChange() {
    if (document.visibilityState === "visible") load().then(schedule);
    else stop();
  }

  return {
    /** False until the first response (or failure) arrives. */
    get loaded(): boolean {
      return loaded;
    },
    /** The day the timetable is for (YYYY-MM-DD), or "" when there is none yet. */
    get date(): string {
      return date;
    },
    get(className: string): ClassSchedule | undefined {
      return classes[className];
    },
    get classes(): ClassSchedule[] {
      return Object.values(classes).sort((a, b) => a.class.localeCompare(b.class, "hr"));
    },
    /** `$effect(() => schedules.subscribe())` polls while a component shows the board. */
    subscribe(): () => void {
      if (consumers++ === 0) {
        document.addEventListener("visibilitychange", onVisibilityChange);
        load().then(schedule);
      }
      return () => {
        consumers = Math.max(0, consumers - 1);
        if (consumers === 0) {
          stop();
          document.removeEventListener("visibilitychange", onVisibilityChange);
        }
      };
    },
  };
}

export const schedules = createScheduleStore();
