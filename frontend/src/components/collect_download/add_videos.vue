<template>
    <FramePage title="添加内容">
        <var-button type="danger" round icon-container style="float: right" @click="() => {
            Dialog('删除当前任务列表').then(result => {
                if (result === 'confirm') {
                    ResetDownloadList().then(() => {
                        emit('updateBadge');
                        Snackbar.success('删除成功');
                    });
                }
            });
        }">
            <var-icon name="window-close" />
        </var-button>

        <var-radio-group v-model="QueryType" @change="queryInfornation">
            <var-radio :checked-value="0">收藏夹</var-radio>
            <var-radio :checked-value="1">视频合集</var-radio>
            <var-radio :checked-value="2">视频链接</var-radio>
            <var-radio :checked-value="3">AUID</var-radio>
            <var-radio :checked-value="4">用户空间</var-radio>
        </var-radio-group>

        <var-input :placeholder="inputTip" v-model="input" clearable
        style="margin-bottom: 25px;" >
            <template #prepend-icon>
                <var-icon name="magnify" />
            </template>
        </var-input>
        <!-- TODO: 一键粘贴剪切板 -->


    </FramePage>
    <var-collapse-transition :expand="CardStatus.InfoCard" style="margin-bottom: 30px;">
        <AdditionCard title="查询结果">
            <div style="display: flex; justify-content: center;">
                <img :src="resp.cover" style="width: 250px; height: 156px; border-radius: 8px;">
                <div style="margin-left: 20px; display: flex; flex-direction: column;">
                    <text style="font-size: 18px; font-weight: 700;">{{ resp.title }}</text>
                    <br>
                    <var-chip type="primary" style="margin-top: 5px;">视频数量：{{ resp.count }}</var-chip>
                    <!-- <br> -->
                    <var-chip type="primary" style="margin-top: 5px;">创建人：{{ resp.up_name }}</var-chip>
                </div>
            </div>

            <div></div>

            <var-space justify="flex-end">
                <var-button type="primary" @click="openRightPanel">添加至列表</var-button>
            </var-space>
        </AdditionCard>
    </var-collapse-transition>

    <var-popup position="right" v-model:show="CardStatus.RightPanel" :overlay-style="{backgroundColor: 'rgba(0, 0, 0, 0.2)'}" style=" height: 75%; right: 35px; top: auto; bottom: 35px; border-radius: 8px;">
        <div class="popup-example-block">
            <h3>添加选项</h3>
            <var-cell> 下载全部
                <template #extra>
                    <var-switch v-model="CardStatus.DownloadAll" @click="parms.options.downCount = 0" variant :disabled="!CardStatus.EnableDownloadAll" />
                </template>
            </var-cell>

            <var-cell> 下载数量
                <template #extra>
                    <var-counter :min="0" :max="resp.count" v-model="parms.options.downCount" :disabled="CardStatus.DownloadAll"/>
                </template>
            </var-cell>

            <var-cell> 下载分集
                <template #extra>
                    <var-switch v-model="parms.options.downPart" variant />
                </template>
            </var-cell>
            
            <h3>元数据选项</h3>
            <var-cell> 歌曲名称
                <template #extra>
                    <var-switch v-model="parms.options.songName" @click.stop variant />
                </template>
            </var-cell>

            <var-cell> 歌曲封面
                <template #extra>
                    <var-switch v-model="parms.options.songCover" @click.stop variant />
                </template>
            </var-cell>

            <var-cell> 歌曲作者
                <template #extra>
                    <var-switch v-model="parms.options.songAuthor" @click.stop variant />
                </template>
            </var-cell>
            <br>
            <var-space style="position: absolute; right: 20px; bottom: 20px;">
                <var-button type="primary" :loading="CardStatus.ConfirmBtnLoadig" loading-type="wave" @click="addItToList">确定添加</var-button>
            </var-space>
        </div>
    </var-popup>
</template>

