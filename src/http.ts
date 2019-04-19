import {
  addHandler,
  sendMessage,
  Message,
  str2ab,
} from './gojs-internal/coordinator';

/* ---------------------------------- */
/* HTTP Server Side Handling          */
/* ---------------------------------- */

export type HTTPHandler = (req: HTTPRequest) => HTTPResponse;

export type RouteHandler = {
  route: string;
  handler: HTTPHandler;
};

export type HTTPRequest = {
  headers: { [header: string]: string[] };
  body: string;
};

type HTTPResponse = {
  headers: {
    [header: string]: string[];
  };
  status: number;
  body: string;
};

type FetchResponse = {
  headers: {
    [header: string]: string[];
  };
  status: number;
  body: string;
  url: string;
  ok: boolean;
  text: () => string;
  json: () => any;
  arrayBuffer: () => ArrayBuffer;
};

let routeHandlers: RouteHandler[] = [];

addHandler('http', (message: Message) => {
  switch (message.type) {
    case 'handle':
      const data: {
        handler: number;
        request: HTTPRequest;
      } = message.data;

      if (
        data.handler >= 0 &&
        routeHandlers[data.handler] &&
        routeHandlers[data.handler].handler
      ) {
        if (data.request) {
          message.replyData(routeHandlers[data.handler].handler(data.request));
          return;
        }

        throw new Error('made request to http_handle without request data');
      }

      throw new Error(
        `handler ${data.handler} for http_handle invalid or not specified`
      );
    default:
      throw new Error(`${message.type} is not a valid type for http handler`);
  }
});

export function handleRoute(route: string, handler: HTTPHandler) {
  const id =
    routeHandlers.push({
      route,
      handler,
    }) - 1;

  sendMessage('http', 'route', {
    route,
    handler: id,
  });
}

/* ---------------------------------- */
/* HTTP Fetch                         */
/* ---------------------------------- */

type FetchOptions = {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | string;
  headers: {
    [header: string]: string[];
  };
  body: string;
};

export function fetch(resource: string, init: FetchOptions): FetchResponse {
  const resp = <FetchResponse>sendMessage('http', 'fetch', {
    resource,
    init,
  }).data;

  resp.text = (): string => {
    return resp.body;
  };
  resp.json = (): any => {
    return JSON.parse(resp.body);
  };
  resp.arrayBuffer = (): ArrayBuffer => {
    return str2ab(resp.body);
  };

  return resp;
}
