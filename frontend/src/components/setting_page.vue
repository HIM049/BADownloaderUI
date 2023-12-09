<script setup>
import { reactive, ref, onMounted } from 'vue'
import { LoadConfig, SaveConfig } from '../../wailsjs/go/main/App'

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
        config.download_path = result.download_path
        config.cache_path = result.cache_path
        config.videolist_path = result.videolist_path
        config.download_threads = result.download_threads
        config.retry_count = result.retry_count
    })
}

function saveConfig() {
    config.download_threads = parseInt(config.download_threads)
    config.retry_count = parseInt(config.retry_count)
    SaveConfig(config)
}

</script>

<template>
    <div id="buttonBoard">
        <button class="harf-button" @click.prevent="saveConfig">
            <p>保存</p>
        </button>
        <button class="harf-button" @click.prevent="loadConfig">
            <p>重载</p>
        </button>
    </div>
    <form class="card">
    <h3>软件选项</h3>
        <ol>
            <li class="option-list config-option">
                <label for="downloadPath">音频保存路径</label> <br>
                <input type="text" class="input" id="downloadPath" v-model="config.download_path">
            </li>
            <li class="option-list config-option">
                <label for="cachePath">下载缓存路径</label> <br>
                <input type="text" class="input" id="cachePath" v-model="config.cache_path">
            </li>
            <li class="option-list config-option">
                <label for="videolistPath">视频列表路径</label> <br>
                <input type="text" class="input" id="videolistPath" v-model="config.videolist_path">
            </li>
            <li class="option-list config-option">
                <label for="downloadThreads">最大线程数</label> <br>
                <input type="text" class="input" id="downloadThreads" v-model="config.download_threads">
            </li>
            <li class="option-list config-option">
                <label for="retryCount">下载重试次数</label> <br>
                <input type="text" class="input" id="retryCount" v-model="config.retry_count">
            </li>
        </ol>


    </form>
</template>

<style>
.config-option {
    display: block;
}
.harf-button{
    width: 45%;
    padding: 10px;
    margin-bottom: 10px;
    display: flex;
    justify-content: space-evenly;
    align-items: center;
}
#buttonBoard {
    display: flex;
    justify-content: space-evenly;
}
</style>