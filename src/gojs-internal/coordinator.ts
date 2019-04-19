export type Message = {
  id: string;
  namespace: string;
  type: string;
  data: any;
  replyData?: (data: any) => Message;
};

let handlers: {
  [namespace: string]: Function;
} = {};

export function sendMessage(
  namespace: string,
  type: string,
  data: any
): Message | null {
  return sendMessageWithID('', namespace, type, data);
}

export function sendMessageWithID(
  id: string,
  namespace: string,
  type: string,
  data: any
): Message | null {
  const m: Message = {
    id,
    namespace,
    type,
    data,
  };

  const msg = str2ab(JSON.stringify(m));

  const ab = V8Worker2.send(msg);
  if (ab) {
    const resp = JSON.parse(ab2str(ab));
    if (resp && resp.error) {
      throw Error(resp.error);
    }
    return resp;
  }

  return null;
}

export function addHandler(namespace: string, handler: Function): void {
  handlers[namespace] = handler;
}

export function log(...args: any[]) {
  V8Worker2.print(...args);
}

V8Worker2.recv(ab => {
  if (ab) {
    const resp = JSON.parse(ab2str(ab));

    if (resp.namespace && handlers[resp.namespace]) {
      resp.replyData = (data: any): Message => {
        return sendMessageWithID(resp.id, resp.namespace, resp.type, data);
      };
      handlers[resp.namespace](resp);
      return;
    }

    throw Error('No handler registered for namespace ' + resp.Namespace);
  }

  throw Error('Sent empty message');
});

export function ab2str(buf: ArrayBuffer): string {
  return String.fromCharCode.apply(null, new Uint8Array(buf));
}

export function str2ab(str: string): ArrayBuffer {
  var buf = new ArrayBuffer(str.length);
  var bufView = new Uint8Array(buf);
  for (var i = 0, strLen = str.length; i < strLen; i++) {
    bufView[i] = str.charCodeAt(i);
  }
  return buf;
}
