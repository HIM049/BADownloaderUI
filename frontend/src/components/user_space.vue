<template>
    <FramePage title="个人空间">
        <var-space style="align-items: center;" justify="center" v-if="!showLoginWindow && !is_login">
            你还没有登录
            <var-button type="primary" @click="login">登录</var-button>
        </var-space>
        <var-collapse-transition :expand="showLoginWindow">
            <var-paper style="background-color: var(--color-primary-container);">
                <h3>扫码登录</h3>
                <img style="display: flex; margin: 0 auto; border-radius: 15px;" :src="qrcodeStr" alt="Image">
                <p style="display: flex; justify-content: center;">{{ loginText }}</p>
            </var-paper>
        </var-collapse-transition>   
        <var-space v-if="is_login && CardStatus.LoadUserInf" justify="center">
            <var-avatar :src="user_Information.avatar" />
            <p>{{ user_Information.name }}</p>
        </var-space>     
    </FramePage>
    <AdditionCard title="创建的收藏夹" v-if="is_login && CardStatus.LoadUsersCollect">
        <var-paper style="background-color: var(--color-primary-container);">
            <div v-for="(collect, index) in user_collect.List" style="margin: 5px 20px;">
                <var-space justify="space-between" align="center">
                    <text style="font-size: 15px; font-weight: 600;">{{ collect.title }}</text> 
                    <var-button type="primary" @click="addToList('https://space.bilibili.com/'+user_collect.user_mid+'/favlist?fid='+collect.id+'&ftype=create', 11)"><var-icon name="plus" />添加至列表</var-button>
                </var-space>
            </div>
        </var-paper>
    </AdditionCard>

    <AdditionCard title="收藏和订阅" v-if="is_login && CardStatus.LoadFavCollect">
        
        <var-paper style="background-color: var(--color-primary-container);">
            <div v-for="(collect, index) in user_Favourite.List" style="margin: 5px 20px;">
                <var-space justify="space-between" align="center">
                    <text style="font-size: 15px; font-weight: 600;">{{ collect.title }}</text> 
                    
                    <var-button type="primary" @click="
                        addToList(collect.attr == 0 ? 'https://space.bilibili.com/'+user_collect.user_mid+'/favlist?fid='+collect.id+'&ftype=collect&ctype=21':'https://space.bilibili.com/'+user_collect.user_mid+'/favlist?fid='+collect.id+'&ftype=collect&ctype=11', collect.attr)
                    "><var-icon name="plus" />添加至列表</var-button>
                </var-space>
            </div>
        </var-paper>

        <var-space style="display: flex; align-items: center;">
            <var-button-group type="primary" size="normal" outline >
                <var-button @click="page_index--">上一页</var-button>
                <var-button @click="page_index++">下一页</var-button>
            </var-button-group>
        </var-space>
    </AdditionCard>
</template>

<script setup>
import FramePage from '../components/modules/frame_page.vue'
import AdditionCard from './modules/addition_card.vue'
import { ref, reactive, onMounted, watch } from 'vue'
import { LoginBilibili } from '../../wailsjs/go/main/App'
import { LoadConfig, GetUsersCollect, GetFavCollect, GetUserInf } from '../../wailsjs/go/wails_api/WailsApi'
import { EventsOn, EventsEmit, ClipboardSetText } from '../../wailsjs/runtime'
import { Snackbar } from '@varlet/ui'

const loginText = ref("请扫描二维码登录") // 登录时的提示字符
const qrcodeStr = ref(null) // 二维码图片
const showLoginWindow = ref(false) // 展示登录窗口
const is_login = ref(false) // 登录状态
const page_index = ref(1) // 订阅的收藏夹的页码

const user_collect = ref([]); // 用户创建的收藏夹列表
const user_Favourite = ref([]); // 用户订阅的收藏夹列表
const user_Information = ref(null); // 用户信息
const fav_count = ref(0);

const CardStatus = reactive({
    LoadUsersCollect: false,
    LoadFavCollect: false,
    LoadUserInf: false,
})

onMounted(() => {
    checkLogin();
    if (is_login) {
        GetUsersCollect().then(result => {
            user_collect.value = result;
            CardStatus.LoadUsersCollect = true;
        })
        getFavCollect();
        GetUserInf().then(result => {
            user_Information.value = result;
            CardStatus.LoadUserInf = true;
        })
    }
});

// 一键添加至列表
function addToList(url, type) {
    EventsEmit('addToList', url, type);
}

// 获取订阅收藏夹列表
function getFavCollect() {
    GetFavCollect(page_index.value).then(result => {
        user_Favourite.value = result;
        fav_count.value = result.count;
        CardStatus.LoadFavCollect = true;
    })
}

// 订阅的收藏夹翻页
watch(page_index, (newValue) => {
    if (newValue <= 0) {
        page_index.value = 1;
    }else if (newValue > fav_count.value / 20 + 1) {
        page_index.value = parseInt(fav_count.value / 20 + 1);
    }
    getFavCollect();

})

// 获取二维码事件
EventsOn("qrcodeStr", (qr)=>{
    qrcodeStr.value = "data:image/png;base64," + qr;
})

// 登录状态事件
EventsOn("loginStatus", (status) => {
    if (status == "登录成功") {
        Snackbar.success("登录成功");
        checkLogin();
        setTimeout(() => {
            window.location.reload();
        }, 1000);
    }
    if (status == "二维码已失效") {
        Snackbar.warning("二维码已失效");
        showLoginWindow.value = false;
    }
    loginText.value = status;
})

// 检查登录状态
function checkLogin() {
    LoadConfig().then(result => {
        is_login.value = result.Account.is_login
    })
}

// 登录账户
function login() {
    LoginBilibili().then(result => {
        showLoginWindow.value = false;
    })
    showLoginWindow.value = true;
}
</script>