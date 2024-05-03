<template>
    <FramePage title="正在下载">
        <var-space direction="column" :size="[12, 12]">
            <var-progress type="info" :value="progress.progressPercent" :line-width="6" label :indeterminate="progress.downFinished == 0" v-show="!progress.successed" />
        </var-space>
        
        <var-result 
            title="下载完成"
            description="现在您可以继续使用软件了"
            v-show="progress.successed"
        />
    </FramePage>
</template>

<script setup>
import FramePage from '../modules/frame_page.vue'
import { reactive, computed } from 'vue'
import { EventsOn } from '../../../wailsjs/runtime'

const props = defineProps(['parms'])
const emit = defineEmits(['update:parms'])

const parms = computed({
    get() {
        return props.parms
    },
    set(parms) {
        emit('update:parms', parms)
    }
})

// 进度条相关
const progress = reactive({
    downFinished: 0,    // 任务完成数量
    progressPercent: 0, // 任务完成百分比
    successed: false,   // 任务完成状态
})

// 下载完成事件
EventsOn("downloadFinish", () => {
    progress.downFinished++
    progress.progressPercent = (progress.downFinished / props.parms.options.downCount) * 100

    // 判断任务是否完成
    if (props.parms.options.downCount == progress.downFinished) {
        progress.successed = true
    } else {
        progress.successed = false
    }
})
</script>