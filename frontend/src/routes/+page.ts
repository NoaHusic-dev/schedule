import type { PageLoad } from './$types';
import { fetchCurrentSchedule } from '$lib/services/schedules';

export const load = (async ({ fetch }) => {
  const classes = await fetchCurrentSchedule(fetch);
  
  return {
    classes
  };
}) satisfies PageLoad;