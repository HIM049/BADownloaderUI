<template>
    <FramePage title="开始下载">
        <var-space justify="center">
            <var-button type="primary" @click="startDownload" size="large" v-if="flow.download === 'idle'"><var-icon name="download" />开始下载</var-button>
        </var-space>
        
        
    
        <div v-if="flow.download === 'running'">
            <var-space direction="column" :size="[12, 12]">
                <var-progress type="info" :value="progress.progressPercent" :line-width="6" label :indeterminate="progress.downFinished == 0" />
            </var-space>

            <var-space justify="center">
                <var-chip type="info" class="mt-2.5">正在下载：{{ progress.downloadingTitle }} ( {{ progress.downFinished }} / {{ flow.listCount }} )</var-chip>
                
            </var-space>
        </div>

        <var-result 
            title="下载完成"
            description="请点击下一步回到主页"
            v-show="flow.download === 'done'"
        >
            <template #footer>
                <var-button @click="OpenDownloadFolader" type="primary">打开下载文件夹</var-button>
            </template>
        </var-result>
    </FramePage>
</template>

<script setup>
import FramePage from '../modules/frame_page.vue'
import { reactive } from 'vue'
import { DownloadTaskList, OpenDownloadFolader } from '../../../wailsjs/go/wails_api/WailsApi'
import { EventsOn } from '../../../wailsjs/runtime'
import { flow } from '../../state'

// 进度条相关
const progress = reactive({
    downFinished: 0,    // 任务完成数量
    progressPercent: 0, // 任务完成百分比
    downloadingTitle: "", // 下载中标题
})

async function startDownload() {
    flow.download = 'running';
    await DownloadTaskList();
    flow.download = 'done';
}


// 下载完成事件
EventsOn("downloadFinish", (task) => {
    progress.downFinished++;
    progress.downloadingTitle = task.title;
    progress.progressPercent = (progress.downFinished / flow.listCount) * 100;
})
</script>