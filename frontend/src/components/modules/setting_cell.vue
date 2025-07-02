<template>
    <var-cell>
        {{ title }}
        <template #extra>
            <slot/>
        </template>
    </var-cell>
</template>

<script setup lang="ts">

defineProps({
    title: String,
})
</script>