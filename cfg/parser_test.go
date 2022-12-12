package cfg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func Test(t *testing.T) {
	XMLParser := XMLParser{
		true,
		[]Line{
			{
				MessageType: "request",
				Columns: []Column{
					{
						"created_at",
						"//TransportationOrderRequest/MessageHeader/CreationDateTime",
						true,
					},
					{
						"external_number",
						"//TransportationOrderRequest/MessageHeader/RecipientParty/InternalID[/@schemeID='FO NUMBER']",
						false,
					},
				},
			},
			{
				MessageType: "order",
				Columns: []Column{
					{
						"created_at",
						"//TransportationOrderQuotationCreateRequest/MessageHeader/CreationDateTime",
						false,
					},
					{
						"external_number",
						"//TransportationOrderQuotationCreateRequest/MessageHeader/RecipientParty/InternalID[/@schemeID='FO NUMBER']",
						false,
					},
				},
			},
			{
				MessageType: "cancel_request",
				Columns: []Column{
					{
						"created_at",
						"//TransportationOrderCancellationRequest/MessageHeader/CreationDateTime",
						false,
					},
					{
						"external_number",
						"//TransportationOrderCancellationRequest/MessageHeader/RecipientParty/InternalID[/@schemeID='FO NUMBER']",
						false,
					},
				},
			},
			{
				MessageType: "cancel_order",
				Columns: []Column{
					{
						"created_at",
						"//TransportationOrderQuotationCancellationRequest/MessageHeader/CreationDateTime",
						false,
					},
					{
						"external_number",
						"//TransportationOrderQuotationCancellationRequest/MessageHeader/RecipientParty/InternalID[/@schemeID='FO NUMBER']",
						false,
					},
				},
			},
		},
	}
	bytes, err := yaml.Marshal(&XMLParser)
	assert.NoError(t, err)
	fmt.Println(string(bytes))
}

func TestXMLParser_Load(t *testing.T) {
	type fields struct {
		Set []Line
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			fields: fields{
				Set: []Line{{
					MessageType: "request",
					Columns: []Column{
						{
							"created_at",
							"//TransportationOrderRequest/MessageHeader/CreationDateTime",
							false,
						},
						{
							"external_number",
							"//TransportationOrderRequest/MessageHeader/RecipientParty/InternalID[/@schemeID='FO NUMBER']",
							false,
						},
					},
				},
					{
						MessageType: "order",
						Columns: []Column{
							{
								"created_at",
								"//TransportationOrderQuotationCreateRequest/MessageHeader/CreationDateTime",
								false,
							},
							{
								"external_number",
								"//TransportationOrderQuotationCreateRequest/MessageHeader/RecipientParty/InternalID[/@schemeID='FO NUMBER']",
								false,
							},
						},
					},
					{
						MessageType: "cancel_request",
						Columns: []Column{
							{
								"created_at",
								"//TransportationOrderCancellationRequest/MessageHeader/CreationDateTime",
								false,
							},
							{
								"external_number",
								"//TransportationOrderCancellationRequest/MessageHeader/RecipientParty/InternalID[/@schemeID='FO NUMBER']",
								false,
							},
						},
					},
					{
						MessageType: "cancel_order",
						Columns: []Column{
							{
								"created_at",
								"//TransportationOrderQuotationCancellationRequest/MessageHeader/CreationDateTime",
								false,
							},
							{
								"external_number",
								"//TransportationOrderQuotationCancellationRequest/MessageHeader/RecipientParty/InternalID[/@schemeID='FO NUMBER']",
								false,
							},
						},
					},
				},
			},
			args: args{
				filename: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xp := &XMLParser{
				Set: tt.fields.Set,
			}
			if err := xp.Load(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("XMLParser.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
