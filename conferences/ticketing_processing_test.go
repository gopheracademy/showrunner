package conferences

import (
	"context"
	"reflect"
	"testing"

	"encore.dev/storage/sqldb"
)

func Test_claimSlots(t *testing.T) {
	type args struct {
		ctx   context.Context
		slots map[*Attendee][]ConferenceSlot
	}

	tx, err := sqldb.Begin(context.TODO())
	if err != nil {
		t.Fatalf("beginning transaction: %v", err)
	}
	att01 := &Attendee{
		Email:       "testmail01@gophercon.com",
		CoCAccepted: true,
	}
	savedAttendee01, err := createAttendee(context.TODO(), tx, att01)
	if err != nil {
		t.Fatalf("creating attendee: %v", err)
	}

	att02 := &Attendee{
		Email:       "testmail02@gophercon.com",
		CoCAccepted: true,
	}
	savedAttendee02, err := createAttendee(context.TODO(), tx, att02)
	if err != nil {
		t.Fatalf("creating attendee: %v", err)
	}

	savedAttendee03, err := readAttendee(context.TODO(), tx, savedAttendee02.Email, savedAttendee02.ID)
	if err != nil {
		t.Fatalf("reading attendee: %v", err)
	}
	if savedAttendee03 == nil {
		t.Fatalf("could not find attendee 03 with email %s and id %d", savedAttendee02.Email, savedAttendee02.ID)
	}

	if err := sqldb.Commit(tx); err != nil {
		t.Fatalf("committing test setup transaction: %v", err)
	}

	// There is an entry for general admision to gophercon 2021 preloaded in the first migration
	cslot, err := readConferenceSlotByID(context.TODO(), nil, 1, false)
	if err != nil {
		t.Fatalf("retrieving conference slot: %v", err)
	}

	tests := []struct {
		name    string
		args    args
		want    map[*Attendee][]SlotClaim
		wantErr map[*Attendee]bool
	}{
		{
			name: "claim admission to general event",
			args: args{
				ctx:   context.TODO(),
				slots: map[*Attendee][]ConferenceSlot{savedAttendee01: []ConferenceSlot{*cslot}},
			},
			want: map[*Attendee][]SlotClaim{savedAttendee01: []SlotClaim{
				{
					ID:             1,
					ConferenceSlot: cslot,
					TicketID:       [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					Redeemed:       false,
				},
			},
			},
			wantErr: map[*Attendee]bool{
				savedAttendee01: false,
			},
		},
		{
			name: "claim more admission to general event",
			args: args{
				ctx: context.TODO(),
				slots: map[*Attendee][]ConferenceSlot{
					savedAttendee02: {*cslot},
					savedAttendee03: {*cslot},
				},
			},
			want: map[*Attendee][]SlotClaim{
				savedAttendee02: {
					{
						ID:             2,
						ConferenceSlot: cslot,
						TicketID:       [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						Redeemed:       false,
					},
				},
				savedAttendee03: {
					{
						ID:             3,
						ConferenceSlot: cslot,
						TicketID:       [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						Redeemed:       false,
					},
				},
			},
			wantErr: map[*Attendee]bool{
				savedAttendee02: false,
				savedAttendee03: false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for attendee, slots := range tt.args.slots {
				t.Logf("processing attendee %s", attendee.Email)
				got, err := claimSlots(tt.args.ctx, attendee, slots)
				if (err != nil) != tt.wantErr[attendee] {
					t.Errorf("claimSlots() error = %v, wantErr %v", err, tt.wantErr[attendee])
					return
				}

				rows, err := sqldb.Query(context.TODO(), "SELECT id, ticket_id FROM slot_claim WHERE attendee_id = $1 ORDER BY id DESC", attendee.ID)
				if err != nil {
					t.Fatalf("retrieving new ticket IDs %v", err)
					return
				}

				for i := len(tt.want[attendee]) - 1; i >= 0; i-- {
					if !rows.Next() {
						t.Logf("missing slot claims, only found %d", i)
						t.FailNow()
						return
					}
					if err = rows.Scan(&tt.want[attendee][i].ID, &tt.want[attendee][i].TicketID); err != nil {
						t.Fatalf("scanning new ticket IDs %v", err)
						return
					}
				}
				rows.Close()
				if !reflect.DeepEqual(got, tt.want[attendee]) {
					t.Errorf("claimSlots() = %v, want %v", got, tt.want[attendee])
				}
			}

		})
	}
}
