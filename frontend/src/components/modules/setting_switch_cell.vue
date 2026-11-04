<template>
    <SettingCell :title="props.title">
        <var-switch v-model="parms" variant @change="handleChange" />
    </SettingCell>
</template>

<script setup lang="ts">
import SettingCell from "./setting_cell.vue"
import { computed, nextTick } from 'vue'

const props = defineProps<{
    title?: string
    parms?: boolean
    onchange?: () => void
}>()

const emit = defineEmits(['update:parms'])

const parms = computed({
    get() {
        return props.parms
    },
    set(parms) {
        emit('update:parms', parms)
    }
})

// var-switch 先触发 change 再更新值，等值更新后再回调
function handleChange() {
    nextTick(() => props.onchange?.())
}

</script>

