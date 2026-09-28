function createClockStore() {
  let now = $state(new Date());
  let timer: ReturnType<typeof setTimeout> | null = null;
  let consumers = 0;

  // A chained timeout aligned to the next second, so the clock never drifts or skips a digit.
  function advance() {
    now = new Date();
    timer = setTimeout(advance, 1000 - (Date.now() % 1000));
  }

  return {
    get now(): Date {
      return now;
    },
    /** Minutes since midnight. */
    get minutes(): number {
      return now.getHours() * 60 + now.getMinutes();
    },
    /** Advances every 2 s; lessons held at the same time take turns on it. */
    get tick(): number {
      return Math.floor(now.getTime() / 2000);
    },
    /** `$effect(() => clock.subscribe())` keeps the clock running while a component shows it. */
    subscribe(): () => void {
      consumers++;
      if (timer === null) advance();
      return () => {
        consumers = Math.max(0, consumers - 1);
        if (consumers === 0 && timer !== null) {
          clearTimeout(timer);
          timer = null;
        }
      };
    },
  };
}

export const clock = createClockStore();
