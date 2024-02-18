<template>
    <StepBar :pageNum="step" />
    <el-main style="display: flex; justify-content: center;">
        <div>
            <li v-for="(video, index) in videoList" style="list-style-type: none;">
                <!-- {{ index }} {{ video.title }} -->
                <el-card class="video-card">
                    <template #header>
                        <div class="card-header">
                            <span>{{ video.title }}</span>
                            <el-button class="button" text>{{ video.bvid }}</el-button>
                        </div>
                    </template>
                    <img :src="video.Meta.cover" style="width: 200px;">
                    <el-form label-position="right" style="width: 70%;" :model="video">
                        <el-form-item label="曲名">
                            <el-input v-model="video.Meta.song_name" />
                        </el-form-item>
                        <el-form-item label="歌手">
                            <el-input v-model="video.Meta.author" />
                        </el-form-item>
                        <!-- <el-form-item label="Activity form">
                            <el-input v-model="video" />
                        </el-form-item> -->
                    </el-form>
                    <template #footer>
                        <el-button class="button" type="success" @click="saveVideoList" plain>Save</el-button>
                    </template>
                </el-card>
            </li>
        </div>
    </el-main>
    <FootBar :status="status" text="开始下载" @back="$emit('back')" @next="$emit('next')" />
</template>

<script setup>
import { reactive, computed, onMounted, ref } from 'vue'
import { MakeAndSaveList, GetVideoList, SaveVideoList } from '../../wailsjs/go/main/App'
import StepBar from '../components/modules/step_bar.vue'
import FootBar from '../components/modules/footer.vue'

const props = defineProps(['parms'])
const emit = defineEmits(['back', 'next'])
const parms = computed({
    get() {
        return props.parms
    },
    set(parms) {
        emit('update:parms', parms)
    }
})


const videoList = ref([])
const step = 2 // 进度条步数

// 底栏状态
const status = reactive({
    showBack: true,
    showNext: true,
    allowBack: true,
    allowNext: false,
})

// 页面挂载时加载列表内容
// 创建视频列表
onMounted(() => {

    const loading = ElLoading.service({
        lock: true,
        text: '正在创建列表',
        background: 'rgba(0, 0, 0, 0.7)',
    })

    MakeAndSaveList(props.parms.favListID, Number(props.parms.options.downCount), props.parms.options.downPart).then(result => {
        if (result != null) {
            // 创建失败
            loading.close()
            ElMessage.error(result);
        } else {
            // 创建成功
            ElMessage.success("创建完成");
        }
        loading.close()

        // 修改组件状态
        status.allowNext = true

        GetVideoList().then(result => {
            videoList.value = result
        })
    })

})

function saveVideoList() {
    SaveVideoList(videoList.value).then(result => {
        if (result != null) {
            ElMessage.error("保存失败" + result);
        } else {
            ElMessage.success("保存成功");
        }
    })
}
</script>

<style>
/* 列表卡片样式 */
.video-card {
    width: 60%;
    margin: 0 auto;
    margin-bottom: 10px;
}

.card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.el-card__body {
    display: flex;
    justify-content: space-between;
}

.el-card__body img {
    border-radius: var(--el-border-radius-small);
}
</style>

