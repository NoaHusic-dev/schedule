import type { Schedule } from './schedule';

export interface Classes {
  classes: Class[];
}

export interface Class {
  class: string;
  schedule: Schedule[];
}
