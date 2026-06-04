export namespace main {
	
	export class Tag {
	    id: number;
	    name: string;
	    color: string;
	
	    static createFrom(source: any = {}) {
	        return new Tag(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.color = source["color"];
	    }
	}
	export class Task {
	    id: number;
	    scheduleId: number;
	    rowIndex: number;
	    dayOfWeek: number;
	    startTime: string;
	    endTime: string;
	    description: string;
	    isCompleted: boolean;
	    createdAt: string;
	    updatedAt: string;
	    tags: Tag[];
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.scheduleId = source["scheduleId"];
	        this.rowIndex = source["rowIndex"];
	        this.dayOfWeek = source["dayOfWeek"];
	        this.startTime = source["startTime"];
	        this.endTime = source["endTime"];
	        this.description = source["description"];
	        this.isCompleted = source["isCompleted"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	        this.tags = this.convertValues(source["tags"], Tag);
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

