import { api } from "#lib/api/client.ts";
import type { Schedule } from "#lib/types/schedule.ts";

export const fetchCurrentSchedule = (fetch?: typeof globalThis.fetch) =>
  api.get<Schedule>("/schedules/current", { fetch });
