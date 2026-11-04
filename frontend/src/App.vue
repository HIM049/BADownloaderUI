<template>
    <MainPage />
</template>

<script setup>
import MainPage from './components/main_page.vue'
import { onMounted } from 'vue'
import { Dialog, Snackbar } from '@varlet/ui'
import { Checkffmpeg, InstallFfmpegByWinget } from '../wailsjs/go/wails_api/WailsApi'
import './App.css'

// 启动时检查 FFmpeg 环境
onMounted(() => {
    Checkffmpeg().then(result => {
        if (!result) {
            Dialog({
                title: '缺少必备环境',
                message: '本软件依赖 FFmpeg 进行音频元数据和转码处理，你需要安装该程序才能让软件正常工作',
                closeOnClickOverlay: false,
                confirmButtonText: '知道了',
                cancelButtonText: '使用 Winget 安装',
                onCancel: () => {
                    Snackbar.loading('正在安装 FFmpeg，请稍候...');
                    InstallFfmpegByWinget().then(() => {
                        Snackbar.clear();
                        Snackbar.success('安装成功，重启软件后生效');
                    }).catch(err => {
                        Snackbar.clear();
                        Snackbar.error('安装失败：' + err);
                    });
                },
            });
        }
    });
});

</script>

<style>
div.var-paper {
    padding: 20px;
    border-radius: 20px;
}
</style>