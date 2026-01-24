package bilibili

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

// Audio.go API

// Query query audio streams
func (audio *Audio) Query(auid string) error {
	// 设置 URL 并发送 GET 请求
	urlPath := "https://www.bilibili.com/audio/music-service-c/web/song/info"

	body, err := Request(urlPath, WithParams(map[string]string{"sid": auid}))
	if err != nil {
		return err
	}
	bodyJson := string(body)

	audio.Auid = auid
	audio.Meta.Title = gjson.Get(bodyJson, "data.title").String()
	audio.Meta.Cover = gjson.Get(bodyJson, "data.cover").String()
	audio.Meta.Lyric = gjson.Get(bodyJson, "data.lyric").String()
	audio.Up.Author = gjson.Get(bodyJson, "data.author").String()

	return nil
}

// GetStream get audio stream url
func (audio *Audio) GetStream(sessdata string) error {
	// 创建请求
	urlStr := "https://api.bilibili.com/audio/music-service-c/url"

	// 设置 URL 参数
	params := map[string]string{
		"songid":    audio.Auid,
		"quality":   "2",
		"privilege": "2",
		"mid":       "2",
		"platform":  "web",
	}

	body, err := Request(urlStr, WithSESSDATA(sessdata), WithParams(params))
	if err != nil {
		return err
	}
	bodyJson := string(body)

	// 错误检查
	if CheckObj(int(gjson.Get(bodyJson, "code").Int())) {
		return errors.New(gjson.Get(bodyJson, "message").String())
	}

	audio.Stream.Type = int(gjson.Get(bodyJson, "data.type").Int())
	audio.Stream.StreamLink = gjson.Get(bodyJson, "data.cdns.0").String()

	return nil
}

// Video.go API

// 请求视频详细信息
// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/video/info.md
// TODO：重新添加字幕信息
func (v *Video) Query(sessdata, bvid string) error {
	// 创建请求
	urlStr := "https://api.bilibili.com/x/web-interface/view"

	body, err := Request(urlStr, WithSESSDATA(sessdata), WithParams(map[string]string{"bvid": bvid}))
	if err != nil {
		return err
	}

	json := string(body)

	// 将信息写入结构体
	v.Bvid = bvid
	v.Meta.Title = gjson.Get(json, "data.title").String()                        // 视频标题
	v.Meta.Cover = gjson.Get(json, "data.pic").String()                          // 视频封面
	v.Meta.LyricsPath = gjson.Get(json, "data.subtitle.0.subtitle_url").String() // 字幕获取（临时）
	v.Up.Mid = int(gjson.Get(json, "data.owner.mid").Int())                      // UP MID
	v.Up.Name = gjson.Get(json, "data.owner.name").String()                      // UP 昵称
	v.Up.Avatar = gjson.Get(json, "data.owner.face").String()                    // UP 头像

	// 根据分 P 数量写入对应信息
	for i := 0; i < int(gjson.Get(json, "data.videos").Int()); i++ {

		// 单个分集视频信息
		videos := Videos{
			Cid:  int(gjson.Get(json, "data.pages."+strconv.Itoa(i)+".cid").Int()),
			Part: gjson.Get(json, "data.pages."+strconv.Itoa(i)+".part").String(),
		}
		v.Videos = append(v.Videos, videos)
	}

	return nil
}

