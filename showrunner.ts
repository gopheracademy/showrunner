export default class Client {
    conferences: conferences.ServiceClient

    constructor(environment: string = "production", token?: string) {
        const base = new BaseClient(environment, token)
        this.conferences = new conferences.ServiceClient(base)
    }
}

export namespace conferences {
    /**
     * AddPaperParams defines the inputs used by the AddPaper API method
     */
    export interface AddPaperParams {
        Paper: Paper
    }

    /**
     * AddPaperResponse defines the output returned by the AddPaper API method
     */
    export interface AddPaperResponse {
        PaperID: number
    }

    /**
     * Conference is an instance like GopherCon 2020
     */
    export interface Conference {
        ID: number
        Name: string
        Slug: string
        StartDate: string
        EndDate: string
        Venue: Venue
    }

    /**
     * ConferenceSlot holds information for any sellable/giftable slot we have in the event for
     * a Talk or any other activity that requires admission.
     * store: "interface"
     */
    export interface ConferenceSlot {
        ID: number
        Name: string
        Description: string
        Cost: number
        Capacity: number
        StartDate: string
        EndDate: string
        /**
         * DependsOn means that these two Slots need to be acquired together, user must either buy
         * both Slots or pre-own one of the one it depends on.
         */
        DependsOn: number

        /**
         * PurchaseableFrom indicates when this item is on sale, for instance early bird tickets are the first
         * ones to go on sale.
         */
        PurchaseableFrom: string

        /**
         * PuchaseableUntil indicates when this item stops being on sale, for instance early bird tickets can
         * no loger be purchased N months before event.
         */
        PurchaseableUntil: string

        /**
         * AvailableToPublic indicates is this is something that will appear on the tickets purchase page (ie, we can
         * issue sponsor tickets and those cannot be bought individually)
         */
        AvailableToPublic: boolean

        Location: Location
        ConferenceID: number
    }

    /**
     * ContactRole defines the type that encapsulates the different contact roles
     */
    export type ContactRole = number

    /**
     * CreateJobParams defines the inputs used by the CreateJob API method
     */
    export interface CreateJobParams {
        Job: Job
    }

    /**
     * CreateJobResponse defines the output returned by the CreateJob API method
     */
    export interface CreateJobResponse {
        Job: Job
    }

    /**
     * DeleteJobParams defines the input used by
     * the DeleteJob API method
     */
    export interface DeleteJobParams {
        JobID: number
    }

    /**
     * Event is a brand like GopherCon
     */
    export interface Event {
        ID: number
        Name: string
        Slug: string
        Conferences: Conference[]
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
        Events: Event[]
    }

    /**
     * GetConferenceSlotsParams defines the inputs used by the GetConferenceSlots API method
     */
    export interface GetConferenceSlotsParams {
        ConferenceID: number
    }

    /**
     * GetConferenceSlotsResponse defines the output returned by the GetConferenceSlots API method
     */
    export interface GetConferenceSlotsResponse {
        ConferenceSlots: ConferenceSlot[]
    }

    /**
     * GetConferenceSponsorsParams defines the inputs used by the GetConferenceSponsors API method
     */
    export interface GetConferenceSponsorsParams {
        ConferenceID: number
    }

    /**
     * GetConferenceSponsorsResponse defines the output returned by the GetConferenceSponsors API method
     */
    export interface GetConferenceSponsorsResponse {
        Sponsors: Sponsor[]
    }

    /**
     * GetCurrentByEventParams defines the inputs used by the GetCurrentByEvent API method
     */
    export interface GetCurrentByEventParams {
        EventID: number
    }

    /**
     * GetCurrentByEventResponse defines the output returned by the GetCurrentByEvent API method
     */
    export interface GetCurrentByEventResponse {
        Event: Event
    }

    /**
     * GetJobParams defines the inputs used by the GetJob API method
     */
    export interface GetJobParams {
        JobID: number
    }

    /**
     * GetJobResponse defines the output returned by the GetJob API method
     */
    export interface GetJobResponse {
        Job: Job
    }

    /**
     * GetPaperParams defines the inputs used by the GetPaper API method
     */
    export interface GetPaperParams {
        PaperID: number
    }

