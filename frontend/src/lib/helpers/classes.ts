/** Rows the board has room for between its header and footer. */
export const ROWS_PER_PAGE = 16;

/**
 * The board shows up to ROWS_PER_PAGE classes at a time and rotates through these pages. Names
 * must match Edupage exactly, spacing included.
 */
export const CLASS_PAGES: readonly (readonly string[])[] = [
  ["1 EL", "2 EL", "3 EL", "4EL", "1 MT", "2 MT", "3 MT", "4 MT", "1 PT", "2 PT", "3 PT", "4 PT", "1 RT", "2 RT", "3 RT", "4 RT"],
  ["1 K", "2 K", "3 K", "1 P", "2 P", "3 P", "1 S", "2 S", "3 S", "1 SE", "2 SE", "3SA", "3 SP", "2M", "3M"],
];

/** Normalizes a class name to "N. XX": "1EL" and "1 EL" both become "1. EL". */
export function formatClassName(className: string): string {
  return className
    .trim()
    .replace(/^(\d+)\s*([a-zA-Z]+)$/, "$1. $2")
    .toUpperCase();
}

/** Departments with their own colour in layout.css (--color-dep-*), keyed by the class name's letters. */
const DEPARTMENT_COLORS: Readonly<Record<string, string>> = {
  EL: "var(--color-dep-el)",
  MT: "var(--color-dep-mt)",
  PT: "var(--color-dep-pt)",
  RT: "var(--color-dep-rt)",
  SP: "var(--color-dep-sp)",
};

/** The colour of a class's row stripe: "1 EL" gets EL's. Undefined leaves the board's accent. */
export function departmentColor(className: string): string | undefined {
  const letters = className.trim().match(/^\d+\s*([a-zA-Z]+)$/)?.[1];
  return letters ? DEPARTMENT_COLORS[letters.toUpperCase()] : undefined;
}

export function formatGroup(groups?: string[]): string {
  if (!groups || groups.length === 0) {
    return "";
  }
  return "(" + groups.join(", ") + ")";
}
