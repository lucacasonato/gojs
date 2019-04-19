import * as console from './lib/console.js';
import { handleRoute, fetch } from './lib/http.js';

handleRoute('/app', req => {
  return {
    statusCode: 200,
    headers: {},
    body: 'Hello',
  };
});

handleRoute('/google', req => {
  const r = fetch('http://www.mocky.io/v2/5185415ba171ea3a00704eed', {
    method: 'GET',
    headers: {},
    body: '',
  });

  console.log(r.json().hello);

  return {
    statusCode: 200,
    headers: {},
    body: 'Google',
  };
});
