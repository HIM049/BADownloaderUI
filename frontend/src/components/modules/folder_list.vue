<template>
    <var-paper class="surface-container">
        <div v-for="folder in folders" :key="folder.id" class="my-1.5 mx-5">
            <var-space justify="space-between" align="center">
                <span class="text-[15px] font-semibold">{{ folder.title }}</span>
                <var-button type="primary" @click="emit('add', folder)"><var-icon name="plus" />添加至列表</var-button>
            </var-space>
        </div>
    </var-paper>
</template>

<script setup lang="ts">
import { biliapi } from '../../../wailsjs/go/models'

defineProps<{ folders: biliapi.FavFolder[] }>()
const emit = defineEmits<{ add: [folder: biliapi.FavFolder] }>()
</script>
