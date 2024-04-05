<template>
    <HeadBar> 
        <var-tabs v-model:active="activePage" style="max-width: 45%; margin: 0px auto;">
            <var-tab><var-icon name="account-circle-outline" /> 个人空间</var-tab>
            <var-tab><var-icon name="content-copy" /> 批量下载</var-tab>
            <var-tab><var-icon name="download-outline" /> 单曲下载</var-tab>
            <var-tab><var-icon name="code-json" /> 软件设置</var-tab>
        </var-tabs>
    </HeadBar>
    <div>
        <div class="page-background"></div>
        <var-tabs-items v-model:active="activePage">
            <var-tab-item>
                <UserSpace/>
            </var-tab-item>

            <var-tab-item>
                <CollectDownload/>
            </var-tab-item>

            <var-tab-item>
                <SingleDownload/>
            </var-tab-item>

            <var-tab-item>
                <SettingPage/>
            </var-tab-item>

        </var-tabs-items>
    </div>
</template>

<script setup>
import HeadBar from '../components/modules/head_bar.vue'
import UserSpace from '../components/user_space.vue'
import CollectDownload from '../components/collect_download.vue'
import SingleDownload from '../components/single_downlaod.vue'
import SettingPage from '../components/setting_page.vue'
import { ref, shallowRef, reactive, computed, watch } from 'vue'

// 分页切换索引
const activePage = ref(1)

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
</style>