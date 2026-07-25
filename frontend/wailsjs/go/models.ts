export namespace gui {
	
	export class DayTotalDTO {
	    date: string;
	    seconds: number;
	
	    static createFrom(source: any = {}) {
	        return new DayTotalDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.seconds = source["seconds"];
	    }
	}
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
	export class StatsDTO {
	    totalSeconds: number;
	    sessionCount: number;
	    topTag: string;
	    streakDays: number;
	
	    static createFrom(source: any = {}) {
	        return new StatsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalSeconds = source["totalSeconds"];
	        this.sessionCount = source["sessionCount"];
	        this.topTag = source["topTag"];
	        this.streakDays = source["streakDays"];
	    }
	}
	export class TagTimeDTO {
	    tag: string;
	    seconds: number;
	
	    static createFrom(source: any = {}) {
	        return new TagTimeDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tag = source["tag"];
	        this.seconds = source["seconds"];
	    }
	}

}

