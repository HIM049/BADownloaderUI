<template>
    <FramePage title="添加内容">
        <template #actions>
            <var-space>
                <var-button type="primary" @click="importList">
                    <var-icon name="upload" />
                    导入列表
                </var-button>
                <var-button type="danger" icon-container @click="resetList">
                    <var-icon name="delete" />
                    清除列表
                </var-button>
            </var-space>
        </template>

        <var-radio-group v-model="sourceKey">
            <var-radio v-for="item in sources" :key="item.key" :checked-value="item.key">{{ item.label }}</var-radio>
        </var-radio-group>

        <var-input :placeholder="source.placeholder" v-model="input" clearable class="mb-6">
            <template #prepend-icon>
                <var-icon name="magnify" />
            </template>
        </var-input>
        <!-- TODO: 一键粘贴剪切板 -->
    </FramePage>

    <var-collapse-transition :expand="showResult" class="mb-8">
        <AdditionCard title="查询结果">
            <div class="flex justify-center">
                <img :src="result.cover" class="w-[250px] h-[156px] rounded-lg">
                <div class="ml-5 flex flex-col">
                    <text class="text-lg font-bold">{{ result.title }}</text>
                    <br>
                    <var-chip type="primary" class="mt-1.5">视频数量：{{ result.count }}</var-chip>
                    <var-chip type="primary" class="mt-1.5">创建人：{{ result.up_name }}</var-chip>
                </div>
            </div>

            <var-space justify="flex-end">
                <var-button type="primary" @click="openPanel">添加至列表</var-button>
            </var-space>
        </AdditionCard>
    </var-collapse-transition>

    <var-popup position="right" v-model:show="showPanel" :overlay-style="{backgroundColor: 'rgba(0, 0, 0, 0.2)'}" class="!h-[75%] !absolute !right-[35px] !bottom-[35px] !rounded-lg !top-auto">
        <div class="p-6">
            <h3>添加选项</h3>
            <var-cell> 下载全部
                <template #extra>
                    <var-switch v-model="options.downloadAll" variant :disabled="!source.ranged" />
                </template>
            </var-cell>

            <var-cell> 起始位置
                <template #extra>
                    <var-counter :min="1" :max="result.count" v-model="options.offset" />
                </template>
            </var-cell>

            <var-cell> 下载数量
                <template #extra>
                    <var-counter :min="0" :max="result.count" v-model="options.count" :disabled="options.downloadAll" />
                </template>
            </var-cell>

            <var-cell> 下载分集
                <template #extra>
                    <var-switch v-model="options.part" variant />
                </template>
            </var-cell>

            <h3>元数据选项</h3>
            <var-cell> 歌曲名称
                <template #extra>
                    <var-switch v-model="options.meta.song_name" @click.stop variant />
                </template>
            </var-cell>

            <var-cell> 歌曲封面
                <template #extra>
                    <var-switch v-model="options.meta.song_cover" @click.stop variant />
                </template>
            </var-cell>

            <var-cell> 歌曲作者
                <template #extra>
                    <var-switch v-model="options.meta.song_author" @click.stop variant />
                </template>
            </var-cell>
            <br>
            <var-space class="absolute right-5 bottom-5">
                <var-button type="primary" :loading="adding" loading-type="wave" @click="addToList">确定添加</var-button>
            </var-space>
        </div>
    </var-popup>
</template>

<script setup lang="ts">
import FramePage from '../modules/frame_page.vue'
import AdditionCard from '../modules/addition_card.vue'
import { reactive, computed, ref, watch } from 'vue'
import {
    QueryVideo,
    QueryCollection,
    QueryCompilation,
    QueryAudio,
    QueryProfileVideo,
    AddVideoToList,
    AddCollectionToList,
    AddCompilationToList,
    AddAudioToList,
    AddProfileVideoToList,
    ResetDownloadList,
    ImportVideoList
} from '../../../wailsjs/go/wails_api/WailsApi'
import { task, wails_api } from '../../../wailsjs/go/models'
import { onEvent, type SourceKey } from '../../events'
import { Dialog, Snackbar } from '@varlet/ui'
import { call } from '../../utils/call'
import { refreshListCount } from '../../state'

// 添加时传给后端的选项
interface AddOptions {
    offset: number
    count: number
    part: boolean
    meta: task.Option
}

// 可添加的来源：parse 从输入解析 ID，query 查询信息，add 添加到列表；ranged 表示支持起始位置和数量
interface Source<ID> {
    key: SourceKey
    label: string
    placeholder: string
    ranged: boolean
    parse(text: string): ID | null
    query(id: ID): Promise<wails_api.QueryResult>
    add(id: ID, options: AddOptions): Promise<void>
}

// defineSource 推断每个来源的 ID 类型，再统一放进来源表
function defineSource<ID>(source: Source<ID>): Source<unknown> {
    return source as unknown as Source<unknown>
}

