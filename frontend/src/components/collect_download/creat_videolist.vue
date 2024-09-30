<template>
    <FramePage title="批量下载">
        <var-space justify="center">                
            <var-button type="primary" @click="creatVideoList" size="large"><var-icon name="plus" />创建新的下载列表</var-button>
            <var-button type="primary" @click="openFile" size="large"><var-icon name="file-document-outline" />打开本地下载列表</var-button>
        </var-space>
    </FramePage>
</template>

<script setup>
import FramePage from '../modules/frame_page.vue'
import { reactive, computed, ref, watch } from 'vue'
import { LoadConfig, CreatVideoList, OpenFileDialog } from '../../../wailsjs/go/main/App'
import { Snackbar } from '@varlet/ui'

const props = defineProps(['parms', 'status'])
const emit = defineEmits(['update:parms', 'update:status', 'nextpage'])

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

// 创建视频列表并保存路径
function creatVideoList() {
    LoadConfig().then(result => {
        parms.value.videoListPath = result.file_config.videolist_path;
        CreatVideoList();
        emit('nextpage');
    })
}

// 选择已有的视频列表并保存路径
function openFile() {
    OpenFileDialog().then(result => {
        if (result == ""){
            Snackbar.warning("未选择文件");
            return
        }
        parms.value.videoListPath = result;
        emit('nextpage');
    })
}
</script>