<template>
    <div class="information-card">
        <h3>收藏夹信息</h3>
        <div class="favlist-box" style="display: flex">
            <img class="cover" :src="resp.cover">
            <div style="padding: 0 20px;">
                <h3>{{ resp.title }}</h3>
                <h4>{{ resp.count }} 个内容</h4>
                <div class="up-box">
                    <img :src="resp.up_avatar">
                    <h4>{{ resp.up_name }}</h4>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
defineProps(['resp'])
</script>

<style>
.information-card {
    background: var(--el-bg-color);
    width: 60%;
    margin: 50px auto;
    padding: 20px;
    border-radius: var(--el-border-radius-base);
    border: 2px solid var(--el-border-color-dark);
}

.information-card .cover {
    width: 45%;
    border-radius: var(--el-border-radius-small);
}

.information-card .up-box {
    display: flex;
    align-items: center;
}

.information-card .up-box img {
    width: 18%;
    border-radius: var(--el-border-radius-circle);
}
</style>