import { fetchCurrentSchedule } from "$lib/services/schedules";
import type { Class } from "$lib/types/class";

export type ClassMap = Record<string, Class>;

export class ScheduleState {
  classes = $state<ClassMap>({});

  lastUpdated = $state<Date | null>(null);
  isLoading = $state(false);
  error = $state<string | null>(null);

  constructor() {
    this.refresh();

    setInterval(() => {
      this.refresh();
    }, 1000 * 60 * 60);
  }

  async refresh() {
    this.isLoading = true;
    this.error = null;

    try {
      const response = await fetchCurrentSchedule(); //
      const newClasses: ClassMap = {};

      if (response && response.classes) {
        for (const item of response.classes) {
          newClasses[item.class] = item;
        }
      }

      this.classes = newClasses;
      this.lastUpdated = new Date();

    } catch (err) {
      console.error("Failed to update schedule:", err);
      this.error = "Failed to load schedule";
    } finally {
      this.isLoading = false;
    }
  }

  getForClass(name: string): Class | undefined {
    return this.classes[name];
  }
}

export const scheduleStore = new ScheduleState();