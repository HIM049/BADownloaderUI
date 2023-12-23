<template>
    <el-header>
        <h1>Bili Audio Downloader</h1>
        <el-button size="large"  circle @click="$emit('page-switch')"><el-icon><Setting/></el-icon> </el-button>
    </el-header>
</template>

<style>
header.el-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 40px;
}

header.el-header h1 {
    font-size: var(--el-font-size-extra-large);
}
</style>