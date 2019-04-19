import { log as v8log } from './gojs-internal/coordinator.js';

export function log(...args) {
  v8log(args);
}
