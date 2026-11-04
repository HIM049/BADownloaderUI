<template>
    <FramePage title="个人空间">
        <var-space class="items-center" justify="center" v-if="!showLogin && !isLogin">
            你还没有登录
            <var-button type="primary" @click="login">登录</var-button>
        </var-space>
        <var-collapse-transition :expand="showLogin">
            <var-paper class="surface-container">
                <h3>扫码登录</h3>
                <img class="flex mx-auto rounded-[15px]" :src="qrcode" alt="Image">
                <p class="flex justify-center">{{ loginText }}</p>
            </var-paper>
        </var-collapse-transition>
        <var-space v-if="isLogin && loaded.account" justify="center">
            <var-avatar :src="account.avatar" />
            <p class="text-md font-bold my-3">{{ account.name }}</p>
        </var-space>
    </FramePage>

    <AdditionCard title="创建的收藏夹" v-if="isLogin && loaded.created">
        <FolderList :folders="created.List" @add="addCreated" />
    </AdditionCard>

    <AdditionCard title="收藏和订阅" v-if="isLogin && loaded.favorites">
        <FolderList :folders="favorites.List" @add="addCollected" />
        <var-paper class="surface-container mt-4 !py-2.5 flex justify-center">
            <var-pagination v-model:current="favPage" :total="favorites.count" :size="pageSize" :show-size-changer="false" :simple="false" />
        </var-paper>
    </AdditionCard>
</template>

<script setup lang="ts">
import FramePage from '../components/modules/frame_page.vue'
import AdditionCard from './modules/addition_card.vue'
import FolderList from './modules/folder_list.vue'
import { ref, reactive, onMounted, watch } from 'vue'
import { LoadConfig, GetUsersCollect, GetFavCollect, GetUserInf, LoginBilibili } from '../../wailsjs/go/wails_api/WailsApi'
import { biliapi, wails_api } from '../../wailsjs/go/models'
import { onEvent, emitEvent } from '../events'
import { Snackbar } from '@varlet/ui'
import { call } from '../utils/call'

// 与后端 GetFavCollect 的每页数量一致
const pageSize = 20

const loginText = ref('请扫描二维码登录') // 登录时的提示文字
const qrcode = ref('') // 二维码图片
const showLogin = ref(false) // 展示登录窗口
const isLogin = ref(false) // 登录状态
const favPage = ref(1) // 订阅的收藏夹的页码

const created = ref(new wails_api.FolderList({ List: [] })) // 用户创建的收藏夹列表
const favorites = ref(new wails_api.FolderList({ List: [] })) // 用户订阅的收藏夹列表
const account = ref(new wails_api.AccountInfo()) // 用户信息

// 各卡片的数据是否已加载
const loaded = reactive({
    created: false,
    favorites: false,
    account: false,
})

onMounted(async () => {
    await checkLogin()
    if (!isLogin.value) {
        return
    }
    const collect = await call(GetUsersCollect(), '获取收藏夹失败')
    if (collect.ok) {
        created.value = collect.data
        loaded.created = true
    }
    getFavCollect()
    const user = await call(GetUserInf(), '获取用户信息失败')
    if (user.ok) {
        account.value = user.data
        loaded.account = true
    }
})

// 一键添加创建的收藏夹
function addCreated(folder: biliapi.FavFolder) {
    emitEvent('addToList', 'favlist', favlistURL(folder.id, 'ftype=create'))
}

// 一键添加收藏的内容，attr 为 0 时是视频合集
function addCollected(folder: biliapi.FavFolder) {
    if (folder.attr == 0) {
        emitEvent('addToList', 'season', favlistURL(folder.id, 'ftype=collect&ctype=21'))
    } else {
        emitEvent('addToList', 'favlist', favlistURL(folder.id, 'ftype=collect&ctype=11'))
    }
}

// 拼接与网页端格式一致的收藏夹链接
function favlistURL(fid: number, query: string) {
    return 'https://space.bilibili.com/' + created.value.user_mid + '/favlist?fid=' + fid + '&' + query
}

// 获取当前页的订阅收藏夹
async function getFavCollect() {
    const { ok, data } = await call(GetFavCollect(favPage.value), '获取收藏夹失败')
    if (ok) {
        favorites.value = data
        loaded.favorites = true
    }
}

watch(favPage, getFavCollect)

// 获取二维码事件
onEvent('qrcodeStr', (qr) => {
    qrcode.value = 'data:image/png;base64,' + qr
})

// 登录状态事件
onEvent('loginStatus', (status) => {
    if (status == '登录成功') {
        Snackbar.success('登录成功')
        checkLogin()
        setTimeout(() => {
            window.location.reload()
        }, 1000)
    }
    if (status == '二维码已失效') {
        Snackbar.warning('二维码已失效')
        showLogin.value = false
    }
    loginText.value = status
})

// 检查登录状态
async function checkLogin() {
    const result = await LoadConfig()
    isLogin.value = result.account.is_login
}

// 登录账户
function login() {
    LoginBilibili().then(() => {
        showLogin.value = false
    })
    showLogin.value = true
}
</script>
