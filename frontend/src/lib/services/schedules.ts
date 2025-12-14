import type { Classes } from "$lib/types/class";
import { api } from "$lib/api/client.js";

export async function fetchCurrentSchedule(customFetch?: typeof fetch) {
  return api.get<Classes>('/schedules/current', customFetch);
}