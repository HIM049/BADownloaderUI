<template>
    <el-main style="display: flex; justify-content: center;">
        <el-button type="primary" size="large" @click="makeUpEditor" >使用系统编辑器打开 <el-icon><Link /></el-icon></el-button>
    </el-main>
</template>

<script setup>
import { MakeUpEditor } from '../../wailsjs/go/main/App'

// 启动文本编辑器（临时解决方案）
function makeUpEditor() {
  MakeUpEditor()
}
</script>

