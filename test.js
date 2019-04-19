import * as console from './lib/console.js';
import { handleRoute } from './lib/http.js';

handleRoute('/app', req => {
  return {
    statusCode: 200,
    headers: {},
    body: 'Hello',
  };
});

handleRoute('/google', req => {
  console.log(req);

  return {
    statusCode: 200,
    headers: {},
    body: 'Hello',
  };
});
