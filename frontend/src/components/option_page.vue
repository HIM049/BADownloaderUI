<template>
    <StepBar :pageNum="step" />
    <el-main>
        <!-- 下载设定 -->
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

            <el-form class="option-card" label-position="right" label-width="100px">
                <h3>元数据选项</h3>

                <el-form-item label="歌曲名称">
                    <el-switch v-model="parms.options.songName" style="--el-switch-on-color: #13ce66;" />
                </el-form-item>
                <el-form-item label="歌曲封面">
                    <el-switch v-model="parms.options.songCover" style="--el-switch-on-color: #13ce66;" />
                </el-form-item>
                <el-form-item label="歌曲作者">
                    <el-switch v-model="parms.options.songAuthor" style="--el-switch-on-color: #13ce66;" />
                </el-form-item>

            </el-form>
        </div>
    </el-main>
    <FootBar :status="status" text="生成列表" @back="$emit('back')" @next="$emit('next')" />
</template>

<script setup>
import StepBar from '../components/modules/step_bar.vue'
import FootBar from '../components/modules/footer.vue'
import { computed, ref, reactive } from 'vue'

// 模块参数
const props = defineProps(['parms'])
const emit = defineEmits(['update:parms', 'back', 'next'])
// 模块参数响应
const parms = computed({
    get() {
        return props.parms
    },
    set(parms) {
        emit('update:parms', parms)
    }
})

const step = 1 // 进度条步数
const downAll = ref(true)

// 底栏状态
const status = reactive({
    showBack: true,
    showNext: true,
    allowBack: true,
    allowNext: true,
})
</script>

<style>
#options {
    display: flex;
    justify-content: center;
}

.option-card {
    width: 250px;
    padding: 30px;
    margin: 10px;
    background-color: var(--el-bg-color);
    border: 2px solid var(--el-border-color-dark);
}
</style>