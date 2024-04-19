<template>
    <!-- <el-main style="display: flex; justify-content: center;">
        <div>
            <li v-for="(video, index) in videoList" style="list-style-type: none;">
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
</el-form>
<template #footer>
                        <el-button class="button" type="success" @click="saveVideoList" plain>Save</el-button>
                    </template>
</el-card>
</li>
</div>
</el-main> -->
    <FramePage title="列表编辑" style="width: 50%; margin: 0 auto;">
        <!-- <var-card outline title="点击按钮生成列表">
            <template #extra>
                <var-space>
                    <var-button type="primary" @click="makeList">生成列表</var-button>
                </var-space>
            </template>
        </var-card> -->

        <li v-for="(video, index) in videoList" style="list-style-type: none;">
            <var-card :title="video.title" :src="video.Meta.cover" layout="row" outlines style="margin-bottom: 20px;">
                <template #extra>
                    <var-space justify="flex-end">
                        <var-button type="primary">编辑</var-button>
                    </var-space>
                </template>
            </var-card>
            <!-- <el-form label-position="right" style="width: 70%;" :model="video">
                    <el-form-item label="曲名">
                        <el-input v-model="video.Meta.song_name" />
                    </el-form-item>
                    <el-form-item label="歌手">
                        <el-input v-model="video.Meta.author" />
                    </el-form-item>
                </el-form> -->
        </li>
    </FramePage>
</template>

<script setup>
import FramePage from '../modules/frame_page.vue'
import { reactive, computed, watch, ref } from 'vue'
import { MakeAndSaveList, GetVideoList, SaveVideoList } from '../../../wailsjs/go/main/App'
import { Snackbar, LoadingBar } from '@varlet/ui'

const videoList = ref([])
const show = ref(true)

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

const status = computed({
    get() {
        return props.status
    },
    set(status) {
        emit('update:status', status)
    }
})

// 检查是否完成列表加载
watch(props.status, (newValue) => {
    if (newValue.loading == false) {
        GetVideoList().then(result => {
            videoList.value = result
        })
    }
})

function saveVideoList() {
    SaveVideoList(videoList.value).then(result => {
        if (result != null) {
            Snackbar.error("保存失败" + result);
        } else {
            Snackbar.success("保存成功");
        }
    })
}
</script>

<style>
.popup-example-block {
  padding: 24px;
  width: 280px;
}
</style>