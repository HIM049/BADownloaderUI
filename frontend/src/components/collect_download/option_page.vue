<template>
    <!-- <el-main>
        <div id="options">
            <el-form class="option-card" label-position="right" label-width="100px">
                <h3>下载选项</h3>

                <el-form-item label="下载全部">
                    <el-switch v-model="downAll" style="--el-switch-on-color: #13ce66;"
                        @click="parms.options.downCount = 0" />
                </el-form-item>

                <el-form-item label="下载数量">
                    <el-input-number v-model="parms.options.downCount" :min="0" :max="parms.count" :disabled="downAll" />
                </el-form-item>

                <el-form-item label="下载分集">
                    <el-switch v-model="parms.options.downPart" style="--el-switch-on-color: #13ce66;" />
                </el-form-item>

            </el-form>

        </div>
    </el-main> -->
    <FramePage title="设定下载选项">
        <var-card 
            title="元数据选项" ripple outline style="margin-bottom: 20px"  >
            <template #description>
                <var-divider />
                <div style="margin-left: 20px;">
                    <var-cell> 下载全部
                        <template #extra>
                            <var-switch v-model="downAll" @click="parms.options.downCount = 0" />
                        </template>
                    </var-cell>

                    <var-cell> 下载数量
                        <template #extra>
                            <var-counter :min="0" :max="parms.count" v-model="parms.options.downCount"/>
                        </template>
                    </var-cell>

                    <var-cell> 下载分集
                        <template #extra>
                            <var-switch v-model="parms.options.downPart" />
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
                            <var-switch v-model="parms.options.songName" @click.stop />
                        </template>
                    </var-cell>

                    <var-cell> 歌曲封面
                        <template #extra>
                            <var-switch v-model="parms.options.songCover" @click.stop />
                        </template>
                    </var-cell>

                    <var-cell> 歌曲作者
                        <template #extra>
                            <var-switch v-model="parms.options.songAuthor" @click.stop />
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
// 下载全部（无作用）
const downAll = ref(true)

</script>