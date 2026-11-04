import { reactive } from 'vue'
import { GetListCount } from '../wailsjs/go/wails_api/WailsApi'

// 批量下载流程共享的状态
export const flow = reactive({
    pageIndex: 0,     // 0 添加内容，1 列表编辑，2 开始下载
    listCount: 0,     // 任务列表中的任务数量
    download: 'idle', // 下载状态：idle、running、done
})

// refreshListCount 从后端同步任务数量
export async function refreshListCount() {
    flow.listCount = await GetListCount()
}
