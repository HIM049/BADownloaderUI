import { EventsOn, EventsEmit } from '../wailsjs/runtime'

// 可添加到列表的来源
export type SourceKey = 'favlist' | 'season' | 'video' | 'audio' | 'space'

// 与后端 task.Progress 对应
export interface Progress {
    index: number
    title: string
}

// 事件名与回调参数
interface EventMap {
    downloadFinish: [progress: Progress]
    addToList: [kind: SourceKey, url: string]
    loginStatus: [status: string]
    qrcodeStr: [base64: string]
    refreshVideoList: []
}

// onEvent 监听事件，返回取消监听的函数
export function onEvent<K extends keyof EventMap>(name: K, callback: (...args: EventMap[K]) => void): () => void {
    return EventsOn(name, callback as (...data: any) => void)
}

// emitEvent 发出事件
export function emitEvent<K extends keyof EventMap>(name: K, ...args: EventMap[K]): void {
    EventsEmit(name, ...args)
}
