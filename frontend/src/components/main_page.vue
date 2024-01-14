<template>
    <el-container direction="vertical">
        <HeadBar @switch-setting="pageNum = -10" @switch-sdown="pageNum = -11" />
        <section id="app-function">
            <HomePage v-if="pageNum == 0" @back="pageNum--" @next="pageNum++" v-model:parms="parms" />
            <keep-alive>
                <OptionPage v-if="pageNum == 1" @back="pageNum--" @next="pageNum++" v-model:parms="parms" />
            </keep-alive>
            <EditListPage v-if="pageNum == 2" @back="pageNum--" @next="pageNum++" v-model:parms="parms" />
            <DownloadPage v-if="pageNum == 3" @back="pageNum--" @next="pageNum = 0" v-model:parms="parms" />
            <SettingPage v-if="pageNum == -10" @back="pageNum = 0" @next="" v-model:parms="parms" />
            <SingleDownload v-if="pageNum == -11" @back="pageNum = 0" @next="" v-model:parms="parms" />
        </section>
    </el-container>
</template>

<script setup>
import SettingPage from '../components/setting_page.vue'
import SingleDownload from '../components/single_downlaod.vue'
import HeadBar from '../components/modules/head_bar.vue'
import HomePage from '../components/home_page.vue'
import OptionPage from '../components/option_page.vue'
import EditListPage from '../components/edit_videolist_page.vue'
import DownloadPage from '../components/download_page.vue'
import { ref, shallowRef, reactive, computed, watch } from 'vue'

const pageNum = ref(0)
const isSetting = ref(false)
const isSdown = ref(false)
const page = shallowRef(HomePage)
const footbarText = ref("继续")

const status = reactive({
    showBack: false,
    allowBack: true,
    showNext: true,
    allowNext: false,
    loadingNext: false,
    showStep: false,
})

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

// // 切换到设置页面
// watch(isSetting, (newValue) => {
//     if (newValue) {
//         page.value = SettingPage;
//         status.showBack = false;
//         status.showNext = false;
//         status.showStep = false;
//     } else {
//         page.value = HomePage;
//         status.showNext = true;
//         pageNum.value = 0;
//     }
// })

// // 切换到单曲下载
// watch(isSdown, (newValue) => {
//     if (newValue) {
//         page.value = SingleDownload;
//         status.showBack = false;
//         status.showNext = false;
//         status.showStep = false;
//     } else {
//         page.value = HomePage;
//         status.showNext = true;
//         pageNum.value = 0;
//     }
// })

// // 下载视频列表
// function downloadVideoList() {

//     // 加载页面
//     const loading = ElLoading.service({
//         lock: true,
//         text: '正在下载音频',
//         background: 'rgba(0, 0, 0, 0.7)',
//     })

//     var opt = {
//         song_name: parms.options.songName,
//         song_cover: parms.options.songCover,
//         song_author: parms.options.songAuthor,
//     }
//     StartDownload(opt).then(result => {
//         status.downloading = false;

//         status.downloadButton = true

//         loading.close()
//         footbarText.value = "回到首页";
//         page.value = DownloadPage;
//     })
// }

</script>