<script setup>
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
    ResetDownloadList
} from '../../../wailsjs/go/wails_api/WailsApi'
import { EventsOn, EventsEmit } from '../../../wailsjs/runtime'
import {Dialog, Snackbar} from '@varlet/ui'

const props = defineProps(['parms', 'status'])
const emit = defineEmits(['update:parms', 'update:status', 'updateBadge'])

const inputTip = ref("请输入 收藏夹网页 URL")
const input = ref("")
const QueryType = ref(0)
const addItToList = ref(null)

// 本页卡片状态
const CardStatus = reactive({
    InfoCard: false,
    InfoVCount: false,
    InfoVUpper: false,
    RightPanel: false,
    DownloadAll: true,
    EnableDownloadAll: true,
    ConfirmBtnLoadig: false
})

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

// 查询函数返回值
const resp = reactive({
    bvid: "",
    fid: 0,
    mid: 0,
    title: "",
    cover: "",
    count: 0,
    up_name: "",
    up_avatar: "",
})

// 输入的 ID 变化时查询歌曲信息
watch(input, () => {
    queryInfornation();
})

// 一键添加事件
EventsOn('addToList', (url, type) => {
    QueryType.value = type === 0 ? 1 : 0;
    input.value = url;
    openRightPanel();
})

// 查询信息函数
function queryInfornation() {

    switch(QueryType.value) {
        case 0: // collect
            inputTip.value = '请输入 收藏夹网页 URL'
            break;
        case 1: // comliation
            inputTip.value = '请输入 视频合集网页 URL'
            break;
        case 2: // video
            inputTip.value = '请输入 网页端视频分享链接'
            break;
        case 3: // audio
            inputTip.value = '请输入 AUID'
            break;
        case 4: // profile
            inputTip.value = '请输入 用户空间 URL'
        break;
    }

    if (input.value === "") {
        // 空输入判断
        CardStatus.InfoCard = false;
        return;
    }
    switch(QueryType.value) {
        case 0: // Collection
            const fid = extractCollect(input.value);
            if (fid == null) {
                CardStatus.InfoCard = false;
                Snackbar.warning("链接匹配失败");
                return;
            }
            QueryCollection(fid).then(result => {
                resp.title = result.Data.Info.title;
                resp.cover = result.Data.Info.cover;
                resp.up_name = result.Data.Info.Upper.name;
                resp.up_avatar = result.Data.Info.Upper.face;
                resp.count = result.Data.Info.media_count;

                resp.fid = fid;
                CardStatus.InfoCard = true;
            });
            break;
        case 1: // Compilation
            const resultId = extractCompilation(input.value);
            if (resultId == null) {
                CardStatus.InfoCard = false;
                Snackbar.warning("链接匹配失败");
                return;
            }
            
            QueryCompilation(Number(resultId.mid), Number(resultId.fid)).then(result => {
                resp.title = result.Data.Meta.name;
                resp.cover = result.Data.Meta.cover;
                resp.up_name = "";
                resp.up_avatar = "";
                resp.count = result.Data.Meta.total;

                resp.fid = Number(resultId.fid);
                resp.mid = Number(resultId.mid);
                CardStatus.InfoCard = true;
            });
            break;
        case 2: // Video
            const bvid = extractBvid(input.value);
            if (bvid == null) {
                CardStatus.InfoCard = false;
                Snackbar.warning("链接匹配失败");
                return;
            }
            QueryVideo(bvid).then(result => {
                resp.title = result.Meta.title;
                resp.cover = result.Meta.cover;
                resp.up_name = result.Up.name;
                resp.up_avatar = result.Up.avatar;
                resp.count = 1;

                resp.bvid = bvid;
                CardStatus.InfoCard = true;
            });
            break;

        case 3:
            const auid = extractAuid(input.value);
            if (auid == null) {
                CardStatus.InfoCard = false;
                Snackbar.warning("链接匹配失败");
                return;
            }
            QueryAudio(auid).then(result => {
                resp.title = result.Meta.title;
                resp.cover = result.Meta.cover;
                resp.up_name = result.Up.author;
                resp.count = 1;

                resp.bvid = auid;
                CardStatus.InfoCard = true;
            });
            break;
        
        case 4:
            const mid = extractMid(input.value);
            if (mid == null) {
                CardStatus.InfoCard = false;
                Snackbar.warning("链接匹配失败");
                return;
            }
            QueryProfileVideo(mid).then(result => {
                resp.title = "用户个人空间";
                resp.cover = "";
                resp.up_name = "";
                resp.count = result;

                resp.bvid = mid;
                CardStatus.InfoCard = true;
            });
            break;
    }
}

