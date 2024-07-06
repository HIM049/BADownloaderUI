<template>
    <FramePage title="列表编辑" v-if="CardStatus.ShowList">

        <li v-for="(video, index) in videoList.List" style="list-style-type: none;">
            <var-card :title="video.title" :src="video.Meta.cover" layout="row" image-width="250px" outlines style="margin-bottom: 20px;">
                <template #description>
                    <var-divider />
                    <div style="margin: 0 10px;">
                        <var-chip plain type="info" style="margin-bottom: 5px;">歌曲名称：{{ video.Meta.song_name }}</var-chip>
                        <var-chip plain type="info" style="margin-bottom: 30px;">歌曲作者：{{ video.Meta.author }}</var-chip>
                    </div>
                </template>

                <template #extra>
                    <var-button type="primary" @click="openRightPanel(index)">编辑</var-button>
                </template>
            </var-card>
        </li>
    </FramePage>

    <var-popup position="right" v-model:show="CardStatus.RightPanel" :overlay-style="{backgroundColor: 'rgba(0, 0, 0, 0.2)'}" style=" height: 75%; right: 35px; top: auto; bottom: 35px; border-radius: 8px;">
        <div class="popup-example-block">
            <h3>元数据编辑</h3>
            <div>
                <var-cell><var-input variant="outlined" placeholder="曲名" size="small" v-model="CardStatus.Meta.song_name"/></var-cell>
                <var-cell><var-input variant="outlined" placeholder="歌手" size="small" v-model="CardStatus.Meta.author"/></var-cell>
            </div>
        </div>
        <var-space style="position: absolute; right: 20px; bottom: 20px;">
            <var-button type="primary" :loading="CardStatus.ConfirmBtnLoadig" loading-type="wave" @click="saveVideoMeta">保存</var-button>
        </var-space>
    </var-popup>

</template>

<script setup>
import FramePage from '../modules/frame_page.vue'
import AdditionCard from '../modules/addition_card.vue'
import { reactive, computed, watch, ref } from 'vue'
import { LoadVideoList, SaveVideoList } from '../../../wailsjs/go/main/App'
import { Snackbar, LoadingBar } from '@varlet/ui'

const videoList = ref([])

const CardStatus = reactive({
    RightPanel: false,
    ShowList: false,
    // RightPanelEdit: true,
    ListIndex: 0,
    Meta: reactive({
        song_name: "",
        author: "",
    })
})

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
watch(props, (newValue) => {
    if (newValue.parms.pageIndex == 2) {
        LoadVideoList(parms.value.videoListPath).then(result => {
            videoList.value = result;
            CardStatus.ShowList = true;
        })
    }
})

// 打开右侧面板
function openRightPanel(index) {
    CardStatus.ListIndex = index;
    CardStatus.Meta.song_name = videoList.value.List[index].Meta.song_name;
    CardStatus.Meta.author = videoList.value.List[index].Meta.author;
    CardStatus.RightPanel = true;

}

// 保存列表修改
function saveVideoMeta() {
    videoList.value.List[CardStatus.ListIndex].Meta.song_name = CardStatus.Meta.song_name;
    videoList.value.List[CardStatus.ListIndex].Meta.author = CardStatus.Meta.author;

    SaveVideoList(videoList.value, parms.value.videoListPath).then(result => {
        if (result != null) {
            Snackbar.error("保存失败" + result);
        } else {
            Snackbar.success("保存成功");
        }
    })

    CardStatus.RightPanel = false;
}
</script>

<style>
.popup-example-block {
  padding: 24px;
  width: 280px;
}
</style>