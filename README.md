<div align="center" style="padding: 20px;">
  <img src="https://github.com/HIM049/BADownloaderUI/assets/67405384/b680bc86-5b41-4238-85ad-f50bf975bd07"/>
</div>

# Bili Audio Downloader UI - 使用 wails CLI 重构的全新版本

考虑到以往 Bili Audio Downloader 使用命令行操作的不便，于是有了这个完全重构的带 UI 版本！  
这是我第一次使用类似 wails 以及 Vue 来进行开发，经验不足，请多包涵。如有好的修改建议欢迎向我提出！  
（受个人安排影响， PR 和信息处理回复的周期大约为 7 天）  

## 下载
软件的预编译版本请移步至 [Releases](https://github.com/HIM049/BADownloaderUI/releases) 页面下载。  
如果你是 scoop 包管理器的用户，也可以在由 [Weidows](https://github.com/Weidows) 整理的软件仓库中下载使用该软件。

```
scoop bucket add apps https://github.com/kkzzhizhou/scoop-apps
scoop install BADownloaderUI
```

## 使用说明
- 下载 Bili Audio Downloader 的可执行文件，并放入到一个文件夹中
- 运行程序，程序会在目录下生成其配置文件以及缓存目录等
- 输入你希望下载的收藏夹编号或 URL（网址），在确定收藏夹信息正确后点击“下一步” **目前仅支持获取公开收藏夹下载** *（补充说明 #1）*
- 编辑下载偏好。下载数量为 0 时会下载收藏夹中的全部内容，其他数量则是按照收藏夹从前到后排序下载。元数据是音乐的标签，音乐 APP 和播放器通常会需要这些数据。打开对应的开关后程序会将对应的视频数据写入歌曲的元数据中。
- 点击 “生成视频列表” 按钮， 软件会将接下来要下载的歌曲制作成 json 格式的信息表保存在本地。
- 编辑列表内容。列表内的是接下来会下载的内容以及对应内容的元数据。你可以根据需要进行修改。
- 随后点击 “开始下载” 按钮。软件会自动完成剩余的步骤。最终歌曲会被默认输出到 `./Downloads` 文件夹中。

## 补充说明
1. 在 B 站查看收藏夹时，浏览器 URL 中靠后部分的 `fid=` 后跟随的数字部分就是收藏夹编号。如 URL 是以 `/favlist` 结尾，请点击一下希望下载的收藏夹
3. **程序目前未对大部分输入框进行输入审核，请注意输入内容符合要求！**

## 发布说明
- 目前“发布”页面内提供 Windows 平台的预编译内容。格式为 `BAdownloader-{ 版本号 }-{ 平台 }-{ 架构 }` 不了解的用户请下载后缀为 `amd64` 的软件包。

