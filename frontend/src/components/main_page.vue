<template>
    <HeadBar> 
        <var-tabs v-model:active="activePage" style="max-width: 45%; margin: 0px auto;">
            <var-tab><var-icon name="account-circle-outline" /> 个人空间</var-tab>
            <var-tab><var-icon name="content-copy" /> 批量下载</var-tab>
            <var-tab><var-icon name="code-json" /> 软件设置</var-tab>
        </var-tabs>
    </HeadBar>
    <div>
        <div id="page-background" class="page-background" style="overflow: auto;">
            <var-tabs-items v-model:active="activePage" style="height: 100%;">
                <var-tab-item>
                    <UserSpace/>
                </var-tab-item>

                <var-tab-item id="scroll-box">
                    <CollectDownload/>
                </var-tab-item>

                <var-tab-item>
                    <SettingPage/>
                </var-tab-item>

            </var-tabs-items>
        </div>
    </div>
</template>

<script setup>
import HeadBar from '../components/modules/head_bar.vue'
import UserSpace from '../components/user_space.vue'
import CollectDownload from '../components/collect_download.vue'
import SettingPage from '../components/setting_page.vue'
import { ref } from 'vue'
import { EventsOn } from '../../wailsjs/runtime'
import { Snackbar } from '@varlet/ui'

// 分页切换索引
const activePage = ref(1)

// 全局错误提示
EventsOn("error", (err) => {
    Snackbar.warning(err)
})
</script>

<style>
/* 背景框 */
.page-background {
    background: var(--color-body);
    position: absolute;
    width: calc(100% - 20px);
    height: calc(100% - 130px);
    border-radius: 28px;
}

/* 修复 input 刷新后不显示提示文字 */
.var-input label {
    max-width: none;
}

/* 页面滚动 */
.var-swipe-item {
    overflow: auto;
}
</style>