// matchID 用正则匹配输入，成功时用 build 生成 ID 对象，失败返回 null
function matchID<ID>(text: string, regex: RegExp, build: (match: RegExpMatchArray) => ID): ID | null {
    const match = text.match(regex)
    return match ? build(match) : null
}

const sources = [
    defineSource({
        key: 'favlist',
        label: '收藏夹',
        placeholder: '请输入 收藏夹网页 URL',
        ranged: true,
        parse: (text) => matchID(text, /fid=(\d+)/, (m) => ({ fid: m[1] })),
        query: (id) => QueryCollection(id.fid),
        add: (id, o) => AddCollectionToList(id.fid, o.offset, o.count, o.part, o.meta),
    }),
    defineSource({
        key: 'season',
        label: '视频合集',
        placeholder: '请输入 视频合集网页 URL',
        ranged: true,
        parse: (text) => matchID(text, /space\.bilibili\.com\/(\d+)\/favlist\?fid=(\d+)&/, (m) => ({ mid: Number(m[1]), sid: Number(m[2]) })),
        query: (id) => QueryCompilation(id.mid, id.sid),
        add: (id, o) => AddCompilationToList(id.mid, id.sid, o.offset, o.count, o.part, o.meta),
    }),
    defineSource({
        key: 'video',
        label: '视频链接',
        placeholder: '请输入 网页端视频分享链接',
        ranged: false,
        parse: (text) => matchID(text, /BV\w+/, (m) => ({ bvid: m[0] })),
        query: (id) => QueryVideo(id.bvid),
        add: (id, o) => AddVideoToList(id.bvid, o.part, o.meta),
    }),
    defineSource({
        key: 'audio',
        label: 'AUID',
        placeholder: '请输入 AUID',
        ranged: false,
        parse: (text) => matchID(text, /^au(\d+)$/, (m) => ({ auid: m[1] })),
        query: (id) => QueryAudio(id.auid),
        add: (id, o) => AddAudioToList(id.auid, o.meta),
    }),
    defineSource({
        key: 'space',
        label: '用户空间',
        placeholder: '请输入 用户空间 URL',
        ranged: true,
        parse: (text) => matchID(text, /space\.bilibili\.com\/(\d+)/, (m) => ({ mid: m[1] })),
        query: (id) => QueryProfileVideo(id.mid),
        add: (id, o) => AddProfileVideoToList(Number(id.mid), o.offset, o.count, o.part, o.meta),
    }),
]

const sourceKey = ref<SourceKey>('favlist')
const source = computed(() => sources.find((item) => item.key === sourceKey.value) ?? sources[0])
const input = ref('')

// 最近一次查询成功的 ID 和结果
const target = ref<{ source: Source<unknown>, id: unknown } | null>(null)
const result = reactive<wails_api.QueryResult>({ title: '', cover: '', count: 0, up_name: '' })
const showResult = ref(false)

// 添加面板的选项
const showPanel = ref(false)
const adding = ref(false)
const options = reactive({
    downloadAll: true,
    offset: 1,
    count: 0,
    part: true,
    meta: { song_name: true, song_cover: true, song_author: true },
})

// 来源或输入变化时重新查询
watch([sourceKey, input], query)

// 个人空间页面的一键添加
onEvent('addToList', (kind, url) => {
    sourceKey.value = kind
    input.value = url
    openPanel()
})

// 解析输入并查询信息，只采用最后一次查询的结果
let querySeq = 0
async function query() {
    const seq = ++querySeq
    showResult.value = false
    target.value = null
    if (input.value === '') {
        return
    }

    const current = source.value
    const id = current.parse(input.value)
    if (id == null) {
        Snackbar.warning('链接匹配失败')
        return
    }

    const { ok, data } = await call(current.query(id), '查询失败')
    if (!ok || seq !== querySeq) {
        return
    }
    Object.assign(result, data)
    target.value = { source: current, id }
    showResult.value = true
}

function openPanel() {
    options.downloadAll = true
    options.part = true
    showPanel.value = true
}

// 按面板选项添加到列表
async function addToList() {
    if (target.value == null) {
        Snackbar.warning('请等待查询完成')
        return
    }

    adding.value = true
    const { ok } = await call(target.value.source.add(target.value.id, {
        offset: options.offset || 1,
        count: options.downloadAll ? 0 : options.count,
        part: options.part,
        meta: { ...options.meta },
    }), '添加失败')
    adding.value = false
    if (!ok) {
        return
    }

    Snackbar.success('添加完成')
    showPanel.value = false
    input.value = ''
    options.count = 0
    refreshListCount()
}

// 导入列表
async function importList() {
    const { ok } = await call(ImportVideoList(), '导入失败')
    if (ok) {
        Snackbar.success('导入成功')
        refreshListCount()
    }
}

// 确认后清除任务列表
async function resetList() {
    if (await Dialog('清除当前任务列表') !== 'confirm') {
        return
    }
    const { ok } = await call(ResetDownloadList(), '清除失败')
    if (ok) {
        Snackbar.success('清除成功')
        refreshListCount()
    }
}
</script>
