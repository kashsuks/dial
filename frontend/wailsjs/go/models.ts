export namespace gui {
	
	export class SessionDTO {
	    id: number;
	    task: string;
	    project: string;
	    tags: string;
	    startedAt: string;
	    endedAt?: string;
	    isPaused: boolean;
	    elapsedSeconds: number;
	
	    static createFrom(source: any = {}) {
	        return new SessionDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.task = source["task"];
	        this.project = source["project"];
	        this.tags = source["tags"];
	        this.startedAt = source["startedAt"];
	        this.endedAt = source["endedAt"];
	        this.isPaused = source["isPaused"];
	        this.elapsedSeconds = source["elapsedSeconds"];
	    }
	}

}

