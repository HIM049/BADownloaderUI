<template>
    <transition name="el-fade-in-linear">
        <el-steps :space="150" :active="pageNum" finish-status="success" align-center style="justify-content: center;">
            <el-step title="查找收藏夹" />
            <el-step title="个性化选项" />
            <el-step title="编辑列表内容" />
            <el-step title="下载" />
        </el-steps>
    </transition>
</template>

<script setup>
const props = defineProps(['pageNum'])
</script>