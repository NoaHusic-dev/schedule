import { TIMETABLE, timeToMinutes } from "$util/timetable";

export class ClockState {
  time = $state(new Date());

  constructor() {
    // Update every second for the UI clock, but logic depends on minutes
    setInterval(() => {
      this.time = new Date();
    }, 1000);
  }

  get formatted() {
		return this.time.toLocaleTimeString('en-US', {
			hour12: false,
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit'
		});
  }

  /**
   * Helper: Current minutes since midnight
   */
  get totalMinutes() {
    return this.time.getHours() * 60 + this.time.getMinutes();
  }

  /**
   * Returns the ID of the current school period.
   * If during a break, it returns the upcoming period ID.
   */
  get currentPeriodId(): number | null {
    const m = this.totalMinutes;
    
    // 1. Check if we are inside a period
    const active = TIMETABLE.find(p => {
      const start = timeToMinutes(p.start);
      const end = timeToMinutes(p.end);
      return m >= start && m <= end;
    });
    if (active) return active.id;

    // 2. Check if we are in a break (find the next one starting today)
    const nextPeriod = TIMETABLE.find(p => timeToMinutes(p.start) > m);
    return nextPeriod ? nextPeriod.id : null;
  }

  /**
   * Returns the ID of the period following the current one.
   */
  get nextPeriodId(): number | null {
    const currentId = this.currentPeriodId;
    if (!currentId) return null;
    
    const next = currentId + 1;
    // Ensure the next ID actually exists in our timetable
    return TIMETABLE.some(p => p.id === next) ? next : null;
  }
}

export const clock = new ClockState();