    /**
     * GetPaperResponse defines the output returned by the GetPaper API method
     */
    export interface GetPaperResponse {
        Paper: Paper
    }

    /**
     * Job represents the necessary information for a Job
     */
    export interface Job {
        ID: number
        CompanyName: string
        Title: string
        Description: string
        Link: string
        Discord: string
        Rank: number
        Approved: boolean
    }

    /**
     * ListApprovedJobsResponse defines the output returned
     * by the ListApprovedJobs API method
     */
    export interface ListApprovedJobsResponse {
        Jobs: Job[]
    }

    /**
     * ListJobsResponse defines the output returned
     * by the ListJobs API method
     */
    export interface ListJobsResponse {
        Jobs: Job[]
    }

    /**
     * ListPapersParams defines the inputs used by the ListPapers API method
     */
    export interface ListPapersParams {
        ConferenceID: number
    }

    /**
     * ListPapersResponse defines the output returned by the ListPapers API method
     */
    export interface ListPapersResponse {
        Papers: Paper[]
    }

    /**
     * Location defines a location for a venue, such as a room or event space
     */
    export interface Location {
        ID: number
        Name: string
        Description: string
        Address: string
        Directions: string
        GoogleMapsURL: string
        Capacity: number
        VenueID: number
    }

    /**
     * Paper holds information about a paper submitted to a conference
     */
    export interface Paper {
        ID: number
        UserID: number
        ConferenceID: number
        Title: string
        ElevatorPitch: string
        Description: string
        Notes: string
    }

    /**
     * Sponsor defines a conference sponsor, such as Google
     */
    export interface Sponsor {
        ID: number
        Name: string
        Address: string
        Website: string
        SponsorshipLevel: SponsorshipLevel
        Contacts: SponsorContactInformation[]
        ConferenceID: number
    }

    /**
     * SponsorContactInformation defines a contact
     * and their information for a sponsor
     */
    export interface SponsorContactInformation {
        ID: number
        Name: string
        Role: ContactRole
        Email: string
        Phone: string
    }

    /**
     * SponsorshipLevel defines the type that encapsulates the different sponsorship levels
     */
    export type SponsorshipLevel = number

    /**
     * UpdateApproveJobParams defines the inputs used by the
     * UpdateApproveJob API method
     */
    export interface UpdateApproveJobParams {
        JobID: number
        ApprovedStatus: boolean
    }

    /**
     * UpdateApproveJobResponse defines the output the returned
     * by the ApproveJob API method
     */
    export interface UpdateApproveJobResponse {
        Job: Job
    }

    /**
     * UpdateJobParams defines the input used by
     * the UpdateJob API method
     */
    export interface UpdateJobParams {
        Job: Job
    }

    /**
     * UpdateJobResponse defines the output returned
     * by the UpdateJob API method
     */
    export interface UpdateJobResponse {
        Job: Job
    }

    /**
     * UpdatePaperParams defines the inputs used by the GetPaper API method
     */
    export interface UpdatePaperParams {
        Paper: Paper
    }

    /**
     * UpdatePaperResponse defines the output received by the UpdatePaper API method
     */
    export interface UpdatePaperResponse {
        Paper: Paper
    }

    /**
     * UpdateSponsorContactParams defines the inputs used by the UpdateSponsorContactParams API method
     */
    export interface UpdateSponsorContactParams {
        SponsorContactInformation: SponsorContactInformation
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
        ID: number
        Name: string
        Description: string
        Address: string
        Directions: string
        GoogleMapsURL: string
        Capacity: number
    }

    export class ServiceClient {
        private baseClient: BaseClient

        constructor(baseClient: BaseClient) {
            this.baseClient = baseClient
        }

        /**
         * AddPaper inserts a paper into the paper_submissions table
         */
        public AddPaper(params: AddPaperParams): Promise<AddPaperResponse> {
            return this.baseClient.do<AddPaperResponse>("conferences.AddPaper", params)
        }

