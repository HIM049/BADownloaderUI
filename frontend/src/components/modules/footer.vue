<template>
    <el-footer>
        <el-button type="primary" size="large" @click="$emit('back')" :style="{ 'visibility': showBack }" :disabled="!props.status.allowBack"><el-icon><ArrowLeftBold /></el-icon> 返回</el-button>
        <el-button type="primary" size="large" @click="$emit('next')" :style="{ 'visibility': showNext }" :disabled="!props.status.allowNext" :loading="props.status.loadingNext">{{ props.text }} <el-icon><ArrowRightBold /></el-icon></el-button>
    </el-footer>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps(['status', 'text'])

const showBack = computed(() => {
  return props.status.showBack ? 'visible' : 'hidden'
})

const showNext = computed(() => {
  return props.status.showNext ? 'visible' : 'hidden'
})
</script>

<style>
footer.el-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 40px;
    position: fixed;
    width: 100%;
    bottom: 0;
}
</style>