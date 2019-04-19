function ab2str(buf) {
  return String.fromCharCode.apply(null, new Uint8Array(buf));
}

function str2ab(str) {
  var buf = new ArrayBuffer(str.length);
  var bufView = new Uint8Array(buf);
  for (var i = 0, strLen = str.length; i < strLen; i++) {
    bufView[i] = str.charCodeAt(i);
  }
  return buf;
}

let handlers = {};

export function sendMessage(namespace, type, data) {
  return sendMessageWithID('', namespace, type, data);
}

export function sendMessageWithID(id, namespace, type, data) {
  const msg = str2ab(
    JSON.stringify({
      id,
      namespace,
      type,
      data,
    })
  );

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

export function addHandler(namespace, handler) {
  handlers[namespace] = handler;
}

export function log(...args) {
  V8Worker2.print(...args);
}

V8Worker2.recv(ab => {
  if (ab) {
    const resp = JSON.parse(ab2str(ab));

    if (resp.namespace && handlers[resp.namespace]) {
      resp.replyData = function(data) {
        sendMessageWithID(resp.id, resp.namespace, resp.type, data);
      };
      handlers[resp.namespace](resp);
      return;
    }

    throw Error('No handler registered for namespace ' + resp.Namespace);
  }

  throw Error('Sent empty message');
});