// 获取视频流
// https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/video/videostream_url.md#%E8%8E%B7%E5%8F%96%E8%A7%86%E9%A2%91%E6%B5%81%E5%9C%B0%E5%9D%80_web%E7%AB%AF
func GetVideoStream(bvid, cid, sessdata string) (string, error) {
	// 创建请求
	urlStr := "https://api.bilibili.com/x/player/wbi/playurl"

	// 设置 URL 参数
	params := map[string]string{
		"bvid":  bvid,
		"cid":   cid,
		"fnval": "16",
	}

	// This particular API uses WbiSignURLParams which was called manually in the original code.
	// Since we use the Wbi option, it should be handled by requests.go logic if we use WithWbi().
	// However, Request() with WithWbi() will sign the params.
	// Original code:
	// signedUrl, err := WbiSignURLParams(request.URL.String())
	// Then creates new request with signedUrl.

	body, err := Request(urlStr, WithParams(params), WithWbi(), WithSESSDATA(sessdata))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// login.go API

// 获取登录密钥
func GetLoginKey() (string, string, error) {
	var obj GetLoginKeyReturn
	body, err := getLoginKey()
	if err != nil {
		return "", "", err
	}
	err = decodeJson(body, &obj)
	if err != nil {
		return "", "", err
	}
	// 错误检查
	if CheckObj(obj.Code) {
		return "", "", errors.New(obj.Message)
	}
	return obj.Data.Url, obj.Data.Qrcode_key, nil
}

// 请求登录密钥
func getLoginKey() (string, error) {
	body, err := Request("https://passport.bilibili.com/x/passport-login/web/qrcode/generate")
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// 检查扫码状态
func checkLoginStatus(qrcode_key string) (string, *[]*http.Cookie, error) {
	// Original code manually handles cookies from response.
	// My Request function returns body, error. It does not return the response object to get cookies.
	// We need to modify Request to return (body, cookies, error) or exposing response?
	// Or handle this case specifically.
	// The requirement "unify request methods" suggests using the standard way.
	// But getting cookies from response is a specific requirement here.

	// Let's implement checkLoginStatus manually here using http.Client to preserve functionality,
	// OR modify Request to return *http.Response (and let user read body)?
	// Or maybe a specialized request in requests.go?

	// Since this is a specific case (polling login status and getting Set-Cookie),
	// I will keep the http.Client usage but try to make it cleaner or just adapt it slightly.
	// Actually, the user asked to "classify request methods and unify to SAME call method (closure call)".
	// If I can't fit it into `Request`, I should extend `Request` or add `RequestWithResponse`.

	// Let's add `RequestWithCookies` to requests.go?
	// Or just use the original logic but inside this file?
	// The prompt implies "unify" -> so I should probably make `Request` capable of returning cookies or headers.
	// But `Request` currently returns `[]byte`.

	// For now, I will manually implement this using `http.Client` here but structure it similarly,
	// OR I will modify `requests.go` to return `(*http.Response, []byte, error)`? That seems too heavy for simple requests.

	// I'll stick to manual implementation for this one but minimalize it, OR add a `RequestRaw` to `requests.go`.
	// Let's stick to standard `http.Client` here for now to avoid over-engineering `requests.go` for one edge case,
	// unless `requests.go` can be easily modified.
	// Actually, `checkLoginStatus` returns `*[]*http.Cookie`.

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://passport.bilibili.com/x/passport-login/web/qrcode/poll", nil)
	if err != nil {
		return "", nil, err
	}

	q := req.URL.Query()
	q.Add("qrcode_key", qrcode_key)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, errors.New("Error:" + strconv.Itoa(resp.StatusCode))
	}

	cookies := resp.Cookies()
	body, _ := io.ReadAll(resp.Body) // Simplify error handling as original code ignored it? No, original: _
	bodyString := string(body)
	return bodyString, &cookies, nil
}

func CheckLoginStatus(qrcode_key string) (*checkLoginReturn, *[]*http.Cookie, error) {
	var obj checkLoginReturn
	body, cookies, err := checkLoginStatus(qrcode_key)
	if err != nil {
		return nil, nil, err
	}
	err = decodeJson(body, &obj)
	if err != nil {
		return nil, nil, err
	}
	// 错误检查
	if CheckObj(obj.Code) {
		return nil, nil, errors.New(obj.Message)
	}

	return &obj, cookies, nil
}

// 获取用户信息
// https://socialsisteryi.github.io/bilibili-API-collect/docs/login/login_info.html
func (accountInf *AccountInformation) GetUserInf(sessdata string) error {
	body, err := Request("https://api.bilibili.com/x/web-interface/nav", WithSESSDATA(sessdata))
	if err != nil {
		return err
	}
	bodyJson := string(body)

	// 错误检查
	if CheckObj(int(gjson.Get(bodyJson, "code").Int())) {
		return errors.New(gjson.Get(bodyJson, "message").String())
	}

	accountInf.Avatar = gjson.Get(bodyJson, "data.face").String()
	accountInf.Name = gjson.Get(bodyJson, "data.uname").String()

	return nil
}

// collect.go API

func getFavList(id, ps, pn, sessdata string) (string, error) {
	params := map[string]string{
		"media_id": id,
		"ps":       ps,
		"pn":       pn,
		"platform": "web",
	}
	body, err := Request("https://api.bilibili.com/x/v3/fav/resource/list", WithSESSDATA(sessdata), WithParams(params))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func GetFavListObj(id, sessdata string, ps, pn int) (*FavList, error) {
	var obj FavList
	body, err := getFavList(id, strconv.Itoa(ps), strconv.Itoa(pn), sessdata)
	if err != nil {
		return nil, err
	}
	err = decodeJson(body, &obj)
	if err != nil {
		return nil, err
	}
	// 错误检查
	if CheckObj(obj.Code) {
		return nil, errors.New(obj.Message)
	}
	return &obj, nil
}

// 获取用户收藏的收藏夹
func (collects *Collects) GetFavCollect(sessdata string, ps, pn int) error {
	json, err := getUserfavoritesCollect(sessdata, strconv.Itoa(collects.UserMid), strconv.Itoa(ps), strconv.Itoa(pn))
	if err != nil {
		return err
	}

	// 错误检查
	if CheckObj(int(gjson.Get(json, "code").Int())) {
		return errors.New(gjson.Get(json, "message").String())
	}

	collects.Count = int(gjson.Get(json, "data.count").Int())
	pageCount := collects.Count

	if collects.Count/20 >= pn {
		pageCount = 20
	} else {
		pageCount = collects.Count % 20
	}

	for i := 0; i < pageCount; i++ {
		meta := new(meta)
		meta.Id = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".id").Int())
		meta.Mid = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".mid").Int())
		meta.Attr = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".attr").Int())
		meta.Title = gjson.Get(json, "data.list."+strconv.Itoa(i)+".title").String()
		meta.Cover = gjson.Get(json, "data.list."+strconv.Itoa(i)+".cover").String()
		meta.MediaCount = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".media_count").Int())
		collects.List = append(collects.List, *meta)
	}

	return nil
}

