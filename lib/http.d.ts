export declare type HTTPHandler = (req: HTTPRequest) => HTTPResponse;
export declare type RouteHandler = {
    route: string;
    handler: HTTPHandler;
};
export declare type HTTPRequest = {
    headers: {
        [header: string]: string[];
    };
    body: string;
};
declare type HTTPResponse = {
    headers: {
        [header: string]: string[];
    };
    status: number;
    body: string;
};
declare type FetchResponse = {
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
export declare function handleRoute(route: string, handler: HTTPHandler): void;
declare type FetchOptions = {
    method: 'GET' | 'POST' | 'PUT' | 'DELETE' | string;
    headers: {
        [header: string]: string[];
    };
    body: string;
};
export declare function fetch(resource: string, init: FetchOptions): FetchResponse;
export {};
