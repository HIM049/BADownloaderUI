<template>
    <el-header>
        <h1>Bili Audio Downloader</h1>
        <div>
            <el-button size="large"  circle @click="$emit('switch-sdown')"><el-icon><Download /></el-icon> </el-button>
            <el-button size="large"  circle @click="$emit('switch-setting')"><el-icon><Setting/></el-icon> </el-button>
        </div>
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