import { api } from "$lib/api/client";
import type { Schedule } from "$lib/types/schedule";

export const fetchCurrentSchedule = (fetch?: typeof globalThis.fetch) =>
  api.get<Schedule>("/schedules/current", { fetch });
