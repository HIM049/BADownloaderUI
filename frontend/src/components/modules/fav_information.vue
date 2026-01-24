<template>
    <div class="w-[60%] mx-auto my-[50px] p-5 border-2" style="background: var(--el-bg-color); border-radius: var(--el-border-radius-base); border-color: var(--el-border-color-dark);">
        <h3>收藏夹信息</h3>
        <div class="flex">
            <img class="w-[45%]" :src="resp.cover" style="border-radius: var(--el-border-radius-small);">
            <div class="px-5">
                <h3>{{ resp.title }}</h3>
                <h4>{{ resp.count }} 个内容</h4>
                <div class="flex items-center">
                    <img :src="resp.up_avatar" class="w-[18%]" style="border-radius: var(--el-border-radius-circle);">
                    <h4>{{ resp.up_name }}</h4>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
defineProps(['resp'])
</script>