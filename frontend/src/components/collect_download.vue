<template>
    <!-- 分页带 transform，fixed 会失效，用 sticky 让导航停在滚动区域底部 -->
    <div class="flex flex-col min-h-full">
        <var-tabs-items v-model:active="flow.pageIndex" class="grow">
            <var-tab-item>
                <AddVideos />
            </var-tab-item>

            <var-tab-item>
                <VideolistEditor />
            </var-tab-item>

            <var-tab-item>
                <DownloadProcess />
            </var-tab-item>
        </var-tabs-items>

        <footer class="sticky bottom-2.5 mx-2.5 mb-2.5">
            <var-space justify="space-between">
                <var-button type="primary" size="large" @click="flow.pageIndex--" :disabled="!nav.allowBack" v-show="nav.showBack">< 上一步</var-button>
                <var-badge type="danger" position="left-top" :value="flow.listCount" :hidden="!nav.showBadge">
                    <var-button type="primary" size="large" @click="next" :disabled="!nav.allowNext">下一步 ></var-button>
                </var-badge>
            </var-space>
        </footer>
    </div>
</template>

<script setup lang="ts">
import AddVideos from '../components/collect_download/add_videos.vue'
import VideolistEditor from '../components/collect_download/videolist_editor.vue'
import DownloadProcess from '../components/collect_download/download_process.vue'
import { computed, watch, onMounted } from 'vue'
import { flow, refreshListCount } from '../state'

// 按当前页面和任务数量计算翻页按钮状态
const nav = computed(() => {
    const hasTasks = flow.listCount > 0
    switch (flow.pageIndex) {
        case 0:
            return { showBack: false, allowBack: false, allowNext: hasTasks, showBadge: hasTasks }
        case 1:
            return { showBack: true, allowBack: true, allowNext: hasTasks, showBadge: hasTasks }
        default:
            return { showBack: true, allowBack: flow.download === 'idle', allowNext: flow.download === 'done', showBadge: false }
    }
})

// 下载完成后的下一步重新加载页面回到首页
function next() {
    if (flow.pageIndex >= 2) {
        window.location.reload()
        return
    }
    flow.pageIndex++
}

watch(() => flow.pageIndex, refreshListCount)

onMounted(refreshListCount)
</script>
