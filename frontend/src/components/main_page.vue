<template>
    <el-container direction="vertical">
        <HeadBar @page-switch="isSetting = !isSetting" />
        <transition name="el-fade-in-linear">
            <el-steps :space="150" :active="pageNum" finish-status="success" align-center style="justify-content: center;"
                v-show="status.showStep">
                <el-step title="查找收藏夹" />
                <el-step title="编辑个性化选项" />
                <el-step title="编辑列表内容" />
                <el-step title="下载完成" />
            </el-steps>
        </transition>
        <section id="app-function">
            <keep-alive>
                <component :is="page" v-model:buttonStatus="status" v-model:parms="parms"></component>
            </keep-alive>
        </section>
        <FootBar :status="status" :text="footbarText" @back="pageNum--" @next="pageNum++" />
    </el-container>
</template>

<script setup>
import SettingPage from '../components/setting_page.vue'
import HeadBar from '../components/modules/head_bar.vue'
import HomePage from '../components/home_page.vue'
import OptionPage from '../components/option_page.vue'
import EditListPage from '../components/edit_videolist_page.vue'
import FinalPage from '../components/final_page.vue'
import FootBar from '../components/modules/footer.vue'
import { ref, shallowRef, reactive, computed, watch } from 'vue'
import { MakeAndSaveList, StartDownload } from '../../wailsjs/go/main/App'

const pageNum = ref(0)
const isSetting = ref(false)
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

// 切换到设置页面
watch(isSetting, (newValue)=>{
    if (newValue) {
        page.value = SettingPage;
        status.showBack = false;
        status.showNext = false;
        status.showStep = false;
    } else {
        page.value = HomePage;
        status.showNext = true;
        pageNum.value = 0;
    }
})

// 根据页面编号切换页面和按钮状态
watch(pageNum, (newValue) => {
    switch (newValue) {
        case 0: // 主页面
            page.value = HomePage;
            status.allowNext = false;
            status.showBack = false;
            footbarText.value = "继续";
            break;
        case 1: // 列表生成选项页
            page.value = OptionPage;
            status.showStep = true;
            status.showBack = true;
            footbarText.value = "生成列表";
            break;
        case 2: // 列表编辑页
            creatVideoList();
            break;
        case 3: // 结束页
            downloadVideoList();
            break;
        default:
            pageNum.value = 0;
            break;

    }
})

// 创建视频列表
function creatVideoList() {
    const loading = ElLoading.service({
        lock: true,
        text: '正在创建列表',
        background: 'rgba(0, 0, 0, 0.7)',
    })

    MakeAndSaveList(parms.favListID, Number(parms.options.downCount), parms.options.downPart).then(result => {
        if (result != null) {
            // 创建失败
            loading.close()
            ElMessage.error(result);
        } else {
            // 创建成功
            loading.close()
            ElMessage.success("创建完成");
        }

        // 修改组件状态
        status.makeingList = false;
        status.makeListButton = true
        status.editlistButton = true;
        status.downloadButton = true;

        page.value = EditListPage;
        footbarText.value = "开始下载";
    })
}

// 下载视频列表
function downloadVideoList() {

    // 加载页面
    const loading = ElLoading.service({
        lock: true,
        text: '正在下载音频',
        background: 'rgba(0, 0, 0, 0.7)',
    })

    var opt = {
        song_name: parms.options.songName,
        song_cover: parms.options.songCover,
        song_author: parms.options.songAuthor,
    }
    StartDownload(opt).then(result => {
        status.downloading = false;

        status.downloadButton = true

        loading.close()
        footbarText.value = "回到首页";
        page.value = FinalPage;
    })
}

</script>