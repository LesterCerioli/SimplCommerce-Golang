package entities

type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "New"
	OrderStatusOnHold          OrderStatus = "OnHold"
	OrderStatusPendingPayment  OrderStatus = "PendingPayment"
	OrderStatusPaymentReceived OrderStatus = "PaymentReceived"
	OrderStatusPaymentFailed   OrderStatus = "PaymentFailed"
	OrderStatusInvoiced        OrderStatus = "Invoiced"
	OrderStatusShipping        OrderStatus = "Shipping"
	OrderStatusShipped         OrderStatus = "Shipped"
	OrderStatusComplete        OrderStatus = "Complete"
	OrderStatusCanceled        OrderStatus = "Canceled"
	OrderStatusRefunded        OrderStatus = "Refunded"
	OrderStatusClosed          OrderStatus = "Closed"
)

type PaymentStatus string

const (
	PaymentStatusSucceeded PaymentStatus = "Succeeded"
	PaymentStatusFailed    PaymentStatus = "Failed"
)

type ReviewStatus string

const (
	ReviewStatusPending     ReviewStatus = "Pending"
	ReviewStatusApproved    ReviewStatus = "Approved"
	ReviewStatusNotApproved ReviewStatus = "NotApproved"
)

type MediaType int

const (
	MediaTypeImage MediaType = iota
	MediaTypeVideo
	MediaTypeFile
)

type ProductLinkType int

const (
	ProductLinkSuper     ProductLinkType = 1
	ProductLinkRelated   ProductLinkType = 2
	ProductLinkCrossSell ProductLinkType = 3
	ProductLinkUpSell    ProductLinkType = 4
)

type AddressType int

const (
	AddressTypeShipping AddressType = iota
	AddressTypeBilling
)