function openRightPanel() {
    CardStatus.DownloadAll = true;
    props.parms.options.downPart = true;
    CardStatus.EnableDownloadAll = true;

    switch (QueryType.value) {
        case 0: // Collection
            addItToList.value = () => {
                CardStatus.ConfirmBtnLoadig = true;
                AddCollectionToList(parms.value.videoListPath, resp.fid, props.parms.options.downCount, props.parms.options.downPart).then(()=>{
                    afterAdd();
                });
            }
            break;
    
        case 1: // Compilation
            addItToList.value = () => {
                CardStatus.ConfirmBtnLoadig = true;
                AddCompilationToList(parms.value.videoListPath, Number(resp.mid), Number(resp.fid), props.parms.options.downCount, props.parms.options.downPart).then(()=>{
                    afterAdd();
                });
            }
            break;

        case 2: // Video
            CardStatus.EnableDownloadAll = false;
            addItToList.value = () => {
                CardStatus.ConfirmBtnLoadig = true;
                AddVideoToList(parms.value.videoListPath, resp.bvid, props.parms.options.downPart).then(()=>{
                    afterAdd();
                });
            }
            break;

        case 3: //Audio
            CardStatus.EnableDownloadAll = false;
            addItToList.value = () => {
                CardStatus.ConfirmBtnLoadig = true;
                AddAudioToList(parms.value.videoListPath, resp.bvid).then(()=>{
                    afterAdd();
                });
            }
            break;

        case 4: //Profile
            addItToList.value = () => {
                CardStatus.ConfirmBtnLoadig = true;
                AddProfileVideoToList(parms.value.videoListPath, Number(resp.bvid), props.parms.options.downCount, props.parms.options.downPart).then(()=>{
                    afterAdd();
                });
            }
            break;
    }
    CardStatus.RightPanel = true;
}

function afterAdd() {
    Snackbar.success("添加完成");
    CardStatus.RightPanel = false;
    input.value = "";
    props.parms.options.downCount = 0;
    emit('updateBadge');
    CardStatus.ConfirmBtnLoadig = false;
}

// 过滤视频分享链接
function extractBvid(url) {
    const regex = /BV\w+/;
    const match = url.match(regex);

    if (match) {
        return match[0];
    } else {
        return null;
    }
}

// 过滤收藏夹链接
function extractCollect(url) {
    const regex = /fid=(\d+)/;
    const match = url.match(regex);

    if (match) {
        return match[1];
    } else {
        return null;
    }
}

// 过滤视频合集链接
function extractCompilation(url) {
    const regex = /space\.bilibili\.com\/(\d+)\/favlist\?fid=(\d+)&/;
    const match = url.match(regex);

    if (match) {
        return {
            mid: match[1],
            fid: match[2]
        };
    } else {
        return null;
    }
}

// 过滤AUID
function extractAuid(auid) {
    const match = auid.match(/^au(\d+)$/);
    return match ? match[1] : null;
}

// 过滤用户空间 MID
function extractMid(url) {
    const regex = /space\.bilibili\.com\/(\d+)/;
    const match = url.match(regex);
    if (match) {
        return match[1]; // 返回匹配到的数字部分
    }
    return null; // 如果没有匹配到，返回 null
}

</script>

<style>
.var-popup--right {
    height: 75%;
    position: absolute;
    right: 35px;
    bottom: 35px;
    border-radius: 8px;
}
</style>