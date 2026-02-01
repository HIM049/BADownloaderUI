package bilibili

// 用于获取收藏夹基本信息的函数
// 传入收藏夹 ID ，ps 单页大小， pn 页码
type CompliationInformation struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Archives []struct {
			Bvid  string `json:"bvid"`
			Pic   string `json:"pic"`
			Title string `json:"title"`
		}
		Meta struct {
			Cover       string `json:"cover"`
			Description string `json:"description"`
			Name        string `json:"name"`
			Total       int    `json:"total"`
		}
	}
}

// Used to be getCompliation and GetCompliationObj here. Moved to api.go
