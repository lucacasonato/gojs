import { addHandler, sendMessage, str2ab, } from './gojs-internal/coordinator';
let routeHandlers = [];
addHandler('http', (message) => {
    switch (message.type) {
        case 'handle':
            const data = message.data;
            if (data.handler >= 0 &&
                routeHandlers[data.handler] &&
                routeHandlers[data.handler].handler) {
                if (data.request) {
                    message.replyData(routeHandlers[data.handler].handler(data.request));
                    return;
                }
                throw new Error('made request to http_handle without request data');
            }
            throw new Error(`handler ${data.handler} for http_handle invalid or not specified`);
        default:
            throw new Error(`${message.type} is not a valid type for http handler`);
    }
});
export function handleRoute(route, handler) {
    const id = routeHandlers.push({
        route,
        handler,
    }) - 1;
    sendMessage('http', 'route', {
        route,
        handler: id,
    });
}
export function fetch(resource, init) {
    const resp = sendMessage('http', 'fetch', {
        resource,
        init,
    }).data;
    resp.text = () => {
        return resp.body;
    };
    resp.json = () => {
        return JSON.parse(resp.body);
    };
    resp.arrayBuffer = () => {
        return str2ab(resp.body);
    };
    return resp;
}
