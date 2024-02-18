<template>
        
    <StepBar :pageNum="step" />
    <el-main>
        <!-- <div id="icon-banner">
            <img src="./image/icon-non-bg.png" style="width: 350px;" v-show="!status.showStep">
        </div> -->
        <!-- 收藏夹信息输入 -->
        <div id="fav-input" class="fav-input">
            <el-input v-model="inputFavID" size="large" placeholder="请输入 收藏夹 ID / 收藏夹 URL" class="input-with-select"
                @input="queryFavListInformation" clearable>

                <template #prepend>
                    <el-button @click="queryFavListInformation"><el-icon>
                            <Search />
                        </el-icon></el-button>
                </template>
            </el-input>
        </div>
        <transition name="el-fade-in-linear">
            <favInf :resp="resp" v-show="status.showStep" />
        </transition>
    </el-main>
    <FootBar :status="status" text="继续" @back="$emit('back')" @next="$emit('next')" />
</template>

<script setup>
import StepBar from '../components/modules/step_bar.vue'
import FootBar from '../components/modules/footer.vue'
import favInf from './modules/fav_information.vue'
import { reactive, computed, ref } from 'vue'
import { SearchFavListInformation } from '../../wailsjs/go/main/App'

const step = 0
const props = defineProps(['parms'])
const emit = defineEmits(['update:parms', 'back', 'next'])

// 底栏状态
const status = reactive({
    showBack: false,
    showNext: true,
    allowBack: true,
    allowNext: false,
    // showStep: false,
})

const inputFavID = ref("")

const parms = computed({
    get() {
        return props.parms
    },
    set(parms) {
        emit('update:parms', parms)
    }
})

// 查询函数返回值
const resp = reactive({
    title: "",
    cover: "",
    count: 0,
    up_name: "",
    up_avatar: "",
})

function extractURL(url) {
    // 尝试创建 URL 对象
    try {
        var parsedUrl = new URL(url);
    } catch (error) {
        // 不是 URL ，直接返回
        console.log(error);
        return url
    }
    // 获取参数部分
    var searchParams = new URLSearchParams(parsedUrl.search);
    // 提取特定参数
    var fid = searchParams.get("fid");
    return fid
}

// 查询收藏夹信息
function queryFavListInformation() {
    // TODO: 输入校验 
    // parms.favListID = parms.favListID.replace(/\D/g, '');
    // if (parms.favListID.length > 18) {
    //     parms.favListID = parms.favListID.slice(0, 100);
    // }
    props.parms.favListID = extractURL(inputFavID.value)
    SearchFavListInformation(props.parms.favListID).then(result => {
        // 判断信息有效性
        // console.log(result);
        if (result.message == "0") {
            resp.title = result.Data.Info.title;
            resp.cover = result.Data.Info.cover;
            resp.count = result.Data.Info.media_count;
            resp.up_name = result.Data.Info.Upper.name;
            resp.up_avatar = result.Data.Info.Upper.face;
            props.parms.count = result.Data.Info.media_count;


            // 开放创建列表按钮
            status.allowNext = true;
        } else {
            // 无效的收藏夹
            resp.title = "无效的收藏夹";
            ElMessage.warning("无效的收藏夹");
            // 关闭创建列表按钮
            status.allowNext = false;
        }
        // 展示收藏夹信息        
        status.showStep = true;
    })
}
</script>

<style>
#icon-banner {
    width: 350px;
    margin: 0 auto;
}

div#fav-input {
    width: 500px;
    margin: 0 auto;
}
</style>