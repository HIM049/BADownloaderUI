<template>
    <FramePage title="列表编辑">
        <template #actions>
            <var-button type="primary" @click="saveListTo">另存为</var-button>
        </template>
    </FramePage>
    <AdditionCard>
        <var-list :finished="list.finished" :loading="list.loading" @load="load" :immediate-check="false">
            <li v-for="(task, index) in list.tasks" :key="task.Index" class="list-none">
                <var-card :title="task.SongName" :src="task.CoverUrl" layout="row" image-width="250px" outlines v-if="!task.IsDelete" class="mb-5">
                    <template #description>
                        <var-divider />
                        <div class="mx-2.5">
                            <var-chip plain type="info" class="mb-1.5">歌曲名称：{{ task.SongName }}</var-chip>
                            <var-chip plain type="info" class="mb-[30px]">歌曲作者：{{ task.SongAuthor }}</var-chip>
                        </div>
                    </template>

                    <template #extra>
                        <div class="flex items-center">
                            <var-button type="danger" size="large" round @click="toggleDelete(index)" class="mr-2.5"> <var-icon name="delete" /> </var-button>
                            <var-button type="primary" @click="openPanel(index)">编辑</var-button>
                        </div>
                    </template>
                </var-card>

                <var-card :title="task.SongName" outlines class="mb-5 h-[187px]" v-else>
                    <template #description>
                        <div class="flex justify-center">
                            <h3>已设为删除</h3>
                        </div>
                    </template>
                    <template #extra>
                        <div class="flex items-center">
                            <var-button type="success" @click="toggleDelete(index)" class="mr-2.5"> <var-icon name="history" />  恢复 </var-button>
                            <var-button type="primary" disabled>编辑</var-button>
                        </div>
                    </template>
                </var-card>
            </li>
            <template #finished>
                <div class="py-5 text-center text-[#888]">已经到底了</div>
            </template>
        </var-list>
    </AdditionCard>

    <var-popup position="right" v-model:show="panel.show" :overlay-style="{backgroundColor: 'rgba(0, 0, 0, 0.2)'}" class="!h-[75%] !absolute !right-[35px] !bottom-[35px] !rounded-lg !top-auto">
        <div class="p-6 w-[280px]">
            <h3>元数据编辑</h3>
            <div>
                <var-cell><var-input variant="outlined" placeholder="曲名" size="small" v-model="panel.songName"/></var-cell>
                <var-cell><var-input variant="outlined" placeholder="歌手" size="small" v-model="panel.author"/></var-cell>
            </div>
        </div>
        <var-space class="absolute right-5 bottom-5">
            <var-button type="primary" :loading="panel.saving" loading-type="wave" @click="saveMeta">保存</var-button>
        </var-space>
    </var-popup>
</template>

<script setup lang="ts">
import FramePage from '../modules/frame_page.vue'
import AdditionCard from '../modules/addition_card.vue'
import { reactive, watch, onMounted } from 'vue'
import { GetTaskListPage, SetTaskDeleteState, UpdateTaskMeta, ExportVideoList } from '../../../wailsjs/go/wails_api/WailsApi'
import { task } from '../../../wailsjs/go/models'
import { onEvent } from '../../events'
import { Snackbar } from '@varlet/ui'
import { call } from '../../utils/call'
import { flow } from '../../state'

// 与后端 GetTaskListPage 的每页数量一致
const pageSize = 20

// 已加载的任务和分页状态
const list = reactive({
    tasks: [] as task.TaskInfo[],
    page: 0,
    loading: false,
    finished: false,
})

// 元数据编辑面板
const panel = reactive({
    show: false,
    index: 0,
    songName: '',
    author: '',
    saving: false,
})

onMounted(load)

// 进入列表编辑页或后端导入列表后刷新
watch(() => flow.pageIndex, (index) => {
    if (index === 1) reload()
})
onEvent('refreshVideoList', reload)

// 重新从第一页加载列表
function reload() {
    list.tasks = []
    list.page = 0
    list.finished = false
    load()
}

// 加载下一页，不足一页时结束
async function load() {
    list.loading = true
    const { ok, data } = await call(GetTaskListPage(list.page), '加载失败')
    list.loading = false
    if (!ok) {
        return
    }
    list.tasks.push(...data)
    list.page++
    list.finished = data.length < pageSize
}

// 另存列表
async function saveListTo() {
    const { ok } = await call(ExportVideoList(), '导出失败')
    if (ok) {
        Snackbar.success('操作完成')
    }
}

// 切换删除状态，失败时恢复
async function toggleDelete(index: number) {
    const task = list.tasks[index]
    task.IsDelete = !task.IsDelete
    const { ok } = await call(SetTaskDeleteState(task.Index, task.IsDelete), '状态更新失败')
    if (!ok) {
        task.IsDelete = !task.IsDelete
    }
}

function openPanel(index: number) {
    const task = list.tasks[index]
    panel.index = index
    panel.songName = task.SongName
    panel.author = task.SongAuthor
    panel.show = true
}

// 保存元数据修改
async function saveMeta() {
    const task = list.tasks[panel.index]
    panel.saving = true
    const { ok } = await call(UpdateTaskMeta(task.Index, panel.songName, panel.author), '保存失败')
    panel.saving = false
    if (!ok) {
        return
    }
    task.SongName = panel.songName
    task.SongAuthor = panel.author
    Snackbar.success('保存成功')
    panel.show = false
}
</script>
