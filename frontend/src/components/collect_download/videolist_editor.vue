<template>
    <FramePage title="列表编辑">
        <template #actions>
            <var-button type="primary" @click="saveListTo">另存为</var-button>
        </template>
        <!-- Removed manual controls -->
    </FramePage>
    <AdditionCard v-if="CardStatus.ShowList">
        <var-list :finished="CardStatus.finished" :loading="CardStatus.loading" @load="load" :immediate-check="false">
            <li v-for="(video, index) in TaskList.tasks" :key="index" class="list-none">
                <var-card :title="video.SongName" :src="video.CoverUrl" layout="row" image-width="250px" outlines v-if="!video.delete" class="mb-5">
                    <template #description>
                        <var-divider />
                        <div class="mx-2.5">
                            <var-chip plain type="info" class="mb-1.5">歌曲名称：{{ video.SongName }}</var-chip>
                            <var-chip plain type="info" class="mb-[30px]">歌曲作者：{{ video.SongAuthor }}</var-chip>
                        </div>
                    </template>

                    <template #extra>
                        <div class="flex items-center">
                            <var-button type="danger" size="large" round @click="setDeleteState(index)" class="mr-2.5"> <var-icon name="delete" /> </var-button>
                            <var-button type="primary" @click="openRightPanel(index)">编辑</var-button>
                        </div>
                    </template>
                </var-card>

                
                <var-card :title="video.title" outlines class="mb-5 h-[187px]" v-if="video.delete">
                    <template #description>
                        <div class="flex justify-center">
                            
                        <h3>已设为删除</h3>
                        </div>
                    </template>
                    <template #extra>
                        <div class="flex items-center">
                            <var-button type="success" @click="setDeleteState(index)" class="mr-2.5"> <var-icon name="history" />  恢复 </var-button>
                            <var-button type="primary" disabled>编辑</var-button>
                        </div>
                    </template>
                </var-card>
            </li>
            <template #finished>
                <div class="py-5 text-center text-[#888]">已经到底了</div>
            </template>
        </var-list>
        
        <!-- Debug: Manual Load Button -->
        <!-- <var-button block type="primary" @click="load" style="margin-top: 10px;">Manually Load More (Debug)</var-button> -->
        
    </AdditionCard>

    <var-popup position="right" v-model:show="CardStatus.RightPanel" :overlay-style="{backgroundColor: 'rgba(0, 0, 0, 0.2)'}" class="!h-[75%] !absolute !right-[35px] !bottom-[35px] !rounded-lg !top-auto">
        <div class="p-6 w-[280px]">
            <h3>元数据编辑</h3>
            <div>
                <var-cell><var-input variant="outlined" placeholder="曲名" size="small" v-model="CardStatus.Meta.song_name"/></var-cell>
                <var-cell><var-input variant="outlined" placeholder="歌手" size="small" v-model="CardStatus.Meta.author"/></var-cell>
            </div>
        </div>
        <var-space class="absolute right-5 bottom-5">
            <var-button type="primary" :loading="CardStatus.ConfirmBtnLoadig" loading-type="wave" @click="saveVideoMeta">保存</var-button>
        </var-space>
    </var-popup>

</template>

<script setup>
import FramePage from '../modules/frame_page.vue'
import AdditionCard from '../modules/addition_card.vue'
import { reactive, computed, watch, ref, onMounted } from 'vue'
import { GetTaskListPage, GetTaskListAll, SetTaskDeleteState, UpdateTaskMeta, SaveVideoListTo, ExportVideoList } from '../../../wailsjs/go/wails_api/WailsApi'
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
    ExportVideoList().then(err => {
        if (err) {
            Snackbar.error("导出失败: " + err);
        } else if (err === null) {
             Snackbar.success("操作完成");
        }
    });
}

// 修改视频删除状态
function setDeleteState(index) {
    const currentState = TaskList.tasks[index].IsDelete;
    // Optimistic update
    TaskList.tasks[index].IsDelete = !currentState;
    TaskList.tasks[index].delete = !currentState; // Keep compatibility if 'delete' prop used in template

    SetTaskDeleteState(TaskList.tasks[index].Index, !currentState).then(() => { // Use real index from task
         // Success
    }).catch(err => {
        // Revert on failure
        TaskList.tasks[index].IsDelete = currentState;
        TaskList.tasks[index].delete = currentState;
        Snackbar.error("状态更新失败");
    });
}

// 打开右侧面板
function openRightPanel(index) {
    CardStatus.ListIndex = index;
    // Assuming TaskList.tasks matches backend order for now, or using Index from task
    const task = TaskList.tasks[index];
    CardStatus.Meta.song_name = task.SongName;
    CardStatus.Meta.author = task.SongAuthor;
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
        
        // Map IsDelete to 'delete' prop for compatibility with template v-ifs
        result.forEach(task => {
            task.delete = task.IsDelete;
        });

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
    const index = CardStatus.ListIndex;
    const taskIndex = TaskList.tasks[index].Index;
    
    UpdateTaskMeta(taskIndex, CardStatus.Meta.song_name, CardStatus.Meta.author).then(() => {
        TaskList.tasks[index].SongName = CardStatus.Meta.song_name;
        TaskList.tasks[index].SongAuthor = CardStatus.Meta.author;
        Snackbar.success("保存成功");
        CardStatus.RightPanel = false;
    }).catch(err => {
         Snackbar.error("保存失败: " + err);
    });
}
</script>
