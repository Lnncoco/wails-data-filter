export namespace main {
	
	export class RegexConfig {
	    regexp: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new RegexConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.regexp = source["regexp"];
	        this.value = source["value"];
	    }
	}
	export class Options {
	    eraseZeroWidthCharacter: boolean;
	    regex: RegexConfig[];
	
	    static createFrom(source: any = {}) {
	        return new Options(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.eraseZeroWidthCharacter = source["eraseZeroWidthCharacter"];
	        this.regex = this.convertValues(source["regex"], RegexConfig);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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
	export class PathStruct {
	    directory: string;
	    filename: string;
	    ext: string;
	    addSuffixStr: string;
	
	    static createFrom(source: any = {}) {
	        return new PathStruct(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.directory = source["directory"];
	        this.filename = source["filename"];
	        this.ext = source["ext"];
	        this.addSuffixStr = source["addSuffixStr"];
	    }
	}

}

