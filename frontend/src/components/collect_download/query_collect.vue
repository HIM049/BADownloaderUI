<template>
    <FramePage title="收藏夹下载">
        <var-input placeholder="收藏夹 ID / 收藏夹 URL" v-model="inputFavID" clearable 
        style="margin-bottom: 25px;" >
            <template #prepend-icon>
                <var-icon name="magnify" />
            </template>
        </var-input>

        <var-card 
            :title="resp.title" :src="resp.cover" layout="row" ripple outline 
            v-show="InfCard">
            <template #description>
                <var-divider />
                <div style="margin-left: 20px;">
                    <label>内容数量：{{ resp.count }}</label>
                    <br>
                    <label>创建者：{{ resp.up_name }}</label>
                </div>
            </template>
        </var-card>
    </FramePage>
</template>

<script setup>
import FramePage from '../modules/frame_page.vue'
import { reactive, computed, ref, watch } from 'vue'
import { SearchFavListInformation, SearchCompListInformation } from '../../../wailsjs/go/main/App'
import { Snackbar } from '@varlet/ui'

const props = defineProps(['parms', 'status'])
const emit = defineEmits(['update:parms', 'update:status'])

const inputFavID = ref("")
const InfCard = ref(false)

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
    title: "搜索结果",
    cover: "",
    count: 0,
    up_name: "",
    up_avatar: "",
})

// 在输入字段中提取收藏夹 ID
// function extractURL(url) {
//     try {
//         var parsedUrl = new URL(url);
//     } catch (error) {
//         // 不是 URL ，直接返回
//         return url
//     }
//     // 提取特定参数
//     var searchParams = new URLSearchParams(parsedUrl.search);
//     var fid = searchParams.get("fid");
//     return fid
// }

function extractURL(url) {
    try {   
        var urlParams = new URLSearchParams(new URL(url).search);
    }catch{
        return null
    }
  var fid = Number(urlParams.get("fid"));
  var ctype = urlParams.get("ctype");

  var isCompliation = ctype == "21";

  var regex = /\/(\d+)\//;
  var matches = url.match(regex);
  var mid = Number(matches[1]);
  return {
    mid,
    fid,
    isCompliation,
  };
}

// 输入的 ID 变化时查询歌曲信息
watch(inputFavID, (newid) => {
    var urlParams = extractURL(newid)
    if (urlParams == null) {
        return
    }
    if (!urlParams.isCompliation) {
        props.parms.isComp = false
        props.parms.favListID = String(urlParams.fid)
        SearchFavListInformation(props.parms.favListID).then(result => {
            // 判断信息有效性
            if (result.message == "0") {
                resp.title = result.Data.Info.title;
                resp.cover = result.Data.Info.cover;
                resp.count = result.Data.Info.media_count;
                resp.up_name = result.Data.Info.Upper.name;
                resp.up_avatar = result.Data.Info.Upper.face;
                props.parms.count = result.Data.Info.media_count;

                // 开放创建列表按钮
                InfCard.value = true;
                props.status.allowNext = true;
            } else {
                // 无效的收藏夹
                // Snackbar.warning("无效的收藏夹");
                
                // 关闭创建列表按钮
                props.status.allowNext = false;
                InfCard.value = false;
            }
        })
    } else {
        props.parms.isComp = true
        props.parms.favListID = String(urlParams.fid)
        props.parms.mid = String(urlParams.mid)

        SearchCompListInformation(urlParams.mid, urlParams.fid).then(result => {
            if (result.message == "0") {
                resp.title = result.Data.Meta.name;
                resp.cover = result.Data.Meta.cover;
                resp.count = result.Data.Meta.total;
                resp.up_name = null;
                resp.up_avatar = null;
                props.parms.count = result.Data.Meta.total;

                // 开放创建列表按钮
                InfCard.value = true;
                props.status.allowNext = true;
            } else {
                // 无效的收藏夹
                Snackbar.warning("无效的合集");
                // 关闭创建列表按钮
                props.status.allowNext = false;
                InfCard.value = false;
            }
        })
    }
})
</script>