

export interface BackendResponseMessage {
    methodName: string;
    data: any;
}

// 定义事件回调函数类型
type Callback = (data: BackendResponseMessage) => void;

// 事件管理器类
export class EventEmitter {
    private static instance: EventEmitter;
    private events: { [key: string]: Callback[] } = {};

    private constructor() {}

    // 获取单例实例
    static getInstance(): EventEmitter {
        if (!EventEmitter.instance) {
            EventEmitter.instance = new EventEmitter();
        }
        return EventEmitter.instance;
    }

    // 订阅事件
    on(methodName: string, callback: Callback) {
        if (!this.events[methodName]) {
            this.events[methodName] = [];
        }
        this.events[methodName].push(callback);
    }

    // 发布事件
    emit(methodName: string, data?: any) {
        const callbacks: Callback[] = this.events[methodName];
        if (callbacks) {
            callbacks.forEach(callback => {
                callback(data);
            });
        }
    }

    // 取消订阅事件
    off(methodName: string, callback: Callback) {
        const callbacks: Callback[] = this.events[methodName];
        if (callbacks) {
            this.events[methodName] = callbacks.filter(cb => cb !== callback);
        }
    }
}

