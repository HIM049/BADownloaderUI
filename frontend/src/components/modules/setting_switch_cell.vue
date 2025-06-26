<template>
    <var-cell>
        {{ props.title }}
        <template #extra>
            <var-switch v-model="parms" variant @change="onchange()" />
        </template>
    </var-cell>
</template>

<script setup lang="ts">
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