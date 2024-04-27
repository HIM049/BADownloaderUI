<template>
    <FramePage title="设定下载选项">
        <var-card 
            title="元数据选项" ripple outline style="margin-bottom: 20px"  >
            <template #description>
                <var-divider />
                <div style="margin-left: 20px;">
                    <var-cell> 下载全部
                        <template #extra>
                            <var-switch v-model="downAll" @click="parms.options.downCount = 0" variant />
                        </template>
                    </var-cell>

                    <var-cell> 下载数量
                        <template #extra>
                            <var-counter :min="0" :max="parms.count" v-model="parms.options.downCount" :disabled="downAll"/>
                        </template>
                    </var-cell>

                    <var-cell> 下载分集
                        <template #extra>
                            <var-switch v-model="parms.options.downPart" variant />
                        </template>
                    </var-cell>
                </div>
            </template>
        </var-card>

        <var-card 
            title="元数据选项" ripple outline style="margin-bottom: 20px"  >
            <template #description>
                <var-divider />
                <div style="margin-left: 20px;">
                    <var-cell> 歌曲名称
                        <template #extra>
                            <var-switch v-model="parms.options.songName" @click.stop variant />
                        </template>
                    </var-cell>

                    <var-cell> 歌曲封面
                        <template #extra>
                            <var-switch v-model="parms.options.songCover" @click.stop variant />
                        </template>
                    </var-cell>

                    <var-cell> 歌曲作者
                        <template #extra>
                            <var-switch v-model="parms.options.songAuthor" @click.stop variant />
                        </template>
                    </var-cell>
                </div>
            </template>
        </var-card>

    </FramePage>
</template>

<script setup>
import FramePage from '../modules/frame_page.vue'
import { computed, ref, reactive } from 'vue'

const props = defineProps(['parms', 'status'])
const emit = defineEmits(['update:parms', 'update:status'])

const parms = computed({
    get() {
        return props.parms
    },
    set(parms) {
        emit('update:parms', parms)
    }
})
// 下载全部开关状态
const downAll = ref(true)

</script>