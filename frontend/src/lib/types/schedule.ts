import type { Lesson } from './lesson';

export interface Schedule {
  start: string;
  end: string;
  lessons: Lesson[];
}