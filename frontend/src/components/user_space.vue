<template>
    <FramePage title="个人空间">
        <var-space class="items-center" justify="center" v-if="!showLoginWindow && !is_login">
            你还没有登录
            <var-button type="primary" @click="login">登录</var-button>
        </var-space>
        <var-collapse-transition :expand="showLoginWindow">
            <var-paper style="background-color: var(--color-primary-container);">
                <h3>扫码登录</h3>
                <img class="flex mx-auto rounded-[15px]" :src="qrcodeStr" alt="Image">
                <p class="flex justify-center">{{ loginText }}</p>
            </var-paper>
        </var-collapse-transition>
        <var-space v-if="is_login && CardStatus.LoadUserInf" justify="center">
            <var-avatar :src="user_Information.avatar" />
            <p class="text-md font-bold my-3">{{ user_Information.name }}</p>
        </var-space>
    </FramePage>
    <AdditionCard title="创建的收藏夹" v-if="is_login && CardStatus.LoadUsersCollect">
        <var-paper style="background-color: var(--color-primary-container);">
            <div v-for="(collect, index) in user_collect.List" class="my-1.5 mx-5">
                <var-space justify="space-between" align="center">
                    <text class="text-[15px] font-semibold">{{ collect.title }}</text>
                    <var-button type="primary" @click="addCreated(collect)"><var-icon name="plus" />添加至列表</var-button>
                </var-space>
            </div>
        </var-paper>
    </AdditionCard>

    <AdditionCard title="收藏和订阅" v-if="is_login && CardStatus.LoadFavCollect">

        <var-paper style="background-color: var(--color-primary-container);">
            <div v-for="(collect, index) in user_Favourite.List" class="my-1.5 mx-5">
                <var-space justify="space-between" align="center">
                    <text class="text-[15px] font-semibold">{{ collect.title }}</text>

                    <var-button type="primary" @click="addCollected(collect)"><var-icon name="plus" />添加至列表</var-button>
                </var-space>
            </div>
        </var-paper>

        <var-space class="flex items-center">
            <var-button-group type="primary" size="normal" outline>
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
import { LoadConfig, GetUsersCollect, GetFavCollect, GetUserInf, LoginBilibili } from '../../wailsjs/go/wails_api/WailsApi'
import { EventsOn, EventsEmit } from '../../wailsjs/runtime'
import { Snackbar } from '@varlet/ui'
import { call } from '../utils/call'

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

onMounted(async () => {
    await checkLogin();
    if (!is_login.value) {
        return;
    }
    const collect = await call(GetUsersCollect(), "获取收藏夹失败");
    if (collect.ok) {
        user_collect.value = collect.data;
        CardStatus.LoadUsersCollect = true;
    }
    getFavCollect();
    const user = await call(GetUserInf(), "获取用户信息失败");
    if (user.ok) {
        user_Information.value = user.data;
        CardStatus.LoadUserInf = true;
    }
});

// 一键添加创建的收藏夹
function addCreated(collect) {
    EventsEmit('addToList', 'favlist', favlistURL(collect.id, 'ftype=create'));
}

// 一键添加收藏的内容，attr 为 0 时是视频合集
function addCollected(collect) {
    if (collect.attr == 0) {
        EventsEmit('addToList', 'season', favlistURL(collect.id, 'ftype=collect&ctype=21'));
    } else {
        EventsEmit('addToList', 'favlist', favlistURL(collect.id, 'ftype=collect&ctype=11'));
    }
}

// 拼接与网页端格式一致的收藏夹链接
function favlistURL(fid, query) {
    return 'https://space.bilibili.com/' + user_collect.value.user_mid + '/favlist?fid=' + fid + '&' + query;
}

// 获取订阅收藏夹列表
async function getFavCollect() {
    const { ok, data } = await call(GetFavCollect(page_index.value), "获取收藏夹失败");
    if (ok) {
        user_Favourite.value = data;
        fav_count.value = data.count;
        CardStatus.LoadFavCollect = true;
    }
}

// 订阅的收藏夹翻页
watch(page_index, (newValue) => {
    if (newValue <= 0) {
        page_index.value = 1;
    } else if (newValue > fav_count.value / 20 + 1) {
        page_index.value = parseInt(fav_count.value / 20 + 1);
    }
    getFavCollect();

})

// 获取二维码事件
EventsOn("qrcodeStr", (qr) => {
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
async function checkLogin() {
    const result = await LoadConfig();
    is_login.value = result.account.is_login;
}

// 登录账户
function login() {
    LoginBilibili().then(result => {
        showLoginWindow.value = false;
    })
    showLoginWindow.value = true;
}
</script>