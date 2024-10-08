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
	
	    static createFrom(source: any = {}) {
	        return new Audio(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.auid = source["auid"];
	    }
	}
	export class Collects {
	    user_mid: number;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new Collects(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user_mid = source["user_mid"];
	        this.count = source["count"];
	    }
	}
	export class CompliationInformation {
	    code: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new CompliationInformation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	    }
	}
	export class FavList {
	    code: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new FavList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	    }
	}
	export class Video {
	    bvid: string;
	
	    static createFrom(source: any = {}) {
	        return new Video(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bvid = source["bvid"];
	    }
	}
	export class Videos {
	    cid: number;
	    part: string;
	
	    static createFrom(source: any = {}) {
	        return new Videos(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cid = source["cid"];
	        this.part = source["part"];
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

}

export namespace main {
	
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
	export class AudioInformation {
	    quality: number;
	    stream: string;
	
	    static createFrom(source: any = {}) {
	        return new AudioInformation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.quality = source["quality"];
	        this.stream = source["stream"];
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
	
	export class DownloadOption {
	    song_name: boolean;
	    song_cover: boolean;
	    song_author: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DownloadOption(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.song_name = source["song_name"];
	        this.song_cover = source["song_cover"];
	        this.song_author = source["song_author"];
	    }
	}
	
	export class MetaInformation {
	    song_name: string;
	    cover: string;
	    author: string;
	    lyrics_path: string;
	
	    static createFrom(source: any = {}) {
	        return new MetaInformation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.song_name = source["song_name"];
	        this.cover = source["cover"];
	        this.author = source["author"];
	        this.lyrics_path = source["lyrics_path"];
	    }
	}
	export class VideoInformation {
	    bvid: string;
	    cid: number;
	    title: string;
	    page_title: string;
	    format: string;
	    part_id: number;
	    is_audio: boolean;
	    delete: boolean;
	
	    static createFrom(source: any = {}) {
	        return new VideoInformation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bvid = source["bvid"];
	        this.cid = source["cid"];
	        this.title = source["title"];
	        this.page_title = source["page_title"];
	        this.format = source["format"];
	        this.part_id = source["part_id"];
	        this.is_audio = source["is_audio"];
	        this.delete = source["delete"];
	    }
	}
	export class VideoList {
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new VideoList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.count = source["count"];
	    }
	}

}

