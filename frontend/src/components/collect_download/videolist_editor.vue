<template>
    <FramePage title="列表编辑" style="width: 50%; margin: 0 auto;" v-if="showList">

        <li v-for="(video, index) in videoList" style="list-style-type: none;">
            <var-card :title="video.title" :src="video.Meta.cover" layout="row" outlines style="margin-bottom: 20px;">
                <template #description>
                    <var-divider />
                    <div>
                        <var-cell><var-input variant="outlined" placeholder="曲名" size="small" v-model="video.Meta.song_name" @change="saveVideoList" /></var-cell>
                        <var-cell><var-input variant="outlined" placeholder="歌手" size="small" v-model="video.Meta.author" @change="saveVideoList" /></var-cell>
                    </div>
                </template>
            </var-card>
        </li>
    </FramePage>
</template>

<script setup>
import FramePage from '../modules/frame_page.vue'
import { reactive, computed, watch, ref } from 'vue'
import { GetVideoList, SaveVideoList } from '../../../wailsjs/go/main/App'
import { Snackbar, LoadingBar } from '@varlet/ui'

const videoList = ref([])
const showList = ref(false)

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
            showList.value = true
        })
    }
})

// 保存视频列表
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