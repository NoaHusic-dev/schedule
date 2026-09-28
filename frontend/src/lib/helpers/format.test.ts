/// <reference types="bun" />
import { describe, expect, test } from "bun:test";
import { formatClock, isoDate } from "./format";

describe("isoDate", () => {
  test("is the local calendar date, zero-padded", () => {
    expect(isoDate(new Date(2026, 8, 5, 23, 59))).toBe("2026-09-05");
    expect(isoDate(new Date(2026, 11, 31, 0, 0))).toBe("2026-12-31");
  });
});

describe("formatClock", () => {
  test("is 24-hour HH:mm:ss", () => {
    expect(formatClock(new Date(2026, 8, 28, 9, 5, 7))).toBe("09:05:07");
    expect(formatClock(new Date(2026, 8, 28, 14, 45, 0))).toBe("14:45:00");
  });
});
