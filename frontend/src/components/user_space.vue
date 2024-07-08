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
        <var-space v-if="is_login" justify="center">
            <var-avatar :src="user_Information.avatar" />
            <p>{{ user_Information.name }}</p>
        </var-space>     
    </FramePage>
    <AdditionCard title="创建的收藏夹" v-if="is_login">
            <var-cell v-for="(collect, index) in user_collect.List" style="margin: 0 10px;">
                {{ collect.title }}
                <template #extra>
                    <var-button type="primary" @click="ClipboardSetText('https://space.bilibili.com/'+user_collect.user_mid+'/favlist?fid='+collect.id+'&ftype=create').then(Snackbar.success('复制成功'))">复制链接</var-button>
                </template>
            </var-cell>
        </AdditionCard>

        <AdditionCard title="收藏和订阅" v-if="is_login">
            <var-cell v-for="(collect, index) in user_Favourite.List" style="margin: 0 10px;">
                {{ collect.title }}
                <template #extra>
                    <var-button type="primary" @click="
                        ClipboardSetText(collect.attr == 0 ? 'https://space.bilibili.com/'+user_collect.user_mid+'/favlist?fid='+collect.id+'&ftype=collect&ctype=21':'https://space.bilibili.com/'+user_collect.user_mid+'/favlist?fid='+collect.id+'&ftype=collect&ctype=11').then(Snackbar.success('复制成功'))
                    ">复制链接</var-button>
                </template>
            </var-cell>
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
import { LoginBilibili, LoadConfig, GetUsersCollect, GetFavCollect, GetUserInf } from '../../wailsjs/go/main/App'
import { EventsOn, ClipboardSetText } from '../../wailsjs/runtime'
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

onMounted(() => {
    checkLogin();
    if (is_login) {
        GetUsersCollect().then(result => {
            user_collect.value = result;
        })
        getFavCollect();
        GetUserInf().then(result => {
            user_Information.value = result;
        })
    }
});

// 获取订阅收藏夹列表
function getFavCollect() {
    GetFavCollect(page_index.value).then(result => {
        user_Favourite.value = result;
        fav_count.value = result.count;
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