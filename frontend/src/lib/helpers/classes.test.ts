/// <reference types="bun" />
import { describe, expect, test } from "bun:test";
import { CLASS_PAGES, formatClassName, formatGroup, ROWS_PER_PAGE } from "./classes";

describe("formatClassName", () => {
  test.each([
    ["1 EL", "1. EL"],
    ["4EL", "4. EL"],
    ["3SA", "3. SA"],
    ["2M", "2. M"],
    ["  3   SP ", "3. SP"],
    ["1el", "1. EL"],
  ])("%p → %p", (raw, formatted) => {
    expect(formatClassName(raw)).toBe(formatted);
  });
});

describe("formatGroup", () => {
  test("no groups is empty", () => {
    expect(formatGroup([])).toBe("");
    expect(formatGroup(undefined)).toBe("");
  });
  test("groups are listed in parentheses", () => {
    expect(formatGroup(["A (AEP)"])).toBe("(A (AEP))");
    expect(formatGroup(["1", "2"])).toBe("(1, 2)");
  });
});

describe("CLASS_PAGES", () => {
  test("each page fits the board", () => {
    for (const page of CLASS_PAGES) expect(page.length).toBeLessThanOrEqual(ROWS_PER_PAGE);
  });
  test("no class appears twice", () => {
    const all = CLASS_PAGES.flat();
    expect(new Set(all).size).toBe(all.length);
  });
});
