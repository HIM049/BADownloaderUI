<template>
    <el-main style="display: flex; justify-content: center;">
        <!-- <el-button type="primary" size="large" @click="makeUpEditor" >使用系统编辑器打开 <el-icon><Link /></el-icon></el-button> -->
        <!-- <el-button type="primary" size="large" @click="getvideolist" >获取 <el-icon><Link /></el-icon></el-button> -->
        <div>
            <li v-for="(video, index) in videoList" style="list-style-type: none;">
                <!-- {{ index }} {{ video.title }} -->
                <el-card class="video-card" style="margin-bottom: 10px;">
                    <template #header>
                        <div class="card-header">
                            <span>{{ video.title }}</span>
                            <el-button class="button" text>{{ video.bvid }}</el-button>
                        </div>
                    </template>
                    <img :src="video.cover" style="width: 200px;">
                    <el-form label-position="right" style="width: 70%;" :model="video" >
                        <el-form-item label="曲名">
                            <el-input v-model="video.song_name" />
                        </el-form-item>
                        <el-form-item label="歌手">
                            <el-input v-model="video.author" />
                        </el-form-item>
                        <!-- <el-form-item label="Activity form">
                            <el-input v-model="video" />
                        </el-form-item> -->
                    </el-form>
                    <template #footer>
                        <el-button class="button" type="success" @click="saveVideoList" plain>Save</el-button>
                    </template>
                </el-card>
            </li>
        </div>
    </el-main>
</template>

<script setup>
import { MakeUpEditor, GetVideoList, SaveVideoList } from '../../wailsjs/go/main/App'
import { reactive, computed, onMounted, ref } from 'vue'
const videoList = ref([])


// 启动文本编辑器（临时解决方案）
function makeUpEditor() {
    MakeUpEditor()
}

// 页面挂载时加载列表内容
onMounted(() => {
    GetVideoList().then(result => {
        videoList.value = result
    })
})

function saveVideoList() {
    SaveVideoList(videoList.value).then(result => {
        if (result != null) {
            ElMessage.error("保存失败"+result);
        } else {
            ElMessage.success("保存成功");
        }
    })
}
</script>

<style>
.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}
.el-card__body {
    display: flex;
    justify-content: space-between;
}

.el-card__body img {
    border-radius: var(--el-border-radius-small);
}
</style>

