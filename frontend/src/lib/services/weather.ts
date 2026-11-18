import { api } from "#lib/api/client.ts";
import type { CurrentWeather } from "#lib/types/weather.ts";

export const fetchCurrentWeather = (fetch?: typeof globalThis.fetch) =>
  api.get<CurrentWeather>("/weather/current", { fetch });
