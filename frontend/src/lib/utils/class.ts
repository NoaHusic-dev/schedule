/**
 * Formats class names to a standardized "N.EL" format.
 * Examples: "1EL" -> "1.EL", "2 EL" -> "2.EL", "3   EL" -> "3.EL"
 * @param className The raw class name string
 * @returns Standardized string
 */
export function formatClassName(className: string): string {
  // 1. Trim whitespace from start and end
  // 2. Use regex to find the boundary between digits and letters
  // (\d+) captures the number
  // \s* matches any optional whitespace
  // ([a-zA-Z]+) captures the class letters (like EL, A, B, etc)
  return className
    .trim()
    .replace(/^(\d+)\s*([a-zA-Z]+)$/, '$1. $2')
    .toUpperCase(); // Optional: ensures "1el" becomes "1.EL"
}

export function formatGroup(groups?: string[]): string {
  if (!groups || groups.length === 0) {
    return "";
  }
  return "(" + groups.join(", ") + ")";
}