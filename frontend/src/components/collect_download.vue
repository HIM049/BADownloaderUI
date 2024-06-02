<template>
    <var-loading description="正在加载" :loading="status.loading">
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
    </var-loading>
</template>

<script setup>
import QueryCollect from '../components/collect_download/query_collect.vue'
import OptionPage from '../components/collect_download/option_page.vue'
import VideolistEditor from '../components/collect_download/videolist_editor.vue'
import DownloadProcess from '../components/collect_download/download_process.vue'
import { ref, reactive, watch, onMounted } from 'vue'
import { MakeAndSaveList, MakeAndSaveCompList, StartDownload } from '../../wailsjs/go/main/App'
import { Snackbar } from '@varlet/ui'

// 页面索引值
const pageIndex = ref(0)

// 底部翻页按钮距离
const scrollTop = ref(10)

// 收藏夹信息
const parms = reactive({
    favListID: "",
    mid: null,
    count: 0,
    isComp: null,
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
    showNext: true,
    allowNext: false,
    loading: false,
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
    // 列表编辑页面
    if (newPageIndex == 2) {
        status.loading = true;
        if (!parms.isComp) {
            MakeAndSaveList(parms.favListID, Number(parms.options.downCount), parms.options.downPart).then(result => {
                if (result != null) {
                    // 创建失败
                    Snackbar.error(result);
                } else {
                    // 创建成功
                    Snackbar.success("创建完成");
                }
                status.loading = false;
            })
        } else {
            MakeAndSaveCompList(Number(parms.mid), Number(parms.favListID), Number(parms.options.downCount), parms.options.downPart).then(result => {
                if (result != null) {
                    // 创建失败
                    Snackbar.error(result);
                } else {
                    // 创建成功
                    Snackbar.success("创建完成");
                }
                status.loading = false;
            })
        }
    }
    // 下载页面
    if (newPageIndex == 3) {
        status.allowNext = false
        status.allowBack = false
        var opt = {
            song_name: parms.options.songName,
            song_cover: parms.options.songCover,
            song_author: parms.options.songAuthor,
        }
        StartDownload(parms.favListID, opt).then(result => {
            status.allowNext = true
        })
    }
    // 回到首页
    if (newPageIndex > 3) {
        pageIndex.value = 0;
        window.location.reload();
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