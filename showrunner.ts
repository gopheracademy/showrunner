export class Client {
    conferences: conferences.ServiceClient;

    constructor(environment: string) {
        const base = new BaseClient(environment)
        this.conferences = new conferences.ServiceClient(base)
    }
}

export namespace conferences {
    export interface ConferenceSlot {
        ID: number;
        Name: string;
        Description: string;
        Cost: number;
        Capacity: number;
        StartDate: string;
        EndDate: string;

        /**
         * DependsOn means that these two Slots need to be acquired together, user must either buy
         * both Slots or pre-own one of the one it depends on.
         * DependsOn *ConferenceSlot // Currently removed as it broke encore
         * PurchaseableFrom indicates when this item is on sale, for instance early bird tickets are the first
         * ones to go on sale.
         */
        PurchaseableFrom: string;

        /**
         * PuchaseableUntil indicates when this item stops being on sale, for instance early bird tickets can
         * no loger be purchased N months before event.
         */
        PurchaseableUntil: string;

        /**
         * AvailableToPublic indicates is this is something that will appear on the tickets purchase page (ie, we can
         * issue sponsor tickets and those cannot be bought individually)
         */
        AvailableToPublic: boolean;
        Location: Location;
        ConferenceID: number;
    }

    export interface GetCurrentByEventParams {
        EventID: number;
    }

    export interface SponsorContactInformation {
        ID: number;
        Name: string;
        Role: number;
        Email: string;
        Phone: string;
    }

    export interface GetConferenceSlotsResponse {
        ConferenceSlots: ConferenceSlot[];
    }

    export interface Event {
        ID: number;
        Name: string;
        Slug: string;
        Conferences: Conference[];
    }

    export interface GetAllParams {
    }

    export interface GetAllResponse {
        Events: Event[];
    }

    export interface UpdateSponsorContactParams {
        SponsorContactInformation: SponsorContactInformation;
    }

    export interface GetConferenceSponsorsParams {
        ConferenceID: number;
    }

    export interface GetConferenceSlotsParams {
        ConferenceID: number;
    }

    export interface Location {
        ID: number;
        Name: string;
        Description: string;
        Address: string;
        Directions: string;
        GoogleMapsURL: string;
        Capacity: number;
        VenueID: number;
    }

    export interface GetConferenceSponsorsResponse {
        Sponsors: Sponsor[];
    }

    export interface Sponsor {
        ID: number;
        Name: string;
        Address: string;
        Website: string;
        SponsorshipLevel: number;
        Contacts: SponsorContactInformation[];
        ConferenceID: number;
    }

    export interface GetCurrentByEventResponse {
        Event: Event;
    }

    export interface Conference {
        ID: number;
        Name: string;
        Slug: string;
        StartDate: string;
        EndDate: string;
        Venue: Venue;
    }

    export interface Venue {
        ID: number;
        Name: string;
        Description: string;
        Address: string;
        Directions: string;
        GoogleMapsURL: string;
        Capacity: number;
    }

    export interface UpdateSponsorContactResponse {
    }

    export class ServiceClient {
        private baseClient: BaseClient;

        constructor(baseClient: BaseClient) {
            this.baseClient = baseClient
        }

        /**
         * GetConferenceSlots retrieves all event slots for a specific event id
         */
        public GetConferenceSlots(params: GetConferenceSlotsParams): Promise<GetConferenceSlotsResponse> {
            return this.baseClient.do<GetConferenceSlotsResponse>("conferences.GetConferenceSlots", params);
        }

        /**
         * GetCurrentByEvent retrieves the current conference and event information for a specific event
         */
        public GetCurrentByEvent(params: GetCurrentByEventParams): Promise<GetCurrentByEventResponse> {
            return this.baseClient.do<GetCurrentByEventResponse>("conferences.GetCurrentByEvent", params);
        }

        /**
         * GetAll retrieves all conferences and events
         */
        public GetAll(params: GetAllParams): Promise<GetAllResponse> {
            return this.baseClient.do<GetAllResponse>("conferences.GetAll", params);
        }

        /**
         * UpdateSponsorContact retrieves all conferences and events
         */
        public UpdateSponsorContact(params: UpdateSponsorContactParams): Promise<UpdateSponsorContactResponse> {
            return this.baseClient.do<UpdateSponsorContactResponse>("conferences.UpdateSponsorContact", params);
        }

        /**
         * GetConferenceSponsors retrieves the sponsors for a specific conference
         */
        public GetConferenceSponsors(params: GetConferenceSponsorsParams): Promise<GetConferenceSponsorsResponse> {
            return this.baseClient.do<GetConferenceSponsorsResponse>("conferences.GetConferenceSponsors", params);
        }
    }
}

class BaseClient {
    baseURL: string;

    constructor(environment: string) {
		if (environment === "dev") {
			this.baseURL = "http://localhost:4060/"
		} else {
			this.baseURL = `https://showrunner-46b2.encoreapi.com/${environment}/`
		}
    }

    public async do<T>(endpoint: string, req?: any): Promise<T> {
        let response = await fetch(this.baseURL + endpoint, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(req || {})
        })
        if (!response.ok) {
            let body = await response.text()
            throw new Error("request failed: " + body)
        }
        return <T>(await response.json())
    }

    public async doVoid(endpoint: string, req?: any): Promise<void> {
        let response = await fetch(this.baseURL + endpoint, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(req || {})
        })
        if (!response.ok) {
            let body = await response.text()
            throw new Error("request failed: " + body)
        }
        await response.text()
    }
}

const client = new Client("azure")
export default client
