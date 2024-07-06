<template>

    <FramePage title="应用设置" v-if="overload">
        <var-form>
            <var-space direction="column" size="large">
                <var-cell> 
                    使用账号获取内容
                    <template #extra>
                        <var-switch v-model="config.Account.use_account" variant @change="setUseAccount"/>
                    </template>
                </var-cell>
                <var-cell>
                    清除保存的账号信息
                    <template #extra>
                        <var-button type="danger" @click="clearAccount" :disabled="!config.Account.is_login">退出登录</var-button>
                    </template>
                </var-cell>
                <var-divider />

                <var-tooltip content="如果您的计算机中安装了 ffmpeg ，可以打开此开关将音频转码为 MP3 格式输出" style="width: 100%;">
                    <var-cell> 使用 ffmpeg 转码音频
                        <template #extra>
                            <var-switch v-model="config.convert_format" variant @change="setConvertFormat" />
                        </template>
                    </var-cell>
                </var-tooltip>
                
                <var-cell> 关闭软件后清除缓存
                    <template #extra>
                        <var-switch v-model="config.delete_cache" variant @change="changeCfg" />
                    </template>
                </var-cell>

                <var-divider />
                
                <var-cell> 
                    最大下载线程数
                    <template #extra>
                        <var-counter v-model="config.download_threads" @change="changeCfg" />
                    </template>
                </var-cell>        

                <var-cell> 
                    下载重试次数
                    <template #extra>
                        <var-counter v-model="config.retry_count" @change="changeCfg" />
                    </template>
                </var-cell>
                <var-divider />

                <var-input variant="outlined" placeholder="音频保存路径" size="small" v-model="config.download_path"
                    :rules="[v => !!v || '该选项不能为空']" @change="changeCfg" />

                <var-input variant="outlined" placeholder="下载缓存路径" size="small" v-model="config.cache_path"
                    :rules="[v => !!v || '该选项不能为空']" @change="changeCfg" />

                <var-input variant="outlined" placeholder="视频列表路径" size="small" v-model="config.videolist_path"
                    :rules="[v => !!v || '该选项不能为空']" @change="changeCfg" />

                <var-divider />

                <var-space justify="flex-end">
                    <var-button type="danger" @click="refreshConfig">重置设置</var-button>
                    <var-button type="success" @click="changeCfg">保存更改</var-button>
                </var-space>
            </var-space>
        </var-form>
    </FramePage>
</template>

<script setup>
import FramePage from '../components/modules/frame_page.vue'
import { reactive, ref, onMounted } from 'vue'
import { LoadConfig, SaveConfig, RefreshConfig, Checkffmpeg } from '../../wailsjs/go/main/App'
import { Snackbar } from '@varlet/ui'

const changeCfg = ref(null) // 修改设置时的响应
const overload = ref(false) // 是否完成页面加载

onMounted(() => {
    loadConfig();
    setTimeout(() => {
        changeCfg.value = saveConfig
    }, 100)
})

// 设置内容
const config = ref([])

// 读取配置文件
function loadConfig() {
    LoadConfig().then(result => {
        config.value = result;
        overload.value = true;
    })
}

// 保存配置文件
function saveConfig() {
    setTimeout(function () {
        SaveConfig(config.value).then(result => {
            Snackbar.success("保存成功");
        })
    },100);
}

// 重置配置文件
function refreshConfig() {
    RefreshConfig().then(result => {
        loadConfig();
    })
    Snackbar.success("已重置配置文件");
}

// 清除账户信息
function clearAccount() {
    config.value.Account.is_login = false;
    config.value.Account.use_account = false;
    config.value.Account.sessdata = "";
    config.value.Account.bili_jct = "";
    config.value.Account.dede_user_id = "";
    config.value.Account.dede_user_id__ck_md5 = "";
    config.value.Account.sid = "";
    saveConfig();
    setTimeout(() => {
        Snackbar.success("已退出登录");
    }, 110);
    setTimeout(() => {
        window.location.reload();
    }, 1000);
}

// 开启账户开关时检查登录
function setUseAccount() {
    setTimeout(() => {
            
        if (config.value.Account.is_login) {
            saveConfig();
        } else {
            config.value.Account.use_account = false;
            Snackbar.warning("您暂未登录");
        }
    }, 0);
}

// 校验是否存在 ffmpeg
function setConvertFormat() {
    Checkffmpeg().then(result => {
        if (result) {
            saveConfig();
        } else {
            config.convert_format = false;
            Snackbar.warning("未检测到 ffmpeg 安装");
        }
    })
}

</script>