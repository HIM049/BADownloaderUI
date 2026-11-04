import { Snackbar } from '@varlet/ui'

type CallResult<T> = { ok: true; data: T } | { ok: false; data?: undefined }

// call 等待绑定方法返回，失败时弹出提示
export async function call<T>(promise: Promise<T>, failText = '操作失败'): Promise<CallResult<T>> {
    try {
        return { ok: true, data: await promise }
    } catch (err) {
        Snackbar.error(`${failText}：${err}`)
        return { ok: false }
    }
}
