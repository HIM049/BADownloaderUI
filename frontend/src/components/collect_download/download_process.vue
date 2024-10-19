<template>
    <FramePage title="开始下载">
        <var-space justify="center">
            <var-button type="primary" @click="startDownload" size="large" v-if="!downloading && !progress.successed"><var-icon name="download" />开始下载</var-button>
        </var-space>
        
        
    
        <div v-if="downloading">
            <var-space direction="column" :size="[12, 12]">
                <var-progress type="info" :value="progress.progressPercent" :line-width="6" label :indeterminate="progress.downFinished == 0" v-show="!progress.successed" />
            </var-space>

            <var-space justify="center">
                <var-chip type="info" style="margin-top: 10px;">正在下载：{{ progress.downloadingTitle }} ( {{ progress.downFinished }} / {{ parms.listCount }} )</var-chip>
                
            </var-space>
        </div>

        <var-result 
            title="下载完成"
            description="请点击下一步回到主页"
            v-show="progress.successed"
        >
            <template #footer>
                <var-button @click="OpenDownloadFolader" type="primary">打开下载文件夹</var-button>
            </template>
        </var-result>
    </FramePage>
</template>

<script setup>
import FramePage from '../modules/frame_page.vue'
import { ref, reactive, computed } from 'vue'
import { ListDownload, OpenDownloadFolader } from '../../../wailsjs/go/main/App'
import { EventsOn } from '../../../wailsjs/runtime'

const downloading = ref(false)

// 进度条相关
const progress = reactive({
    downFinished: 0,    // 任务完成数量
    progressPercent: 0, // 任务完成百分比
    successed: false,   // 任务完成状态
    downloadingTitle: "", // 下载中标题
})

const props = defineProps(['parms', 'status'])
const emit = defineEmits(['update:parms', 'update:status'])

const parms = computed({
    get() {
        return props.parms
    },
    set(parms) {
        emit('update:parms', parms)
    }
})

const status = computed({
    get() {
        return props.status
    },
    set(status) {
        emit('update:status', status)
    }
})

function startDownload() {
    status.value.allowBack = false;
    downloading.value = true;
    var opt = {
        song_name: parms.value.options.songName,
        song_cover: parms.value.options.songCover,
        song_author: parms.value.options.songAuthor,
    }
    ListDownload(parms.value.videoListPath, opt).then(result => {
        downloading.value = false;
        status.value.allowNext = true;
        progress.successed = true;
    })
}


// 下载完成事件
EventsOn("downloadFinish", (title) => {
    progress.downFinished++;
    progress.downloadingTitle = title;
    progress.progressPercent = (progress.downFinished / parms.value.listCount) * 100;
})
</script>