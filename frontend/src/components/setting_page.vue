<template>

    <body class="page-body">
        <var-paper :elevation="2">
            <h2>应用设置</h2>
            <var-form>
                <var-space direction="column" size="large">
                    <var-input variant="outlined" placeholder="音频保存路径" size="small" v-model="config.download_path"
                        :rules="[v => !!v || '该选项不能为空']" @change="changeCfg" />

                    <var-input variant="outlined" placeholder="下载缓存路径" size="small" v-model="config.cache_path"
                        :rules="[v => !!v || '该选项不能为空']" @change="changeCfg" />

                    <var-input variant="outlined" placeholder="视频列表路径" size="small" v-model="config.videolist_path"
                        :rules="[v => !!v || '该选项不能为空']" @change="changeCfg" />

                    <var-divider />

                    <var-cell> 使用 ffmpeg 转码音频
                        <template #extra>
                            <var-switch v-model="config.convert_format" variant @change="setConvertFormat" />
                        </template>
                    </var-cell>

                    <var-divider />

                    <label>最大线程数</label>
                    <var-counter v-model="config.download_threads" @change="changeCfg" />

                    <label>下载重试次数</label>
                    <var-counter v-model="config.retry_count" @change="changeCfg" />

                    <var-space justify="flex-end">
                        <var-button type="danger" @click="refreshConfig">重置设置</var-button>
                        <var-button type="success" @click="changeCfg">保存更改</var-button>
                    </var-space>
                </var-space>
            </var-form>

        </var-paper>
    </body>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { LoadConfig, SaveConfig, RefreshConfig, Checkffmpeg } from '../../wailsjs/go/main/App'
import { Snackbar } from '@varlet/ui'

onMounted(() => {
    loadConfig()
    setTimeout(() => {
        changeCfg.value = saveConfig
    }, 100)
})

const changeCfg = ref(null)

const config = reactive({
    download_path: "",
    cache_path: "",
    videolist_path: "",
    download_threads: 0,
    retry_count: 0,
})

// 读取配置文件
function loadConfig() {
    LoadConfig().then(result => {
        config.download_path = result.download_path
        config.cache_path = result.cache_path
        config.videolist_path = result.videolist_path
        config.download_threads = result.download_threads
        config.retry_count = result.retry_count
        config.convert_format = result.convert_format
    })
}

// 保存配置文件
function saveConfig() {
    setTimeout(function () {
        SaveConfig(config).then(result => {
            Snackbar.success("保存成功")
        })
    },100);
}

// 重置配置文件
function refreshConfig() {
    RefreshConfig().then(result => {
        loadConfig()
    })
    Snackbar.success("已重置配置文件")
}

// 校验是否存在 ffmpeg
function setConvertFormat() {
    Checkffmpeg().then(result => {
        if (result) {
            saveConfig()
        } else {
            config.convert_format = false
            Snackbar.warning("未检测到 ffmpeg 安装，请检查环境变量")
        }
    })
}

</script>