package cfg

const defaultYaml = `includeFilename: true
set:
- messageType: request
  columns:
  - name: created_at
    xpath: //TransportationOrderRequest/MessageHeader/CreationDateTime
  - name: external_number
    xpath: //TransportationOrderRequest/MessageHeader/RecipientParty/InternalID[/@schemeID='FO NUMBER']
- messageType: order
  columns:
  - name: created_at
    xpath: //TransportationOrderQuotationCreateRequest/MessageHeader/CreationDateTime
  - name: external_number
    xpath: //TransportationOrderQuotationCreateRequest/MessageHeader/RecipientParty/InternalID[/@schemeID='FO NUMBER']
- messageType: cancel_request
  columns:
  - name: created_at
    xpath: //TransportationOrderCancellationRequest/MessageHeader/CreationDateTime
  - name: external_number
    xpath: //TransportationOrderCancellationRequest/MessageHeader/RecipientParty/InternalID[/@schemeID='FO NUMBER']
- messageType: cancel_order
  columns:
  - name: created_at
    xpath: //TransportationOrderQuotationCancellationRequest/MessageHeader/CreationDateTime
  - name: external_number
    xpath: //TransportationOrderQuotationCancellationRequest/MessageHeader/RecipientParty/InternalID[/@schemeID='FO NUMBER']
`
