<template>
    <var-steps :active="pageIndex" style=" max-width: 50%; margin: 10px auto;">
        <var-step>查找收藏夹</var-step>
        <var-step>个性化选项</var-step>
        <var-step>编辑列表</var-step>
        <var-step>下载文件</var-step>
    </var-steps>
    <var-tabs-items v-model:active="pageIndex">
        <var-tab-item>
            <QueryCollect v-model:parms="parms" v-model:status="status" />
        </var-tab-item>

        <var-tab-item>
            <OptionPage v-model:parms="parms" v-model:status="status"  />
        </var-tab-item>

        <var-tab-item>
            <VideolistEditor v-model:parms="parms" v-model:status="status"  />
        </var-tab-item>

        <var-tab-item>
            <DownloadProcess v-model:parms="parms" v-model:status="status"  />
        </var-tab-item>
    </var-tabs-items>

    <footer class="page-turning" :style="{bottom: scrollTop + 'px'}">
        <var-space justify="space-between">
            <var-button type="primary" size="large" @click="pageIndex--" :disabled="!status.allowBack" v-show="status.showBack">< 上一步</var-button>
            <var-button type="primary" size="large" @click="pageIndex++" :disabled="!status.allowNext" v-show="status.showNext">下一步 ></var-button>
        </var-space>
    </footer>
</template>

<script setup>
import QueryCollect from '../components/collect_download/query_collect.vue'
import OptionPage from '../components/collect_download/option_page.vue'
import VideolistEditor from '../components/collect_download/videolist_editor.vue'
import DownloadProcess from '../components/collect_download/download_process.vue'
import { ref, reactive, watch, onMounted } from 'vue'

// 页面索引值
const pageIndex = ref(0)

// 底部翻页按钮距离
const scrollTop = ref(10)

// 收藏夹信息
const parms = reactive({
    favListID: "",
    count: 0,
    // 下载设置
    options: reactive({
        downCount: 0,
        downPart: true,
        songName: true,
        songCover: true,
        songAuthor: true,
    })
})

const status = reactive({
    showBack: false,
    allowBack: true,
    showNext: true,
    allowNext: false,
})

// 调节底部导航按钮位置
onMounted(() => {
    var scrollableElement = document.getElementById('scroll-box');

    scrollableElement.addEventListener('scroll', function(e) {
        scrollTop.value = 10 - e.target.scrollTop;
    });

})

// 导航按钮控制
watch(pageIndex, (newPageIndex) => {
    // 是否是第一页
    if (newPageIndex > 0) {
        status.showBack = true;
    } else {
        status.showBack = false;
    }
    // 是否进入下载页面
    if (newPageIndex > 2) {
        status.allowNext = false
    } else {
        status.allowNext = true

    }
    // 回到首页
    if (newPageIndex > 3) {
        pageIndex.value = 0;
    }
})
</script>

<style>
footer.page-turning {
    position: fixed;
    bottom: 10px;
    right: 10px;
    width: calc( 100% - 20px );
}
</style>