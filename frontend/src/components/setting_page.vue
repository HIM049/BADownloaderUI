<template>

    <FramePage title="应用设置" v-if="overload">
        <var-form>
            <var-paper ripple style="background-color: var(--color-primary-container); margin-bottom: 10px">
                <var-collapse v-model="CardStatus.configClass0" :offset="true" elevation="0">
                    <var-collapse-item title="软件行为与外观" name="1" style="background: none; font-size: 1.19em; font-weight: bold;">
                        <var-cell> 
                            主题颜色
                            <template #extra>
                                
                                <var-select variant="outlined" size="small" placeholder="主题色" v-model="config.theme" style="width: 150px;" @change="setTheme">
                                    <var-option label="粉色" :value="'lightPink'" />
                                    <var-option label="蓝色" :value="'lightBlue'" />
                                </var-select>
                            </template>
                        </var-cell>

                        <var-cell> 关闭软件后清除缓存
                            <template #extra>
                                <var-switch v-model="config.delete_cache" variant @change="changeCfg" />
                            </template>
                        </var-cell>
                    </var-collapse-item>
                </var-collapse>
            </var-paper>

            <var-paper ripple style="background-color: var(--color-primary-container); margin-bottom: 10px">
                <var-collapse v-model="CardStatus.configClass1" :offset="true" elevation="0">
                    <var-collapse-item title="账号使用" name="1" style="background: none; font-size: 1.19em; font-weight: bold;">
                        <var-cell> 
                            获取内容时使用账号
                            <template #extra>
                                <var-switch v-model="config.Account.use_account" variant @change="setUseAccount"/>
                            </template>
                        </var-cell>
                        <var-cell>
                            清除保存的账号信息
                            <template #extra>
                                <var-button type="danger" @click="logoutAccount" :disabled="!config.Account.is_login">退出登录</var-button>
                            </template>
                        </var-cell>
                    </var-collapse-item>
                </var-collapse>
            </var-paper>

            <var-paper ripple style="background-color: var(--color-primary-container); margin-bottom: 10px">
                <var-collapse v-model="CardStatus.configClass2" :offset="true" elevation="0">
                    <var-collapse-item title="软件下载行为" name="1" style="background: none; font-size: 1.19em; font-weight: bold;">
                        <var-cell> 
                            最大下载线程数
                            <template #extra>
                                <var-counter v-model="config.download_config.download_threads" @change="changeCfg" />
                            </template>
                        </var-cell>        

                        <var-cell> 
                            下载重试次数
                            <template #extra>
                                <var-counter v-model="config.download_config.retry_count" @change="changeCfg" />
                            </template>
                        </var-cell>
                    </var-collapse-item>
                </var-collapse>
            </var-paper>

            <var-paper ripple style="background-color: var(--color-primary-container); margin-bottom: 10px">
                <var-collapse v-model="CardStatus.configClass3" :offset="true" elevation="0">
                    <var-collapse-item title="文件与路径" name="1" style="background: none; font-size: 1.19em; font-weight: bold;">

                    <var-tooltip content="如果您的计算机中安装了 ffmpeg ，可以打开此开关将音频转码为 MP3 格式输出" style="width: 100%; margin-bottom: 10px;">
                        <var-cell> 使用 ffmpeg 转码音频
                            <template #extra>
                                <var-switch v-model="config.file_config.convert_format" variant @change="setConvertFormat" />
                            </template>
                        </var-cell>
                    </var-tooltip>

                    <var-tooltip content="双大括号中的为文件名变量，可以通过自行修改或删除自定义文件名" style="width: 100%; margin-bottom: 10px;" trigger="click">
                        <template #content>
                            
                            <p v-pre>双大括号中的为文件名变量，可以通过自行修改或删除自定义文件名</p>
                            <p v-pre>列表编号{{.ID}} 视频标题{{.Title}} 单集标题{{.Subtitle}} 音频质量{{.Quality}} 格式后缀名{{.Format}}</p>
                        </template>
                        <var-input style="margin: 10px;" variant="outlined" placeholder="文件命名方式" size="small" v-model="config.file_config.file_name_template"
                            :rules="[v => !!v || '该选项不能为空']" @change="changeCfg" />
                    </var-tooltip>
                    

                    <var-input style="margin: 10px" variant="outlined" placeholder="音频保存路径" size="small" v-model="config.file_config.download_path"
                        :rules="[v => !!v || '该选项不能为空']" @change="changeCfg" />

                    <var-input style="margin: 10px" variant="outlined" placeholder="下载缓存路径" size="small" v-model="config.file_config.cache_path"
                        :rules="[v => !!v || '该选项不能为空']" @change="changeCfg" />

                    <var-input style="margin: 10px" variant="outlined" placeholder="视频列表路径" size="small" v-model="config.file_config.videolist_path"
                        :rules="[v => !!v || '该选项不能为空']" @change="changeCfg" />
                    </var-collapse-item>
                </var-collapse>
            </var-paper>

            <var-space direction="column" size="large">
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
import { Dialog, Snackbar } from '@varlet/ui'

const changeCfg = ref(null) // 修改设置时的响应
const overload = ref(false) // 是否完成页面加载

const CardStatus = reactive({
    configClass0: [],
    configClass1: [],
    configClass2: [],
    configClass3: [],
})

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

// 切换主题
function setTheme() {
    SaveConfig(config.value).then(result => {
        Dialog('立即重新加载主题？').then(result => {
            if (result == 'confirm') {
                window.location.reload();
            }
            return;
        });
        Snackbar.success("保存成功");
    });
}

// 登出账户
function logoutAccount() {
    Dialog('要退出登录吗？').then(result => {
        if (result == 'confirm') {
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
        return;
    })
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
            config.file_config.convert_format = false;
            Snackbar.warning("未检测到 ffmpeg 安装");
        }
    })
}

</script>

<style>
.var-swipe-item::-webkit-scrollbar {
    display: none;
}
</style>