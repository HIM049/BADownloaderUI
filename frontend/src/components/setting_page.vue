<template>
    <el-main style="width: 75%; margin: 0 auto;">
        <el-row class="button-board">            
            <el-button type="success" @click="saveConfig" plain><el-icon><Check /></el-icon> 保存</el-button>
            <el-button type="danger" @click="refreshConfig" plain><el-icon><Refresh /></el-icon> 重置</el-button>
        </el-row>
        <form class="card">
            <h3>软件选项</h3>
            
            <el-form label-position="right" label-width="100px">
                <el-form-item label="音频保存路径">
                    <el-input v-model="config.download_path" @change="saveConfig" />
                </el-form-item>
                <el-form-item label="下载缓存路径">
                    <el-input v-model="config.cache_path" @change="saveConfig" />
                </el-form-item>
                <el-form-item label="视频列表路径">
                    <el-input v-model="config.videolist_path" @change="saveConfig" />
                </el-form-item>
                <el-form-item label="最大线程数">
                    <el-input-number v-model="config.download_threads" :min="1" @change="saveConfig" />
                </el-form-item>
                <el-form-item label="下载重试次数">
                    <el-input-number v-model="config.retry_count" :min="1" @change="saveConfig" />
                </el-form-item>
            </el-form>
        </form>
    </el-main>
    <FootBar :status="status" text="" @back="$emit('back')" @next="$emit('next')" />
</template>

<script setup>
import FootBar from '../components/modules/footer.vue'
import { reactive, ref, onMounted } from 'vue'
import { LoadConfig, SaveConfig, RefreshConfig } from '../../wailsjs/go/main/App'

// 底栏状态
const status = reactive({
    showBack: true,
    showNext: false,
    allowBack: true,
    allowNext: false,
})

onMounted(() => {
    loadConfig()
})

const config = reactive({
    download_path: "",
    cache_path: "",
    videolist_path: "",
    download_threads: 0,
    retry_count: 0,
})

function loadConfig() {
    LoadConfig().then(result => {
        config.download_path     = result.download_path
        config.cache_path        = result.cache_path
        config.videolist_path    = result.videolist_path
        config.download_threads  = result.download_threads
        config.retry_count       = result.retry_count
    })
}

function saveConfig() {
    config.download_threads = parseInt(config.download_threads)
    config.retry_count = parseInt(config.retry_count)
    SaveConfig(config)
    ElMessage.success("保存成功")
}

function refreshConfig() {
    RefreshConfig().then(result => {
        loadConfig()
    })
    
    ElMessage.warning("已重置配置文件")
}

</script>

<style>
.button-board {
    padding-bottom: 10px;
}
</style>