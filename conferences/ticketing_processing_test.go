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

func Test_payClaims(t *testing.T) {
	type args struct {
		ctx      context.Context
		attendee *Attendee
		claims   []SlotClaim
		payments []FinancialInstrument
	}
	type want struct {
		totalDue  int64
		fullyPaid bool
		fulfilled bool
	}

	att01 := &Attendee{
		Email:       "testmail01@gophercon.com",
		CoCAccepted: true,
	}
	savedAttendee01, err := createAttendee(context.TODO(), nil, att01)
	if err != nil {
		t.Fatalf("creating attendee: %v", err)
	}

	// There is an entry for general admision to gophercon 2021 preloaded in the first migration
	cslot, err := readConferenceSlotByID(context.TODO(), nil, 1, false)
	if err != nil {
		t.Fatalf("retrieving conference slot: %v", err)
	}
	// test01 setup
	claims01, err := claimSlots(context.TODO(), savedAttendee01, []ConferenceSlot{*cslot})
	if err != nil {
		t.Fatalf("claiming conference slot 1 of: %v", err)
	}

	// test02 setup
	claims02, err := claimSlots(context.TODO(), savedAttendee01, []ConferenceSlot{*cslot})
	if err != nil {
		t.Fatalf("claiming conference slot 2 of: %v", err)
	}

	// test03 setup
	claims03, err := claimSlots(context.TODO(), savedAttendee01, []ConferenceSlot{*cslot})
	if err != nil {
		t.Fatalf("claiming conference slot 3 of: %v", err)
	}

	// test04 setup
	claims04, err := claimSlots(context.TODO(), savedAttendee01, []ConferenceSlot{*cslot})
	if err != nil {
		t.Fatalf("claiming conference slot 3 of: %v", err)
	}

	// test05 setup
	claims05, err := claimSlots(context.TODO(), savedAttendee01, []ConferenceSlot{*cslot})
	if err != nil {
		t.Fatalf("claiming conference slot 3 of: %v", err)
	}

	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{name: "full payment",
			args: args{ctx: context.TODO(),
				attendee: savedAttendee01,
				claims:   claims01,
				payments: []FinancialInstrument{&PaymentMethodMoney{
					PaymentRef: "somethingbystripe",
					Amount:     400, // total of initial slot
				}}},
			want: want{totalDue: 400, fullyPaid: true, fulfilled: true},
		},
		{name: "full mixed payment",
			args: args{ctx: context.TODO(),
				attendee: savedAttendee01,
				claims:   claims02,
				payments: []FinancialInstrument{
					&PaymentMethodMoney{
						PaymentRef: "somethingbystripe",
						Amount:     200,
					},
					&PaymentMethodConferenceDiscount{
						Detail: "nice people discount",
						Amount: 200,
					},
				}},
			want: want{totalDue: 400, fullyPaid: true, fulfilled: true},
		},
		{name: "partial credit payment",
			args: args{ctx: context.TODO(),
				attendee: savedAttendee01,
				claims:   claims03,
				payments: []FinancialInstrument{
					&PaymentMethodMoney{
						PaymentRef: "somethingbystripe",
						Amount:     200,
					},
					&PaymentMethodCreditNote{
						Detail: "IOU from sponsor",
						Amount: 200,
					},
				}},
			want: want{totalDue: 400, fullyPaid: false, fulfilled: true},
		},
		{name: "full credit payment",
			args: args{ctx: context.TODO(),
				attendee: savedAttendee01,
				claims:   claims04,
				payments: []FinancialInstrument{
					&PaymentMethodMoney{
						PaymentRef: "somethingbystripe",
						Amount:     200,
					},
					&PaymentMethodMoney{
						PaymentRef: "somethingbystripe01",
						Amount:     200,
					},
					&PaymentMethodCreditNote{
						Detail: "IOU from sponsor",
						Amount: 200,
					},
				}},
			want: want{totalDue: 400, fullyPaid: true, fulfilled: true},
		},
		{name: "partial payment",
			args: args{ctx: context.TODO(),
				attendee: savedAttendee01,
				claims:   claims05,
				payments: []FinancialInstrument{
					&PaymentMethodMoney{
						PaymentRef: "somethingbystripe",
						Amount:     200,
					},
				}},
			want: want{totalDue: 400, fullyPaid: false, fulfilled: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := payClaims(tt.args.ctx, tt.args.attendee, tt.args.claims, tt.args.payments)
			if (err != nil) != tt.wantErr {
				t.Errorf("payClaims() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Paid() != tt.want.fullyPaid {
				t.Errorf("payClaims() payment = %v, want %v", got.Paid(), tt.want.fullyPaid)
			}
			if got.Fulfilled() != tt.want.fulfilled {
				t.Errorf("payClaims() fulfillment = %v, want %v", got.Fulfilled(), tt.want.fulfilled)
			}
			if got.TotalDue() != tt.want.totalDue {
				t.Errorf("payClaims() payment due = %d , want %d", got.TotalDue(), tt.want.totalDue)
			}
		})
	}
}
