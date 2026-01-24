<template>
    <FramePage title="列表编辑">
        <!-- Removed manual controls -->
    </FramePage>
    <AdditionCard v-if="CardStatus.ShowList">
        <var-list :finished="CardStatus.finished" :loading="CardStatus.loading" @load="load" :immediate-check="false">
            <li v-for="(video, index) in TaskList.tasks" :key="index" style="list-style-type: none;">
                <var-card :title="video.SongName" :src="video.CoverUrl" layout="row" image-width="250px" outlines v-if="!video.delete" style="margin-bottom: 20px;">
                    <template #description>
                        <var-divider />
                        <div style="margin: 0 10px;">
                            <var-chip plain type="info" style="margin-bottom: 5px;">歌曲名称：{{ video.SongName }}</var-chip>
                            <var-chip plain type="info" style="margin-bottom: 30px;">歌曲作者：{{ video.SongAuthor }}</var-chip>
                        </div>
                    </template>

                    <template #extra>
                        <div style="display: flex; align-items: center;">
                            <var-button type="danger" size="large" round @click="setDeleteState(index)" style="margin-right: 10px;"> <var-icon name="delete" /> </var-button>
                            <var-button type="primary" @click="openRightPanel(index)">编辑</var-button>
                        </div>
                    </template>
                </var-card>

                
                <var-card :title="video.title" outlines style="margin-bottom: 20px; height: 187px;" v-if="video.delete">
                    <template #description>
                        <div style="display: flex; justify-content: center;">
                            
                        <h3>已设为删除</h3>
                        </div>
                    </template>
                    <template #extra>
                        <div style="display: flex; align-items: center;">
                            <var-button type="success" @click="setDeleteState(index)" style="margin-right: 10px;"> <var-icon name="history" />  恢复 </var-button>
                            <var-button type="primary" disabled>编辑</var-button>
                        </div>
                    </template>
                </var-card>
            </li>
            <template #finished>
                <div class="footer-text">已经到底了</div>
            </template>
        </var-list>
        
        <!-- Debug: Manual Load Button -->
        <!-- <var-button block type="primary" @click="load" style="margin-top: 10px;">Manually Load More (Debug)</var-button> -->
        
    </AdditionCard>

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
import { reactive, computed, watch, ref, onMounted } from 'vue'
// import { SaveVideoListTo } from '../../../wailsjs/go/main/App'
import { GetTaskListPage, GetTaskListAll } from '../../../wailsjs/go/wails_api/WailsApi'
import { EventsOn } from '../../../wailsjs/runtime'
import { Snackbar, LoadingBar, Dialog } from '@varlet/ui'

const videoList = ref([])

const TaskList = reactive({
    tasks: [],
    index: 0,
})

const CardStatus = reactive({
    RightPanel: false,
    ShowList: true, // Init to true to show list immediately
    // RightPanelEdit: true,
    ListIndex: 0,
    Meta: reactive({
        song_name: "",
        author: "",
    }),
    loading: false,
    finished: false,
})

const props = defineProps(['parms', 'status'])
const emit = defineEmits(['update:parms', 'update:status', 'refresh'])
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

onMounted(() => {
    load() // Initial load
})

EventsOn('refreshVideoList', () => {
    TaskList.index = 0
    TaskList.tasks = []
    CardStatus.finished = false
    load();
})


// 另存列表
function saveListTo() {
    // Dialog('清理删除项并导出列表？').then(result => {
    //     if (result == 'confirm') {
    //         tidyAndRefresh(() => {
    //             SaveVideoListTo(videoList.value).then(() => {
    //                 Snackbar.success('导出成功');
    //             });
    //         });
    //     }
    //     return;
    // });
}

// 修改视频删除状态
function setDeleteState(index) {
    // videoList.value.List[index].delete = !videoList.value.List[index].delete;
    // SaveVideoList(videoList.value, parms.value.videoListPath).then(result => {
    //     if (result != null) {
    //         Snackbar.error("保存失败" + result);
    //     } else {
    //         Snackbar.success("保存成功");
    //     }
    // })
}

// 打开右侧面板
function openRightPanel(index) {
    CardStatus.ListIndex = index;
    // Note: This logic assumes videoList.value.List is updated, but current code uses TaskList.tasks.
    // Need to check where videoList is populated or update to use TaskList.tasks
    // Assuming TaskList.tasks is the source of truth for display
    if (TaskList.tasks[index].Meta) {
        CardStatus.Meta.song_name = TaskList.tasks[index].Meta.song_name;
        CardStatus.Meta.author = TaskList.tasks[index].Meta.author;
    } else {
         // Fallback if Meta matches original structure
        CardStatus.Meta.song_name = TaskList.tasks[index].SongName;
        CardStatus.Meta.author = TaskList.tasks[index].SongAuthor;
    }
    
    CardStatus.RightPanel = true;

}

// 动态加载数据
function load() {
    console.log("Loading page:", TaskList.index);
    CardStatus.loading = true
    GetTaskListPage(TaskList.index).then(result => {
        console.log("Loaded result:", result ? result.length : 0);
        if (!result || result.length === 0) {
            console.log("Finished");
            CardStatus.finished = true
            CardStatus.loading = false
            return
        }
        
        TaskList.tasks.push(...result)
        TaskList.index++
        CardStatus.loading = false
    }).catch(err => {
        console.log("Load error:", err);
        CardStatus.loading = false
        Snackbar.error("加载失败: " + err)
    })
}


// 保存列表修改
function saveVideoMeta() {
    // Logic needs to be adapted to backend sync
    // videoList.value.List[CardStatus.ListIndex].Meta.song_name = CardStatus.Meta.song_name;
    // videoList.value.List[CardStatus.ListIndex].Meta.author = CardStatus.Meta.author;

    // SaveVideoList(videoList.value, parms.value.videoListPath).then(result => {
    //     if (result != null) {
    //         Snackbar.error("保存失败" + result);
    //     } else {
    //         Snackbar.success("保存成功");
    //     }
    // })

    CardStatus.RightPanel = false;
}
</script>

<style>
.popup-example-block {
  padding: 24px;
  width: 280px;
}
</style>
