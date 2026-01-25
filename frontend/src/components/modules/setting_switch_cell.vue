<template>
    <SettingCell :title="props.title">
        <var-switch v-model="parms" variant @change="handleChange" />
    </SettingCell>
</template>

<script setup lang="ts">
import SettingCell from "./setting_cell.vue"
import { computed } from 'vue'

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

function handleChange() {
    props.onchange?.()
}

</script>

