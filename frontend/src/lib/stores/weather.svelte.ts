import { fetchCurrentWeather } from "#lib/services/weather.ts";
import type { CurrentWeather } from "#lib/types/weather.ts";

const POLL_INTERVAL = 10 * 60 * 1000;
const POLL_ERROR = 60 * 1000;

function createWeatherStore() {
  let current = $state.raw<CurrentWeather | null>(null);
  let loaded = $state(false);

  let pending: Promise<void> | null = null;
  let failed = false;
  let timer: ReturnType<typeof setTimeout> | null = null;
  let consumers = 0;

  function load(): Promise<void> {
    pending ??= fetchCurrentWeather()
      .then((weather) => {
        current = weather;
        failed = false;
      })
      .catch((e) => {
        failed = true;
        console.warn("Weather refresh failed, keeping the last one", e);
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

  function schedule() {
    stop();
    if (consumers === 0 || document.visibilityState === "hidden") return;
    timer = setTimeout(() => load().then(schedule), failed ? POLL_ERROR : POLL_INTERVAL);
  }

  function onVisibilityChange() {
    if (document.visibilityState === "visible") load().then(schedule);
    else stop();
  }

  return {
    get current(): CurrentWeather | null {
      return current;
    },
    get loaded(): boolean {
      return loaded;
    },
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

export const weather = createWeatherStore();
