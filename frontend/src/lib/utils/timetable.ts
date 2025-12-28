// Define the timetable structure
const TIMETABLE = [
  { id: 1, start: "08:00", end: "08:45" },
  { id: 2, start: "08:50", end: "09:35" },
  { id: 3, start: "09:40", end: "10:25" },
  { id: 4, start: "10:40", end: "11:25" },
  { id: 5, start: "11:30", end: "12:15" },
  { id: 6, start: "12:20", end: "13:05" },
  { id: 7, start: "13:10", end: "13:55" },
  { id: 8, start: "14:00", end: "14:45" },
];

/**
 * Converts a HH:mm string to total minutes from start of day
 */
const timeToMinutes = (timeStr: string): number => {
  const [hours, minutes] = timeStr.split(':').map(Number);
  return hours * 60 + minutes;
};

/**
 * Maps a time range to school period IDs
 * @param startIso ISO Date string (e.g., 2025-12-22T08:00:00+01:00)
 * @param endIso ISO Date string
 * @returns string of comma-separated period IDs
 */
export function getPeriodsFromRange(startIso: string, endIso: string): string {
  const startDate = new Date(startIso);
  const endDate = new Date(endIso);

  // Convert current input times to minutes since midnight
  const startMins = startDate.getHours() * 60 + startDate.getMinutes();
  const endMins = endDate.getHours() * 60 + endDate.getMinutes();

  const activePeriods = TIMETABLE.filter(period => {
    const periodStart = timeToMinutes(period.start);
    const periodEnd = timeToMinutes(period.end);

    // Logic: A period is included if the range overlaps with it.
    // Specifically, we check if the period is within the bounds of the provided start/end.
    return periodStart >= startMins && periodEnd <= endMins;
  });

  if (activePeriods.length === 7) {
    return "Svi sati";
  }

  return activePeriods.map(p => p.id).join(',');
}