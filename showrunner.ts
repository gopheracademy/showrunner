export default class Client {
    conferences: conferences.ServiceClient;

    constructor(environment: string = "production", token?: string) {
        const base = new BaseClient(environment, token)
        this.conferences = new conferences.ServiceClient(base)
    }
}

export namespace conferences {
    /**
     * Conference is an instance like GopherCon 2020
     */
    export interface Conference {
        ID: number;
        Name: string;
        Slug: string;
        StartDate: string;
        EndDate: string;
        Venue: Venue;
    }

    /**
     * ConferenceSlot holds information for any sellable/giftable slot we have in the event for
     * a Talk or any other activity that requires admission.
     * store: "interface"
     */
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
         */
        DependsOn: number;

        /**
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

    /**
     * ContactRole defines the type that encapsulates the different contact roles
     */
    export interface ContactRole number

    /**
     * Event is a brand like GopherCon
     */
    export interface Event {
        ID: number;
        Name: string;
        Slug: string;
        Conferences: Conference[];
    }

    /**
     * GetAllParams defines the inputs used by the GetAll API method
     */
    export interface GetAllParams {
    }

    /**
     * GetAllResponse defines the output returned by the GetAll API method
     */
    export interface GetAllResponse {
        Events: Event[];
    }

    /**
     * GetConferenceSlotsParams defines the inputs used by the GetConferenceSlots API method
     */
    export interface GetConferenceSlotsParams {
        ConferenceID: number;
    }

    /**
     * GetConferenceSlotsResponse defines the output returned by the GetConferenceSlots API method
     */
    export interface GetConferenceSlotsResponse {
        ConferenceSlots: ConferenceSlot[];
    }

    /**
     * GetConferenceSponsorsParams defines the inputs used by the GetConferenceSponsors API method
     */
    export interface GetConferenceSponsorsParams {
        ConferenceID: number;
    }

    /**
     * GetConferenceSponsorsResponse defines the output returned by the GetConferenceSponsors API method
     */
    export interface GetConferenceSponsorsResponse {
        Sponsors: Sponsor[];
    }

    /**
     * GetCurrentByEventParams defines the inputs used by the GetCurrentByEvent API method
     */
    export interface GetCurrentByEventParams {
        EventID: number;
    }

    /**
     * GetCurrentByEventResponse defines the output returned by the GetCurrentByEvent API method
     */
    export interface GetCurrentByEventResponse {
        Event: Event;
    }

    /**
     * Location defines a location for a venue, such as a room or event space
     */
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

    /**
     * Sponsor defines a conference sponsor, such as Google
     */
    export interface Sponsor {
        ID: number;
        Name: string;
        Address: string;
        Website: string;
        SponsorshipLevel: SponsorshipLevel;
        Contacts: SponsorContactInformation[];
        ConferenceID: number;
    }

    /**
     * SponsorContactInformation defines a contact
     * and their information for a sponsor
     */
    export interface SponsorContactInformation {
        ID: number;
        Name: string;
        Role: ContactRole;
        Email: string;
        Phone: string;
    }

    /**
     * SponsorshipLevel defines the type that encapsulates the different sponsorship levels
     */
    export interface SponsorshipLevel number

    /**
     * UpdateSponsorContactParams defines the inputs used by the UpdateSponsorContactParams API method
     */
    export interface UpdateSponsorContactParams {
        SponsorContactInformation: SponsorContactInformation;
    }

    /**
     * UpdateSponsorContactResponse defines the output returned by the UpdateSponsorContactResponse API method
     */
    export interface UpdateSponsorContactResponse {
    }

    /**
     * Venue defines a venue that hosts a conference, such as DisneyWorld
     */
    export interface Venue {
        ID: number;
        Name: string;
        Description: string;
        Address: string;
        Directions: string;
        GoogleMapsURL: string;
        Capacity: number;
    }

    export class ServiceClient {
        private baseClient: BaseClient;

        constructor(baseClient: BaseClient) {
            this.baseClient = baseClient
        }

        /**
         * GetAll retrieves all conferences and events
         */
        public GetAll(params: GetAllParams): Promise<GetAllResponse> {
            return this.baseClient.do<GetAllResponse>("conferences.GetAll", params);
        }

        /**
         * GetConferenceSlots retrieves all event slots for a specific event id
         */
        public GetConferenceSlots(params: GetConferenceSlotsParams): Promise<GetConferenceSlotsResponse> {
            return this.baseClient.do<GetConferenceSlotsResponse>("conferences.GetConferenceSlots", params);
        }

        /**
         * GetConferenceSponsors retrieves the sponsors for a specific conference
         */
        public GetConferenceSponsors(params: GetConferenceSponsorsParams): Promise<GetConferenceSponsorsResponse> {
            return this.baseClient.do<GetConferenceSponsorsResponse>("conferences.GetConferenceSponsors", params);
        }

        /**
         * GetCurrentByEvent retrieves the current conference and event information for a specific event
         */
        public GetCurrentByEvent(params: GetCurrentByEventParams): Promise<GetCurrentByEventResponse> {
            return this.baseClient.do<GetCurrentByEventResponse>("conferences.GetCurrentByEvent", params);
        }

        /**
         * UpdateSponsorContact retrieves all conferences and events
         */
        public UpdateSponsorContact(params: UpdateSponsorContactParams): Promise<UpdateSponsorContactResponse> {
            return this.baseClient.do<UpdateSponsorContactResponse>("conferences.UpdateSponsorContact", params);
        }
    }
}

class BaseClient {
	baseURL: string;
	headers: {[key: string]: string};

    constructor(environment: string, token?: string) {
		this.headers = {"Content-Type": "application/json"}
		if (token !== undefined) {
			this.headers["Authorization"] = "Bearer " + token
		}
		if (environment === "dev") {
			this.baseURL = "http://localhost:4060/"
		} else {
			this.baseURL = `https://showrunner-46b2.encoreapi.com/${environment}/`
		}
    }

    public async do<T>(endpoint: string, req?: any): Promise<T> {
        let response = await fetch(this.baseURL + endpoint, {
            method: "POST",
            headers: this.headers,
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
            headers: this.headers,
            body: JSON.stringify(req || {})
        })
        if (!response.ok) {
            let body = await response.text()
            throw new Error("request failed: " + body)
        }
        await response.text()
    }
}
