import { addHandler, sendMessage } from './gojs-internal/coordinator.js';

let handlers = [];

addHandler('http', function(message) {
  switch (message.type) {
    case 'handle':
      if (
        message.data.handler >= 0 &&
        handlers[message.data.handler] &&
        handlers[message.data.handler].handler
      ) {
        if (message.data.request) {
          message.replyData(
            handlers[message.data.handler].handler(message.data.request)
          );
          return;
        }

        throw new Error('made request to http_handle without request data');
      }

      throw new Error(
        'handler ' +
          message.data.handler +
          ' for http_handle invalid or not specified'
      );
    default:
      throw new Error(message.type + ' is not a valid type for http handler');
  }
});

export function handleRoute(route, handler) {
  const id =
    handlers.push({
      route,
      handler,
    }) - 1;

  sendMessage('http', 'route', {
    route,
    handler: id,
  });
}
