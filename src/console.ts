import { log as v8log } from './gojs-internal/coordinator';

export function log(...args: any[]) {
  args = args.map(arg => {
    if (typeof arg == 'object') return JSON.stringify(arg, null, 2);
    return arg;
  });
  v8log(
    `[${new Date().getHours()}:${new Date().getMinutes()}:${new Date().getSeconds()}] INFO:`,
    ...args
  );
}
