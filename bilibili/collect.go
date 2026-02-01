package bilibili

// 用于获取收藏夹基本信息的函数
// 传入收藏夹 ID ，ps 单页大小， pn 页码
// 获得如下结构体
type FavList struct {
	Code    int    `json:"code"`    // 状态码
	Message string `json:"message"` // 错误消息
	Data    struct {
		Info struct { // 收藏夹信息
			Title       string `json:"title"`       // 收藏夹标题
			Cover       string `json:"cover"`       // 收藏夹封面
			Media_count int    `json:"media_count"` // 收藏夹数量
			Upper       struct {
				Name string `json:"name"` // 创建者昵称
				Face string `json:"face"` // 创建者头像 url
			}
		}
		Medias []struct { // 收藏夹中的视频
			Id    int    `json:"id"`    // 稿件 avid
			Type  int    `json:"type"`  // 内容类型 （视频稿件2 音频12 合集21）
			Title string `json:"title"` // 标题
			Cover string `json:"cover"` // 封面 url
			Page  int    `json:"page"`  // 视频分P数
			Bvid  string `json:"bvid"`  // BV 号
		}
	}
}

// Used to be getFavList here. Moved to api.go

// 获取用户创建的收藏夹
type Collects struct {
	UserMid int `json:"user_mid"`
	Count   int `json:"count"`
	List    []meta
}
type meta struct {
	Id         int    `json:"id"`   // 收藏夹 ID
	Mid        int    `json:"mid"`  // 创建者 MID
	Attr       int    `json:"attr"` // 属性
	Title      string `json:"title"`
	Cover      string `json:"cover"`
	MediaCount int    `json:"media_count"`
}

// Used to be GetFavCollect, getUserfavoritesCollect, GetUsersCollect, getUsersCollect here. Moved to api.go
