export namespace main {
	
	export class Account {
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
	        this.sessdata = source["sessdata"];
	        this.bili_jct = source["bili_jct"];
	        this.dede_user_id = source["dede_user_id"];
	        this.dede_user_id__ck_md5 = source["dede_user_id__ck_md5"];
	        this.sid = source["sid"];
	    }
	}
	export class AudioInf {
	    code: number;
	    msg: string;
	
	    static createFrom(source: any = {}) {
	        return new AudioInf(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.msg = source["msg"];
	    }
	}
	export class Config {
	    download_path: string;
	    cache_path: string;
	    videolist_path: string;
	    download_threads: number;
	    retry_count: number;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.download_path = source["download_path"];
	        this.cache_path = source["cache_path"];
	        this.videolist_path = source["videolist_path"];
	        this.download_threads = source["download_threads"];
	        this.retry_count = source["retry_count"];
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
	export class VideoInformationList {
	    bvid: string;
	    cid: number;
	    title: string;
	    videos: number;
	    list_id: number;
	    is_page: boolean;
	    page_title: string;
	    page_id: number;
	
	    static createFrom(source: any = {}) {
	        return new VideoInformationList(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bvid = source["bvid"];
	        this.cid = source["cid"];
	        this.title = source["title"];
	        this.videos = source["videos"];
	        this.list_id = source["list_id"];
	        this.is_page = source["is_page"];
	        this.page_title = source["page_title"];
	        this.page_id = source["page_id"];
	    }
	}

}

