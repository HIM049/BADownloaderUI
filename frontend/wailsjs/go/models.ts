export namespace adapter {
	
	export class TaskInfo {
	    Index: number;
	    SongName: string;
	    SongAuthor: string;
	    CoverUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new TaskInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Index = source["Index"];
	        this.SongName = source["SongName"];
	        this.SongAuthor = source["SongAuthor"];
	        this.CoverUrl = source["CoverUrl"];
	    }
	}

}

export namespace bilibili {
	
	export class AccountInformation {
	    avatar: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new AccountInformation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.avatar = source["avatar"];
	        this.name = source["name"];
	    }
	}
	export class Audio {
	    auid: string;
	    // Go type: struct { Title string "json:\"title\""; Cover string "json:\"cover\""; Lyric string "json:\"lyric\"" }
	    Meta: any;
	    // Go type: struct { Author string "json:\"author\"" }
	    Up: any;
	    // Go type: struct { Type int "json:\"type\""; StreamLink string "json:\"stream_link\"" }
	    Stream: any;
	
	    static createFrom(source: any = {}) {
	        return new Audio(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.auid = source["auid"];
	        this.Meta = this.convertValues(source["Meta"], Object);
	        this.Up = this.convertValues(source["Up"], Object);
	        this.Stream = this.convertValues(source["Stream"], Object);
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
	export class meta {
	    id: number;
	    mid: number;
	    attr: number;
	    title: string;
	    cover: string;
	    media_count: number;
	
	    static createFrom(source: any = {}) {
	        return new meta(source);
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
	export class Collects {
	    user_mid: number;
	    count: number;
	    List: meta[];
	
	    static createFrom(source: any = {}) {
	        return new Collects(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user_mid = source["user_mid"];
	        this.count = source["count"];
	        this.List = this.convertValues(source["List"], meta);
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
	export class CompliationInformation {
	    code: number;
	    message: string;
	    // Go type: struct { Archives []struct { Bvid string "json:\"bvid\""; Pic string "json:\"pic\""; Title string "json:\"title\"" }; Meta struct { Cover string "json:\"cover\""; Description string "json:\"description\""; Name string "json:\"name\""; Total int "json:\"total\"" } }
	    Data: any;
	
	    static createFrom(source: any = {}) {
	        return new CompliationInformation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
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
	export class FavList {
	    code: number;
	    message: string;
	    // Go type: struct { Info struct { Title string "json:\"title\""; Cover string "json:\"cover\""; Media_count int "json:\"media_count\""; Upper struct { Name string "json:\"name\""; Face string "json:\"face\"" } }; Medias []struct { Id int "json:\"id\""; Type int "json:\"type\""; Title string "json:\"title\""; Cover string "json:\"cover\""; Page int "json:\"page\""; Bvid string "json:\"bvid\"" } }
	    Data: any;
	
	    static createFrom(source: any = {}) {
	        return new FavList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
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
	export class Videos {
	    cid: number;
	    part: string;
	    // Go type: struct { SongName string "json:\"song_name\"" }
	    Meta: any;
	
	    static createFrom(source: any = {}) {
	        return new Videos(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cid = source["cid"];
	        this.part = source["part"];
	        this.Meta = this.convertValues(source["Meta"], Object);
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
	export class Video {
	    bvid: string;
	    // Go type: struct { Title string "json:\"title\""; Cover string "json:\"cover\""; Author string "json:\"author\""; LyricsPath string "json:\"lyrics_path\"" }
	    Meta: any;
	    // Go type: struct { Mid int "json:\"mid\""; Name string "json:\"name\""; Avatar string "json:\"avatar\"" }
	    Up: any;
	    Videos: Videos[];
	
	    static createFrom(source: any = {}) {
	        return new Video(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bvid = source["bvid"];
	        this.Meta = this.convertValues(source["Meta"], Object);
	        this.Up = this.convertValues(source["Up"], Object);
	        this.Videos = this.convertValues(source["Videos"], Videos);
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
	    videolist_path: string;
	
	    static createFrom(source: any = {}) {
	        return new FileConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.convert_format = source["convert_format"];
	        this.file_name_template = source["file_name_template"];
	        this.download_path = source["download_path"];
	        this.cache_path = source["cache_path"];
	        this.videolist_path = source["videolist_path"];
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

