<template>
    <!-- <el-main style="display: flex; justify-content: center;">
        <transition name="el-fade-in-linear">
            <div v-show="!progress.successed" style="margin: 40px;">
                <el-progress type="circle" :percentage="progress.progressPercent" :status="progressState" />
                <p style="display: flex; justify-content: center;">正在下载 ({{ parms.options.downCount }} / {{ progress.downFinished }})</p>
            </div>
        </transition>
        <transition name="el-fade-in-linear">
            <el-result v-show="progress.successed" icon="success" title="下载完成" sub-title="您可以退出程序或再次下载" />
        </transition>
    </el-main> -->
    <FramePage title="正在下载">
        <h3>制作中</h3>
    </FramePage>
</template>

<script setup>
import FramePage from '../modules/frame_page.vue'
import { ref, onMounted, reactive, computed, watch } from 'vue'
import { StartDownload } from '../../../wailsjs/go/main/App'
import { EventsOn } from '../../../wailsjs/runtime'

// const props = defineProps(['parms', 'status'])
// const emit = defineEmits(['update:parms', 'update:status'])

// const parms = computed({
//     get() {
//         return props.parms
//     },
//     set(parms) {
//         emit('update:parms', parms)
//     }
// })

// // 底栏状态
// const status = computed({
//     get() {
//         return props.status
//     },
//     set(status) {
//         emit('update:status', status)
//     }
// })

// // 计算进度条是否为完成
// const progressState = computed(() => {
//     return progress.successed ? 'success' : ''
// })

// // 进度条相关
// const progress = reactive({
//     downFinished: 0,    // 任务完成数量
//     progressPercent: 0, // 任务完成百分比
//     successed: false,   // 任务完成状态
// })

// onMounted(() => {
//     // downloadVideoList()
// })

// // 下载视频列表
// function downloadVideoList() {

//     var opt = {
//         // song_name: parms.options.songName,
//         // song_cover: parms.options.songCover,
//         // song_author: parms.options.songAuthor,
//         song_name: true,
//         song_cover: true,
//         song_author: true,
//     }
//     StartDownload(opt).then(result => {

//         status.allowBack = true
//         status.allowNext = true

//     })
// }

// // 下载完成事件
// EventsOn("downloadFinish", () => {
//     progress.downFinished++
//     progress.progressPercent = (progress.downFinished / props.parms.options.downCount) * 100
//     // 判断任务是否完成
//     if (props.parms.options.downCount == progress.downFinished) {
//         progress.successed = true
//     } else {
//         progress.successed = false
//     }
// })
</script>