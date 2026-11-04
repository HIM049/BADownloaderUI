<template>
    <div class="media-card flex overflow-hidden rounded-lg">
        <div class="relative w-[224px] shrink-0 media-card__cover" :style="{ aspectRatio: aspect }">
            <img v-if="cover" :src="cover" class="w-full h-full object-cover transition-[filter,opacity]" :class="{ 'grayscale opacity-50': muted }">
            <div v-if="muted" class="absolute inset-0 flex items-center justify-center">
                <span class="px-3 py-1 rounded-full bg-black/55 text-white text-sm">{{ mutedLabel }}</span>
            </div>
        </div>

        <div class="flex flex-col grow min-w-0 px-4 py-3">
            <div class="transition-opacity" :class="{ 'opacity-50': muted }">
                <p class="truncate text-base font-bold" :title="title">{{ title }}</p>
                <p v-for="item in items" :key="item.label" class="truncate text-sm mt-1.5 media-card__item" :title="String(item.value)">
                    {{ item.label }}：{{ item.value }}
                </p>
            </div>
            <div class="mt-auto flex items-center justify-end gap-2.5">
                <slot name="actions" />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

// 带封面的信息卡片，封面按固定比例裁剪使卡片高度一致；mutedLabel 非空时显示为置灰状态
const props = withDefaults(defineProps<{
    title: string
    cover?: string
    items?: { label: string, value: string | number }[]
    aspect?: string
    mutedLabel?: string
}>(), {
    cover: '',
    items: () => [],
    aspect: '16 / 10',
    mutedLabel: '',
})

const muted = computed(() => props.mutedLabel !== '')
</script>

<style scoped>
.media-card {
    background: var(--card-background);
    border: 1px solid rgba(255, 255, 255, 0.6);
    color: var(--color-text);
}

.media-card__cover {
    background: var(--color-surface-container);
}

.media-card__item {
    color: var(--color-on-surface-variant);
}
</style>
