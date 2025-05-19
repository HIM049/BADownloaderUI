<template>
    <var-tabs-items v-model:active="parms.pageIndex">
        <var-tab-item>
            <CreatVideolist v-model:parms="parms" v-model:status="status" @nextpage="parms.pageIndex++"/>
        </var-tab-item>

        <var-tab-item>
            <AddVideos v-model:parms="parms" v-model:status="status" @updateBadge="updateBadge" />
        </var-tab-item>

        <var-tab-item>
            <VideolistEditor v-model:parms="parms" v-model:status="status" @refresh="updateBadge" />
        </var-tab-item>

        <var-tab-item>
            <DownloadProcess v-model:parms="parms" v-model:status="status" />
        </var-tab-item>
    </var-tabs-items>

    <footer class="page-turning" :style="{bottom: scrollTop + 'px'}">
        <var-space justify="space-between">
            <var-button type="primary" size="large" @click="parms.pageIndex--" :disabled="!status.allowBack" v-show="status.showBack">< 上一步</var-button>
            <var-badge type="danger" position="left-top" :value="parms.listCount" :hidden="!status.showBadge">
                <var-button type="primary" size="large" @click="nextButton" :disabled="!status.allowNext" v-show="status.showNext">下一步 ></var-button>
            </var-badge>
        </var-space>
    </footer>
</template>

<script setup>
import CreatVideolist from '../components/collect_download/creat_videolist.vue'
import AddVideos from '../components/collect_download/add_videos.vue'
import VideolistEditor from '../components/collect_download/videolist_editor.vue'
import DownloadProcess from '../components/collect_download/download_process.vue'
import { ref, reactive, watch, onMounted } from 'vue'
import { GetListCount } from '../../wailsjs/go/wails_api/WailsApi'
import { Dialog, Snackbar } from '@varlet/ui'
import { EventsOn, EventsEmit } from '../../wailsjs/runtime'
// 页面索引值
const pageIndex = ref(0)

// 底部翻页按钮距离
const scrollTop = ref(10)

const nextButton = ref(() => {
    parms.pageIndex++
})

// 收藏夹信息
const parms = reactive({
    pageIndex: 0,
    videoListPath: "",
    listCount: 0,
    // 下载设置
    options: reactive({
        downCount: 0,
        downPart: true,
        songName: true,
        songCover: true,
        songAuthor: true,
    })
})

// 页面翻页按钮状态
const status = reactive({
    showBack: false,
    allowBack: true,
    showNext: false,
    allowNext: false,
    showBadge: false,
})

// 调节底部导航按钮位置
onMounted(() => {
    var scrollableElement = document.getElementById('scroll-box');

    scrollableElement.addEventListener('scroll', function(e) {
        scrollTop.value = 10 - e.target.scrollTop;
    });

})

watch(parms, (newValue) => {
    if (newValue.listCount <= 0) {
        status.allowNext = false;
    } else {
        status.allowNext =  true;
    }
});

// 导航按钮控制
watch(parms, (newPageIndex) => {
    status.showBack = true;
    status.showNext = true;
    status.allowBack = true;
    status.allowNext = true;
    status.showBadge = false;

    nextButton.value = () => {
        parms.pageIndex++
    };
    // 列表选择或创建
    if (newPageIndex.pageIndex == 0) {
        status.showBack = false;
        status.showNext = false;
    }
    // 添加视频
    if (newPageIndex.pageIndex == 1) {
        status.allowNext = false;
        updateBadge();
    }
    // 列表编辑页面
    if (newPageIndex.pageIndex == 2) {
        nextButton.value = () => {
            // Dialog('清理删除项并下一步？').then(result => {
            //     if (result == 'confirm') {
            //         // TidyVideoList(parms.videoListPath);
            //         parms.pageIndex++;
            //     }
            //
            //     return;
            // });
            parms.pageIndex++;
        };
        EventsEmit('refreshVideoList');
        updateBadge();
    }
    // 下载页面
    if (newPageIndex.pageIndex == 3) {
        status.allowNext = false
        // status.allowBack = false
    }
    // 回到首页
    if (newPageIndex.pageIndex > 3) {
        pageIndex.value = 0;
        window.location.reload();
    }
})

// 更新列表数量显示
function updateBadge() {
    GetListCount(parms.videoListPath).then(result => {
        // badgeValue.value = result;
        parms.listCount = result;
        if (parms.listCount <= 0) {
            status.showBadge = false;
            status.allowNext = false;
        } else {
            status.showBadge = true;
            status.allowNext = true;
        }
    });
}
</script>

<style>
footer.page-turning {
    position: fixed;
    bottom: 10px;
    right: 10px;
    width: calc( 100% - 20px );
}
</style>