export namespace main {
	
	export class GenerationResult {
	    outputPath: string;
	    totalImages: number;
	    landscapeCount: number;
	    portraitCount: number;
	
	    static createFrom(source: any = {}) {
	        return new GenerationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.outputPath = source["outputPath"];
	        this.totalImages = source["totalImages"];
	        this.landscapeCount = source["landscapeCount"];
	        this.portraitCount = source["portraitCount"];
	    }
	}
	export class ProcessUploadsResult {
	    added: string[];
	    ignored: string[];
	    totalUploads: number;
	
	    static createFrom(source: any = {}) {
	        return new ProcessUploadsResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.added = source["added"];
	        this.ignored = source["ignored"];
	        this.totalUploads = source["totalUploads"];
	    }
	}

}

