<template>
    <el-main>
        <div style="width: 700px; margin: 0 auto 10px auto;">
            <h1 style="font-size: 20px; padding-bottom: 10px;">单曲下载器</h1>
            <el-input v-model="auid" size="large" placeholder="请输入 AUID 的数字部分" class="input-with-select"
                @input="searchSongInformation" clearable>

                <template #prepend>
                    <el-button @click="searchSongInformation"><el-icon>
                            <Search />
                        </el-icon></el-button>
                </template>
            </el-input>
        </div>

        <transition name="el-fade-in-linear">
            <el-card class="video-card" style="width: 600px; margin: 0 auto;" v-show="status.showInf">
                <template #header>
                    <div class="card-header">
                        <span>{{ songInf.title }}</span>
                        <el-button class="button" text>AU{{ songInf.auid }}</el-button>
                    </div>
                </template>
                <img :src="songInf.cover" style="width: 200px;">
                <el-form label-position="right" style="width: 50%;" :model="songInf">
                    <el-form-item label="曲名">
                        <el-input v-model="parms.song_name" />
                    </el-form-item>
                    <el-form-item label="歌手">
                        <el-input v-model="parms.author" />
                    </el-form-item>
                    <el-form-item label="歌曲名称">
                        <el-switch v-model="parms.options.songName" style="--el-switch-on-color: #13ce66;" />
                    </el-form-item>
                    <el-form-item label="歌曲封面">
                        <el-switch v-model="parms.options.songCover" style="--el-switch-on-color: #13ce66;" />
                    </el-form-item>
                    <el-form-item label="歌曲作者">
                        <el-switch v-model="parms.options.songAuthor" style="--el-switch-on-color: #13ce66;" />
                    </el-form-item>
                </el-form>
                <template #footer>
                    <el-button class="button" type="success" @click="audioDownload" plain>开始下载</el-button>
                </template>
            </el-card>
        </transition>

    </el-main>
    <FootBar :status="status" text="" @back="$emit('back')" @next="$emit('next')" />
</template>

<script setup>
import FootBar from '../components/modules/footer.vue'
import { reactive, computed, ref } from 'vue'
import { AudioDownload, SearchSongInformation } from '../../wailsjs/go/main/App'
import { ElMessage } from 'element-plus';

// 底栏状态
const status = reactive({
    showBack: true,
    showNext: false,
    allowBack: true,
    allowNext: false,
    showInf: false,
})

const parms = reactive({
    song_name: "",
    author: "",
    // 下载设置
    options: reactive({
        songName: true,
        songCover: true,
        songAuthor: true,
    })
})

const auid = ref("")

const songInf = reactive({
    auid: "0",
    title: "",
    cover: "",
})

function searchSongInformation() {
    SearchSongInformation(auid.value).then(result => {
        if (result.code == "0") {
            songInf.auid = result.Data.id
            songInf.title = result.Data.title
            songInf.cover = result.Data.cover
            parms.song_name = result.Data.title
            parms.author = result.Data.author

            status.showInf = true
        } else {
            ElMessage.warning("无效的收藏夹");
        }
    })
}

function audioDownload() {
    ElMessage.success("开始下载")
    var opt = {
        song_name: parms.options.songName,
        song_cover: parms.options.songCover,
        song_author: parms.options.songAuthor,
    }
    AudioDownload(opt, auid.value, parms.song_name, parms.author, songInf.title).then(result => {
        ElMessage.success("下载完成")
    })
}
</script>

<style>
.el-card__body {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.el-card__footer {
    display: flex;
    justify-content: flex-end;
}
</style>