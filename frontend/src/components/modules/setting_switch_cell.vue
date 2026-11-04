<template>
    <SettingCell :title="title">
        <var-switch v-model="value" variant @change="handleChange" />
    </SettingCell>
</template>

<script setup lang="ts">
import SettingCell from './setting_cell.vue'
import { nextTick } from 'vue'

defineProps<{ title?: string }>()
const value = defineModel<boolean>()
const emit = defineEmits<{ change: [] }>()

// var-switch 先触发 change 再更新值，等值更新后再发出 change
function handleChange() {
    nextTick(() => emit('change'))
}
</script>
