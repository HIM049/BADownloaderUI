<template>
    <var-paper ripple class="surface-container mb-4">
        <var-collapse v-model="expanded" :offset="true" elevation="0">
            <var-collapse-item :title="title" name="1" class="font-bold text-lg bg-transparent!">
                <slot />
            </var-collapse-item>
        </var-collapse>
    </var-paper>
</template>

<script setup lang="ts">
import { ref } from 'vue'

defineProps<{ title: string }>()

// 默认收起
const expanded = ref<string[]>([])
</script>
