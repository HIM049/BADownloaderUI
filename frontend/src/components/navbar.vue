<script setup>
</script>

<template>
    <!-- 导航栏 -->
    <nav id="navbar" class="navbar">
      <h2>Bili Audio Downloader</h2>
      <!-- <button "><el-icon size="30"><Setting /></el-icon></button> -->
      <el-button size="large"  circle @click="$emit('page-switch')"><el-icon><Setting style="width: 30px;" /></el-icon> </el-button>
    </nav>
</template>

<style>
.navbar {
    display: flex;
    justify-content: space-between;
    margin: 20px 40px;
    align-items: center;
}
</style>