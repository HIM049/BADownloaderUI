export namespace main {
	
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

}

