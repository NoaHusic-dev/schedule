import { isoDate } from "#lib/helpers/format.ts";
import { fetchCurrentSchedule } from "#lib/services/schedules.ts";
import type { ClassSchedule, Schedule } from "#lib/types/schedule.ts";

const POLL_INTERVAL = 60 * 1000;
const CACHE_KEY = "schedule-api-cache";
const POLL_ERROR = 60 * 1000; // the server is unreachable, e.g. still booting

function createScheduleStore() {
  let classes = $state.raw<Record<string, ClassSchedule>>({});
  let date = $state("");
  let loaded = $state(false);

  let pending: Promise<void> | null = null;
  let failed = false;
  let refreshing = $state(false);
  let lastRefreshed = $state<number | null>(null);
  let offline = $state(false);
  let timer: ReturnType<typeof setTimeout> | null = null;
  let consumers = 0;

  try {
    const raw = localStorage.getItem(CACHE_KEY);
    if (raw) {
      const cached = JSON.parse(raw) as Schedule & { cachedAt?: number };
      classes = Object.fromEntries(cached.classes.map((c) => [c.class, c]));
      date = cached.date;
      lastRefreshed = cached.cachedAt ?? null;
      loaded = true;
    }
  } catch {
    localStorage.removeItem(CACHE_KEY);
  }

  // Single flight: overlapping callers share one request. A failure keeps the last good schedule
  // on screen; nobody interacts with the TV, so the next poll is the only retry that matters.
  function load(): Promise<void> {
    refreshing = true;
    pending ??= fetchCurrentSchedule()
      .then((schedule) => {
        classes = Object.fromEntries(schedule.classes.map((c) => [c.class, c]));
        date = schedule.date;
        failed = false;
        offline = false;
        lastRefreshed = Date.now();
        localStorage.setItem(CACHE_KEY, JSON.stringify({ ...schedule, cachedAt: Date.now() }));
      })
      .catch((e) => {
        failed = true;
        offline = !navigator.onLine;
        console.warn("Schedule refresh failed, keeping the last one", e);
      })
      .finally(() => {
        loaded = true;
        refreshing = false;
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
    const delay = failed ? POLL_ERROR : POLL_INTERVAL;
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
    get refreshing(): boolean {
      return refreshing;
    },
    get lastRefreshed(): number | null {
      return lastRefreshed;
    },
    get offline(): boolean {
      return offline;
    },
    refresh(): Promise<void> {
      if (pending) return pending;
      return load();
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
