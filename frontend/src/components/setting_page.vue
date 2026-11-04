<template>
    <FramePage title="应用设置" v-if="overload">
        <var-form>
            <SettingGroup title="软件行为与外观">
                <SettingCell title="主题颜色">
                    <var-select variant="outlined" size="small" placeholder="主题色" v-model="config.theme" class="w-[150px]" @change="setTheme">
                        <var-option label="粉色" :value="'lightPink'" />
                        <var-option label="蓝色" :value="'lightBlue'" />
                    </var-select>
                </SettingCell>
                <CellSwitch title="关闭软件后清除缓存" v-model="config.delete_cache" @change="handleSave" />
            </SettingGroup>

            <SettingGroup title="账号使用">
                <CellSwitch title="获取内容时使用账号" v-model="config.account.use_account" @change="setUseAccount" />
                <SettingCell title="清除保存的账号信息">
                    <var-button type="danger" @click="logoutAccount" :disabled="!config.account.is_login">退出登录</var-button>
                </SettingCell>
            </SettingGroup>

            <SettingGroup title="软件下载行为">
                <SettingCell title="最大下载线程数">
                    <var-counter v-model="config.download_config.download_threads" @change="handleSave" />
                </SettingCell>
                <SettingCell title="下载重试次数">
                    <var-counter v-model="config.download_config.retry_count" @change="handleSave" />
                </SettingCell>
            </SettingGroup>

            <SettingGroup title="文件与路径">
                <var-tooltip content="如果您的计算机中安装了 ffmpeg ，可以打开此开关将音频转码为 MP3 格式输出" class="w-full mb-4">
                    <CellSwitch title="使用 ffmpeg 转码音频" v-model="config.file_config.convert_format" @change="setConvertFormat" />
                </var-tooltip>

                <var-tooltip class="w-full mb-4" trigger="click">
                    <template #content>
                        <p v-pre>双大括号中的为文件名变量，可以通过自行修改或删除自定义文件名</p>
                        <p v-pre>列表编号{{.ID}} 视频标题{{.Title}} 单集标题{{.Subtitle}} 音频质量{{.Quality}} 格式后缀名{{.Format}}</p>
                    </template>
                    <var-input class="m-2.5" variant="outlined" placeholder="文件命名方式" size="small" v-model="config.file_config.file_name_template"
                        :rules="[v => !!v || '该选项不能为空']" @change="handleSave" />
                </var-tooltip>

                <div class="flex items-center">
                    <var-input class="m-2.5 w-full" variant="outlined" placeholder="音频保存路径" size="small" readonly v-model="config.file_config.download_path" />
                    <var-button type="primary" @click="setDownloadPathDialog">更改</var-button>
                </div>

                <div class="flex items-center">
                    <var-input class="m-2.5 w-full" variant="outlined" placeholder="下载缓存路径" size="small" readonly v-model="config.file_config.cache_path" />
                    <var-button type="primary" disabled>更改</var-button>
                </div>
            </SettingGroup>

            <var-space justify="flex-end">
                <var-button type="danger" icon-container @click="resetConfig"><var-icon name="refresh" />重置设置</var-button>
                <var-button type="success" icon-container @click="handleSave"><var-icon name="file-document-outline" />保存更改</var-button>
            </var-space>
        </var-form>
    </FramePage>
</template>

<script setup lang="ts">
import FramePage from '../components/modules/frame_page.vue'
import CellSwitch from './modules/setting_switch_cell.vue'
import SettingCell from './modules/setting_cell.vue'
import SettingGroup from './modules/setting_group.vue'
import { ref, onMounted, nextTick } from 'vue'
import { LoadConfig, SaveConfig, ResetConfig, RefreshConfig, SetDownloadPathDialog, Checkffmpeg } from '../../wailsjs/go/wails_api/WailsApi'
import { config as ConfigModels } from '../../wailsjs/go/models'
import { Dialog, Snackbar } from '@varlet/ui'
import { call } from '../utils/call'

// 计数器在值被程序修改时也会触发 change，载入配置期间不保存
const allowSave = ref(false)
const overload = ref(false) // 配置是否已载入

// 设置内容
const config = ref(new ConfigModels.Config())

onMounted(loadConfig)

// 从配置文件重新读取设置
async function loadConfig() {
    allowSave.value = false
    const refreshed = await call(RefreshConfig(), "读取设置失败")
    if (!refreshed.ok) return
    config.value = await LoadConfig()
    overload.value = true
    await nextTick()
    allowSave.value = true
}

async function setDownloadPathDialog() {
    const { ok } = await call(SetDownloadPathDialog(), "设置下载路径失败")
    if (ok) {
        loadConfig()
    }
}

// 表单控件变化时保存
function handleSave() {
    if (allowSave.value) {
        saveConfig()
    }
}

// 保存配置文件，返回是否成功
async function saveConfig(): Promise<boolean> {
    const { ok } = await call(SaveConfig(config.value), "保存失败")
    if (ok) Snackbar.success("保存成功")
    return ok
}

// 重置配置文件
async function resetConfig() {
    const { ok } = await call(ResetConfig(), "重置失败")
    if (ok) {
        loadConfig()
        Snackbar.success("已重置配置文件")
    }
}

// 切换主题
async function setTheme() {
    if (!await saveConfig()) return
    if (await Dialog('立即重新加载主题？') === 'confirm') {
        window.location.reload()
    }
}

// 登出账户
async function logoutAccount() {
    if (await Dialog('要退出登录吗？') !== 'confirm') return
    config.value.account = new ConfigModels.Account({})
    if (!await saveConfig()) return
    Snackbar.success("已退出登录")
    setTimeout(() => window.location.reload(), 1000)
}

// 开启账户开关时检查登录
function setUseAccount() {
    if (config.value.account.is_login) {
        handleSave()
    } else {
        config.value.account.use_account = false
        Snackbar.warning("您暂未登录")
    }
}

// 校验是否存在 ffmpeg
async function setConvertFormat() {
    if (await Checkffmpeg()) {
        handleSave()
    } else {
        config.value.file_config.convert_format = false
        Snackbar.warning("未检测到 ffmpeg，无法启用转码")
    }
}

</script>
