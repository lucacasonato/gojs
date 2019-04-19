export declare type Message = {
    id: string;
    namespace: string;
    type: string;
    data: any;
    replyData?: (data: any) => Message;
};
export declare function sendMessage(namespace: string, type: string, data: any): Message | null;
export declare function sendMessageWithID(id: string, namespace: string, type: string, data: any): Message | null;
export declare function addHandler(namespace: string, handler: Function): void;
export declare function log(...args: any[]): void;
export declare function ab2str(buf: ArrayBuffer): string;
export declare function str2ab(str: string): ArrayBuffer;
