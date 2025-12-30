import { TIMETABLE, timeToMinutes } from "$util/timetable";

export class ClockState {
  time = $state(new Date());
  tick = $state(0);

  constructor() {
    setInterval(() => {
      this.time = new Date();
    }, 1000);

    setInterval(() => {
      this.tick += 1;
    }, 5000);
  }

  get formatted() {
    return this.time.toLocaleTimeString('en-US', {
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  }

  get totalMinutes() {
    return this.time.getHours() * 60 + this.time.getMinutes();
  }

  get currentPeriodId(): number | null {
    const m = this.totalMinutes;
    
    const active = TIMETABLE.find(p => {
      const start = timeToMinutes(p.start);
      const end = timeToMinutes(p.end);
      return m >= start && m <= end;
    });
    if (active) return active.id;

    const nextPeriod = TIMETABLE.find(p => timeToMinutes(p.start) > m);
    return nextPeriod ? nextPeriod.id : null;
  }

  get nextPeriodId(): number | null {
    const currentId = this.currentPeriodId;
    if (!currentId) return null;
    
    const next = currentId + 1;
    return TIMETABLE.some(p => p.id === next) ? next : null;
  }
}

export const clock = new ClockState();