// 获取用户收藏的收藏夹
func getUserfavoritesCollect(sessdata, mid, pageSize, pageNumber string) (string, error) {
	params := map[string]string{
		"ps":       pageSize,
		"pn":       pageNumber,
		"up_mid":   mid,
		"platform": "web",
	}
	body, err := Request("https://api.bilibili.com/x/v3/fav/folder/collected/list", WithSESSDATA(sessdata), WithParams(params))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// 获取用户创建的收藏夹
func (collects *Collects) GetUsersCollect(sessdata string) error {
	json, err := getUsersCollect(sessdata, strconv.Itoa(collects.UserMid))
	if err != nil {
		return err
	}

	// 错误检查
	if CheckObj(int(gjson.Get(json, "code").Int())) {
		return errors.New(gjson.Get(json, "message").String())
	}

	collects.Count = int(gjson.Get(json, "data.count").Int())
	for i := 0; i < collects.Count; i++ {
		meta := new(meta)
		meta.Id = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".id").Int())
		meta.Mid = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".mid").Int())
		meta.Attr = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".attr").Int())
		meta.Title = gjson.Get(json, "data.list."+strconv.Itoa(i)+".title").String()
		meta.MediaCount = int(gjson.Get(json, "data.list."+strconv.Itoa(i)+".media_count").Int())
		collects.List = append(collects.List, *meta)
	}

	return nil
}

// 获取用户创建的收藏夹
func getUsersCollect(sessdata, mid string) (string, error) {
	params := map[string]string{
		"up_mid":   mid,
		"platform": "web",
	}
	body, err := Request("https://api.bilibili.com/x/v3/fav/folder/created/list-all", WithSESSDATA(sessdata), WithParams(params))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// wbi.go API

func getWbiKeys() (string, string) {
	// REPLACEMENT calling Request:
	body, err := Request("https://api.bilibili.com/x/web-interface/nav")
	if err != nil {
		return "", ""
	}

	json := string(body)
	imgURL := gjson.Get(json, "data.wbi_img.img_url").String()
	subURL := gjson.Get(json, "data.wbi_img.sub_url").String()
	// Check if imgURL/subURL are empty to avoid panic on Split? Original code didn't check.
	if imgURL == "" || subURL == "" {
		return "", ""
	}

	return parseWbiKeys(imgURL, subURL)
}

func parseWbiKeys(imgURL, subURL string) (string, string) {
	imgKey := strings.Split(strings.Split(imgURL, "/")[len(strings.Split(imgURL, "/"))-1], ".")[0]
	subKey := strings.Split(strings.Split(subURL, "/")[len(strings.Split(subURL, "/"))-1], ".")[0]
	return imgKey, subKey
}

// compilation.go API

func getCompliation(mid, sid, ps, pn string) (string, error) {
	params := map[string]string{
		"mid":       mid,
		"season_id": sid,
		"page_size": ps,
		"page_num":  pn,
	}
	// Original code sets referer manually, Request() does it too.
	body, err := Request("https://api.bilibili.com/x/polymer/web-space/seasons_archives_list", WithParams(params))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func GetCompliationObj(mid, sid, ps, pn int) (*CompliationInformation, error) {
	var obj CompliationInformation
	body, err := getCompliation(strconv.Itoa(mid), strconv.Itoa(sid), strconv.Itoa(ps), strconv.Itoa(pn))
	if err != nil {
		return nil, err
	}
	err = decodeJson(body, &obj)
	if err != nil {
		return nil, err
	}

	// 错误检查
	if CheckObj(obj.Code) {
		return nil, errors.New(obj.Message)
	}
	return &obj, nil
}

// profile.go API

// 获取用户投稿列表
// https://socialsisteryi.github.io/bilibili-API-collect/docs/user/space.html#%E6%9F%A5%E8%AF%A2%E7%94%A8%E6%88%B7%E6%8A%95%E7%A8%BF%E8%A7%86%E9%A2%91%E6%98%8E%E7%BB%86
func GetProfileVideo(mid, pn, ps, sessdata string) (string, error) {
	params := map[string]string{
		"mid":   mid,
		"order": "pubdate",
		"pn":    pn,
		"ps":    ps,
	}
	// This uses Wbi signing!
	body, err := Request("https://api.bilibili.com/x/space/wbi/arc/search", WithParams(params), WithWbi(), WithSESSDATA(sessdata))
	if err != nil {
		return "", err
	}
	return string(body), nil
}
