<template>
    <SettingCell :title="props.title">
        <var-switch v-model="parms" variant @change="onchange()" />
    </SettingCell>
</template>

<script setup lang="ts">
import SettingCell from "./setting_cell.vue"
import { computed } from 'vue'

const props = defineProps({
    title: String,
    parms: Boolean,
    onchange: Function,
})
const emit = defineEmits(['update:parms'])

const parms = computed({
    get() {
        return props.parms
    },
    set(parms) {
        emit('update:parms', parms)
    }
})

</script>

<style>
.add-body {
    width: 60%;
    margin: 0px auto 20px auto;;
}
</style>