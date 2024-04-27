<template>
    <FramePage title="单曲下载">
        <var-input placeholder="纯数字 AUID" v-model="auid" clearable 
        style="margin-bottom: 25px;" >
            <template #prepend-icon>
                <var-icon name="magnify" />
            </template>
        </var-input>

        <var-card 
            title="搜索结果" :src="songInf.cover" layout="row" ripple outline
            v-show="InfCard"
            style="margin-bottom: 20px" >
            <template #description>
                <var-divider />
                <div>
                    <var-cell><var-input variant="outlined" placeholder="曲名" size="small" v-model="parms.song_name" /></var-cell>
                    <var-cell><var-input variant="outlined" placeholder="歌手" size="small" v-model="parms.author" /></var-cell>
                </div>
            </template>            
        </var-card>

        <var-card 
            title="下载选项" ripple outline
            style="margin-bottom: 20px"  >
            <template #description>
                <var-divider />
                <div style="margin-left: 20px;">
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
                </div>
            </template>
        </var-card>

        <var-space justify="flex-end">
            <var-button type="primary" @click="audioDownload">开始下载</var-button>
        </var-space>

    </FramePage>
</template>

<script setup>
import FramePage from '../components/modules/frame_page.vue'
import { reactive, ref, watch } from 'vue'
import { AudioDownload, SearchSongInformation } from '../../wailsjs/go/main/App'
import { Snackbar } from '@varlet/ui'


// 展示状态
const InfCard = ref(false)

const auid = ref('')

// 歌曲信息
const songInf = reactive({
    title: "",
    cover: "",
})

// 下载参数
const parms = reactive({
    song_name: "",
    author: "",
    // 下载设置
    options: reactive({
        songName: true,
        songCover: true,
        songAuthor: true,
    })
})

// 输入的 ID 变化时查询歌曲信息
watch(auid, (newid) => {
    SearchSongInformation(newid).then(result => {
        if (result.msg == "success") {
            songInf.title = result.Data.title
            songInf.cover = result.Data.cover
            parms.song_name = result.Data.title
            parms.author = result.Data.author

            InfCard.value = true
        } else {
            Snackbar.warning("无效的 AUID");
            InfCard.value = false
        }
    })
})

// 下载歌曲
function audioDownload() {
    if (InfCard.value == false) {
        Snackbar.warning("无效的 AUID");
        return
    }
    Snackbar.success("开始下载")
    var opt = {
        song_name: parms.options.songName,
        song_cover: parms.options.songCover,
        song_author: parms.options.songAuthor,
    }
    AudioDownload(opt, auid, parms.song_name, parms.author, songInf.title).then(result => {
        Snackbar.success("下载完成")
    })
}
</script>