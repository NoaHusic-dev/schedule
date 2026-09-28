/** GET /schedules/current: today's timetable for every class. */
export interface Schedule {
  /** The day the timetable is for (YYYY-MM-DD), or "" before the server's first successful fetch. */
  date: string;
  classes: ClassSchedule[];
}

export interface ClassSchedule {
  class: string;
  schedule: TimeSlot[];
}

/** One or more lessons held at the same time, e.g. a class split into groups. */
export interface TimeSlot {
  start: string; // ISO 8601 with offset
  end: string;
  lessons: Lesson[];
}

export interface Lesson {
  subject: Subject;
  teacher: Teacher;
  room: string;
  groups: string[];
}

export interface Subject {
  name: string;
  short: string;
}

export interface Teacher {
  firstName: string;
  lastName: string;
  short: string;
}
