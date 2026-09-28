/** The local calendar date as YYYY-MM-DD, the format of Schedule.date. */
export function isoDate(date: Date): string {
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`;
}

export function formatClock(date: Date): string {
  return date.toLocaleTimeString("hr-HR", { hour: "2-digit", minute: "2-digit", second: "2-digit", hour12: false });
}

/** "Utorak, 29. rujna". */
export function formatDate(date: Date): string {
  const text = date.toLocaleDateString("hr-HR", { weekday: "long", day: "numeric", month: "long" });
  return text.charAt(0).toUpperCase() + text.slice(1);
}
