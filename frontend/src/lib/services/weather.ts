import { api } from "$lib/api/client";
import type { CurrentWeather } from "$lib/types/weather";

export const fetchCurrentWeather = (fetch?: typeof globalThis.fetch) =>
  api.get<CurrentWeather>("/weather/current", { fetch });