        /**
         * CreateJob inserts a job into the job_board table
         */
        public CreateJob(params: CreateJobParams): Promise<CreateJobResponse> {
            return this.baseClient.do<CreateJobResponse>("conferences.CreateJob", params)
        }

        /**
         * DeleteJob deletes a job by id from the
         * job_board table
         */
        public DeleteJob(params: DeleteJobParams): Promise<void> {
            return this.baseClient.doVoid("conferences.DeleteJob", params)
        }

        /**
         * GetAll retrieves all conferences and events
         */
        public GetAll(params: GetAllParams): Promise<GetAllResponse> {
            return this.baseClient.do<GetAllResponse>("conferences.GetAll", params)
        }

        /**
         * GetConferenceSlots retrieves all event slots for a specific event id
         */
        public GetConferenceSlots(params: GetConferenceSlotsParams): Promise<GetConferenceSlotsResponse> {
            return this.baseClient.do<GetConferenceSlotsResponse>("conferences.GetConferenceSlots", params)
        }

        /**
         * GetConferenceSponsors retrieves the sponsors for a specific conference
         */
        public GetConferenceSponsors(params: GetConferenceSponsorsParams): Promise<GetConferenceSponsorsResponse> {
            return this.baseClient.do<GetConferenceSponsorsResponse>("conferences.GetConferenceSponsors", params)
        }

        /**
         * GetCurrentByEvent retrieves the current conference and event information for a specific event
         */
        public GetCurrentByEvent(params: GetCurrentByEventParams): Promise<GetCurrentByEventResponse> {
            return this.baseClient.do<GetCurrentByEventResponse>("conferences.GetCurrentByEvent", params)
        }

        /**
         * GetJob retrieves a job posting by JobID
         */
        public GetJob(params: GetJobParams): Promise<GetJobResponse> {
            return this.baseClient.do<GetJobResponse>("conferences.GetJob", params)
        }

        /**
         * GetPaper retrieves information for a specific paper id
         */
        public GetPaper(params: GetPaperParams): Promise<GetPaperResponse> {
            return this.baseClient.do<GetPaperResponse>("conferences.GetPaper", params)
        }

        /**
         * ListApprovedJobs retrieves all jobs (approved or not) from
         * the job_board table
         */
        public ListApprovedJobs(): Promise<ListApprovedJobsResponse> {
            return this.baseClient.do<ListApprovedJobsResponse>("conferences.ListApprovedJobs")
        }

        /**
         * ListJobs retrieves all jobs (approved or not) from
         * the job_board table
         */
        public ListJobs(): Promise<ListJobsResponse> {
            return this.baseClient.do<ListJobsResponse>("conferences.ListJobs")
        }

        /**
         * ListPapers retrieves all the papers submitted for a specific conference
         */
        public ListPapers(params: ListPapersParams): Promise<ListPapersResponse> {
            return this.baseClient.do<ListPapersResponse>("conferences.ListPapers", params)
        }

        /**
         * UpdateApproveJob sets the approval of a job to true
         * or false depending on input
         */
        public UpdateApproveJob(params: UpdateApproveJobParams): Promise<UpdateApproveJobResponse> {
            return this.baseClient.do<UpdateApproveJobResponse>("conferences.UpdateApproveJob", params)
        }

        /**
         * UpdateJob updates a job entry based on id in the
         * job_board table
         */
        public UpdateJob(params: UpdateJobParams): Promise<UpdateJobResponse> {
            return this.baseClient.do<UpdateJobResponse>("conferences.UpdateJob", params)
        }

        /**
         * UpdatePaper updates a paper submission for a specific paper id
         */
        public UpdatePaper(params: UpdatePaperParams): Promise<UpdatePaperResponse> {
            return this.baseClient.do<UpdatePaperResponse>("conferences.UpdatePaper", params)
        }

        /**
         * UpdateSponsorContact retrieves all conferences and events
         */
        public UpdateSponsorContact(params: UpdateSponsorContactParams): Promise<UpdateSponsorContactResponse> {
            return this.baseClient.do<UpdateSponsorContactResponse>("conferences.UpdateSponsorContact", params)
        }
    }
}

class BaseClient {
    baseURL: string
    headers: {[key: string]: string}

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
