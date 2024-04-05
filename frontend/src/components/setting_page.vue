<template>
    <body class="page-body">
        <var-paper :elevation="2" >
            <h2>应用设置</h2>
            <var-form >
                <var-space direction="column" size="large">
                    <var-input variant="outlined" 
                    placeholder="音频保存路径" 
                    size="small" 
                    v-model="config.download_path" 
                    :rules="[v => !!v || '该选项不能为空']"
                    @change="saveConfig" 
                    />

                    <var-input variant="outlined" 
                    placeholder="下载缓存路径" 
                    size="small" 
                    v-model="config.cache_path" 
                    :rules="[v => !!v || '该选项不能为空']"
                    @change="saveConfig" 
                    />

                    <var-input variant="outlined" 
                    placeholder="视频列表路径" 
                    size="small" 
                    v-model="config.videolist_path" 
                    :rules="[v => !!v || '该选项不能为空']"
                    @change="saveConfig" 
                    />

                    <label>最大线程数</label>
                    <var-counter v-model="config.download_threads"/>

                    <label>下载重试次数</label>
                    <var-counter v-model="config.retry_count"/>

                    <var-space justify="flex-end">
                        <var-button type="danger" @click="refreshConfig">重置设置</var-button>
                        <var-button type="primary" @click="loadConfig">放弃更改</var-button>
                        <var-button type="success" @click="saveConfig">保存更改</var-button>
                    </var-space>
                </var-space>
            </var-form>
        
        </var-paper>
    </body>
</template>

<script setup>
import FootBar from '../components/modules/footer.vue'
import { reactive, ref, onMounted } from 'vue'
import { LoadConfig, SaveConfig, RefreshConfig } from '../../wailsjs/go/main/App'
import { Snackbar } from '@varlet/ui'

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
    Snackbar.success("保存成功")
}

// 重置配置文件
function refreshConfig() {
    RefreshConfig().then(result => {
        loadConfig()
    })
    Snackbar.success("已重置配置文件")
}

</script>