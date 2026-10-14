/// <reference types="bun" />
import { describe, expect, test } from "bun:test";
import { CLASS_PAGES, departmentColor, formatClassName, formatGroup, ROWS_PER_PAGE } from "./classes";

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

describe("departmentColor", () => {
  test.each([
    ["1 EL", "var(--color-dep-el)"],
    ["4EL", "var(--color-dep-el)"],
    ["2 MT", "var(--color-dep-mt)"],
    ["3 PT", "var(--color-dep-pt)"],
    ["4 RT", "var(--color-dep-rt)"],
    ["  3   sp ", "var(--color-dep-sp)"],
  ])("%p → %p", (name, color) => {
    expect(departmentColor(name)).toBe(color);
  });
  test("a department without its own colour falls back to the accent", () => {
    expect(departmentColor("1 K")).toBeUndefined();
    expect(departmentColor("3SA")).toBeUndefined();
    expect(departmentColor("")).toBeUndefined();
  });
});
