export namespace biliapi {
	
	export class FavFolder {
	    id: number;
	    mid: number;
	    attr: number;
	    title: string;
	    cover: string;
	    media_count: number;
	
	    static createFrom(source: any = {}) {
	        return new FavFolder(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.mid = source["mid"];
	        this.attr = source["attr"];
	        this.title = source["title"];
	        this.cover = source["cover"];
	        this.media_count = source["media_count"];
	    }
	}

}

export namespace config {
	
	export class Account {
	    is_login: boolean;
	    use_account: boolean;
	    sessdata: string;
	    bili_jct: string;
	    dede_user_id: string;
	    dede_user_id__ck_md5: string;
	    sid: string;
	
	    static createFrom(source: any = {}) {
	        return new Account(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.is_login = source["is_login"];
	        this.use_account = source["use_account"];
	        this.sessdata = source["sessdata"];
	        this.bili_jct = source["bili_jct"];
	        this.dede_user_id = source["dede_user_id"];
	        this.dede_user_id__ck_md5 = source["dede_user_id__ck_md5"];
	        this.sid = source["sid"];
	    }
	}
	export class FileConfig {
	    convert_format: boolean;
	    file_name_template: string;
	    download_path: string;
	    cache_path: string;
	
	    static createFrom(source: any = {}) {
	        return new FileConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.convert_format = source["convert_format"];
	        this.file_name_template = source["file_name_template"];
	        this.download_path = source["download_path"];
	        this.cache_path = source["cache_path"];
	    }
	}
	export class DownloadConfig {
	    download_threads: number;
	    retry_count: number;
	
	    static createFrom(source: any = {}) {
	        return new DownloadConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.download_threads = source["download_threads"];
	        this.retry_count = source["retry_count"];
	    }
	}
	export class Config {
	    config_version: number;
	    delete_cache: boolean;
	    theme: string;
	    download_config: DownloadConfig;
	    file_config: FileConfig;
	    Account: Account;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.config_version = source["config_version"];
	        this.delete_cache = source["delete_cache"];
	        this.theme = source["theme"];
	        this.download_config = this.convertValues(source["download_config"], DownloadConfig);
	        this.file_config = this.convertValues(source["file_config"], FileConfig);
	        this.Account = this.convertValues(source["Account"], Account);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

export namespace task {
	
	export class MetaData {
	    Title: string;
	    PageTitle: string;
	    PartId: number;
	    SongName: string;
	    Author: string;
	    LyricsUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new MetaData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Title = source["Title"];
	        this.PageTitle = source["PageTitle"];
	        this.PartId = source["PartId"];
	        this.SongName = source["SongName"];
	        this.Author = source["Author"];
	        this.LyricsUrl = source["LyricsUrl"];
	    }
	}
	export class TaskInfo {
	    Index: number;
	    SongName: string;
	    SongAuthor: string;
	    CoverUrl: string;
	    IsDelete: boolean;
	    type?: number;
	    bvid?: string;
	    cid?: number;
	    auid?: string;
	    sessdata?: string;
	    meta_data?: MetaData;
	
	    static createFrom(source: any = {}) {
	        return new TaskInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Index = source["Index"];
	        this.SongName = source["SongName"];
	        this.SongAuthor = source["SongAuthor"];
	        this.CoverUrl = source["CoverUrl"];
	        this.IsDelete = source["IsDelete"];
	        this.type = source["type"];
	        this.bvid = source["bvid"];
	        this.cid = source["cid"];
	        this.auid = source["auid"];
	        this.sessdata = source["sessdata"];
	        this.meta_data = this.convertValues(source["meta_data"], MetaData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace wails_api {
	
	export class AccountInfo {
	    avatar: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new AccountInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.avatar = source["avatar"];
	        this.name = source["name"];
	    }
	}
	export class MediaMeta {
	    title: string;
	    cover: string;
	
	    static createFrom(source: any = {}) {
	        return new MediaMeta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.cover = source["cover"];
	    }
	}
	export class AudioInfo {
	    Meta: MediaMeta;
	    // Go type: struct { Author string "json:\"author\"" }
	    Up: any;
	
	    static createFrom(source: any = {}) {
	        return new AudioInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Meta = this.convertValues(source["Meta"], MediaMeta);
	        this.Up = this.convertValues(source["Up"], Object);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CompilationInfo {
	    // Go type: struct { Meta struct { Name string "json:\"name\""; Cover string "json:\"cover\""; Total int "json:\"total\"" } "json:\"Meta\"" }
	    Data: any;
	
	    static createFrom(source: any = {}) {
	        return new CompilationInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Data = this.convertValues(source["Data"], Object);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FavListInfo {
	    // Go type: struct { Info struct { Title string "json:\"title\""; Cover string "json:\"cover\""; MediaCount int "json:\"media_count\""; Upper struct { Name string "json:\"name\""; Face string "json:\"face\"" } "json:\"Upper\"" } "json:\"Info\"" }
	    Data: any;
	
	    static createFrom(source: any = {}) {
	        return new FavListInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Data = this.convertValues(source["Data"], Object);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FolderList {
	    user_mid: number;
	    count: number;
	    List: biliapi.FavFolder[];
	
	    static createFrom(source: any = {}) {
	        return new FolderList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user_mid = source["user_mid"];
	        this.count = source["count"];
	        this.List = this.convertValues(source["List"], biliapi.FavFolder);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class VideoInfo {
	    Meta: MediaMeta;
	    // Go type: struct { Name string "json:\"name\""; Avatar string "json:\"avatar\"" }
	    Up: any;
	
	    static createFrom(source: any = {}) {
	        return new VideoInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Meta = this.convertValues(source["Meta"], MediaMeta);
	        this.Up = this.convertValues(source["Up"], Object);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

