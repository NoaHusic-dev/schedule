export interface Lesson {
  subject: {
    name: string;
    short: string;
  };
  teacher: {
    firstName: string;
    lastName: string;
    short: string;
  };
  room: string;
}