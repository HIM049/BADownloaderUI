<template>
    <el-container direction="vertical">
        <HeadBar v-model="pageNum" />
        <section id="app-function">
            <HomePage v-if="pageNum == 0" @back="pageNum--" @next="pageNum++" v-model:parms="parms" />
            <keep-alive>
                <OptionPage v-if="pageNum == 1" @back="pageNum--" @next="pageNum++" v-model:parms="parms" />
            </keep-alive>
            <EditListPage v-if="pageNum == 2" @back="pageNum--" @next="pageNum++" v-model:parms="parms" />
            <DownloadPage v-if="pageNum == 3" @back="pageNum--" @next="pageNum = 0" v-model:parms="parms" />
            <SettingPage v-if="pageNum == -10"/>
            <SingleDownload v-if="pageNum == -11"/>
            <UserSpace v-if="pageNum == -12" />
        </section>
    </el-container>
</template>

<script setup>
import SettingPage from '../components/setting_page.vue'
import SingleDownload from '../components/single_downlaod.vue'
import UserSpace from '../components/user_space.vue'
import HeadBar from '../components/modules/head_bar.vue'
import HomePage from '../components/home_page.vue'
import OptionPage from '../components/option_page.vue'
import EditListPage from '../components/edit_videolist_page.vue'
import DownloadPage from '../components/download_page.vue'
import { ref, shallowRef, reactive, computed, watch } from 'vue'

const pageNum = ref(0)
// const isSetting = ref(false)
// const isSdown = ref(false)
// const page = shallowRef(HomePage)
// const footbarText = ref("继续")

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

</script>