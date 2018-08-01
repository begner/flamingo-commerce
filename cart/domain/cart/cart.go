package cart

import (
	"fmt"
	"log"

	"time"

	"math"

	"github.com/pkg/errors"
)

type (
	//CartProvider should be used to create the cart Value objects
	CartProvider func() *Cart

	// Cart Value Object (immutable data - because the cartservice is responsible to return a cart).
	Cart struct {
		//ID is the main identifier of the cart
		ID string
		//EntityID is a second identifier that may be used by some backends
		EntityID string
		//Cartitems - list of cartitems
		Cartitems []Item
		//CartTotals - the cart totals (contain summary costs and discounts etc)
		CartTotals CartTotals
		//BillingAdress - the main billing address (relevant for all payments/invoices)
		BillingAdress Address

		//Purchaser - additional infos for the legal contact person in this order
		Purchaser Person

		//DeliveryInfos - list of desired Deliverys (or Shippments) involved in this cart - referenced from the items
		DeliveryInfos []DeliveryInfo
		//AdditionalData   can be used for Custom attributes
		AdditionalData map[string]string

		//BelongsToAuthenticatedUser - false = Guest Cart true = cart from the authenticated user
		BelongsToAuthenticatedUser bool
		AuthenticatedUserId        string

		AppliedCouponCodes []CouponCode
	}

	CouponCode struct {
		Code string
	}

	Person struct {
		Address         *Address
		PersonalDetails PersonalDetails
		//ExistingCustomerData if the current purchaser is an existing customer - this contains infos about existing customer
		ExistingCustomerData *ExistingCustomerData
	}

	ExistingCustomerData struct {
		//ID of the customer
		ID string
	}

	PersonalDetails struct {
		DateOfBirth     string
		PassportCountry string
		PassportNumber  string
		Nationality     string
	}

	//DeliveryInfo - represents the Delivery
	DeliveryInfo struct {
		ID               string
		Method           string
		Carrier          string
		DeliveryLocation DeliveryLocation
		ShippingItem     ShippingItem
		DesiredTime      time.Time
		AdditionalData   map[string]string
		RelatedFlight    *FlightData
	}

	// TODO: FlightData and RelatedFlight in the DeliveryInfo struct should not be forced on Flamingo Users here.
	// Should move to OM3 somehow
	FlightData struct {
		ScheduledDateTime  time.Time
		Direction          string
		FlightNumber       string
		AirportName        string
		DestinationCountry string
	}

	DeliveryLocation struct {
		Type string
		//Address - only set for type adress
		Address *Address
		//Code - optional idendifier of this location/destination - is used in special destination Types

		Code string
	}

	CartTotals struct {
		Totalitems        []Totalitem
		TotalShippingItem ShippingItem
		//Final sum that need to be payed: GrandTotal = SubTotal + TaxAmount - DiscountAmount + SOME of Totalitems = (Sum of Items RowTotalWithDiscountInclTax) + SOME of Totalitems
		GrandTotal float64
		//SubTotal = SUM of Item RowTotal
		SubTotal float64
		//SubTotalInclTax = SUM of Item RowTotalInclTax
		SubTotalInclTax float64
		//SubTotalWithDiscounts = SubTotal - Sum of Item ItemRelatedDiscountAmount
		SubTotalWithDiscounts float64
		//SubTotalWithDiscountsAndTax= Sum of RowTotalWithItemRelatedDiscountInclTax
		SubTotalWithDiscountsAndTax float64

		//TotalDiscountAmount = SUM of Item TotalDiscountAmount
		TotalDiscountAmount float64
		//TotalNonItemRelatedDiscountAmount= SUM of Item NonItemRelatedDiscountAmount
		TotalNonItemRelatedDiscountAmount float64

		//DEPRICATED
		//DiscountAmount float64

		//TaxAmount = Sum of Item TaxAmount
		TaxAmount float64
		//CurrencyCode of the Total positions
		CurrencyCode string
	}

	// Item for Cart
	Item struct {
		ID              string
		MarketplaceCode string
		//VariantMarketPlaceCode is used for Configurable products
		VariantMarketPlaceCode string
		ProductName            string

		// Source Id of where the items should be initial picked - This is set by the SourcingLogic
		SourceId string

		Qty int

		//DEPRICATED
		//Price float64
		// DEPRICATED
		//DiscountAmount float64
		// DEPRICATED
		//PriceInclTax float64

		DeliveryInfoReference *DeliveryInfo
		CurrencyCode          string

		//OriginalDeliveryIntent can be "delivery" for homedelivery or "pickup_locationcode" or "collectionpoint_locationcode"
		OriginalDeliveryIntent *DeliveryIntent

		AdditionalData map[string]string
		//brutto for single item
		SinglePrice float64
		//netto for single item
		SinglePriceInclTax float64
		//RowTotal = SinglePrice * Qty
		RowTotal float64
		//TaxAmount=Qty * (SinglePriceInclTax-SinglePrice)
		TaxAmount float64
		//RowTotalInclTax= RowTotal + TaxAmount
		RowTotalInclTax float64
		//AppliedDiscounts contains the details about the discounts applied to this item - they can be "itemrelated" or not
		AppliedDiscounts []ItemDiscount
		// TotalDiscountAmount = Sum of AppliedDiscounts = ItemRelatedDiscountAmount +NonItemRelatedDiscountAmount
		TotalDiscountAmount float64
		// ItemRelatedDiscountAmount = Sum of AppliedDiscounts where IsItemRelated = True
		ItemRelatedDiscountAmount float64
		//NonItemRelatedDiscountAmount = Sum of AppliedDiscounts where IsItemRelated = false
		NonItemRelatedDiscountAmount float64
		//RowTotalWithItemRelatedDiscountInclTax=RowTotal-ItemRelatedDiscountAmount
		RowTotalWithItemRelatedDiscount float64
		//RowTotalWithItemRelatedDiscountInclTax=RowTotalInclTax-ItemRelatedDiscountAmount
		RowTotalWithItemRelatedDiscountInclTax float64
		//This is the price the customer finaly need to pay for this item:  RowTotalWithDiscountInclTax = RowTotalInclTax-TotalDiscountAmount
		RowTotalWithDiscountInclTax float64
	}

	// DiscountItem
	ItemDiscount struct {
		Code  string
		Title string
		Price float64
		//IsItemRelated is a flag indicating if the discount should be displayed in the item or if it the result of a cart discount
		IsItemRelated bool
	}

	// Totalitem for totals
	Totalitem struct {
		Code  string
		Title string
		Price float64
		Type  string
	}

	// ShippingItem
	ShippingItem struct {
		Title string
		Price float64

		TaxAmount      float64
		DiscountAmount float64

		CurrencyCode string
	}
)

const (
	DELIVERY_METHOD_PICKUP      = "pickup"
	DELIVERY_METHOD_DELIVERY    = "delivery"
	DELIVERY_METHOD_UNSPECIFIED = "unspecified"

	DELIVERYLOCATION_TYPE_COLLECTIONPOINT = "collection-point"
	DELIVERYLOCATION_TYPE_STORE           = "store"
	DELIVERYLOCATION_TYPE_ADDRESS         = "address"
	DELIVERYLOCATION_TYPE_FREIGHTSTATION  = "freight-station"

	TOTALS_TYPE_DISCOUNT      = "totals_type_discount"
	TOTALS_TYPE_VOUCHER       = "totals_type_voucher"
	TOTALS_TYPE_TAX           = "totals_type_tax"
	TOTALS_TYPE_LOYALTYPOINTS = "totals_loyaltypoints"
	TOTALS_TYPE_SHIPPING      = "totals_type_shipping"

	FLIGHT_DATE_FORMAT = time.RFC3339
)

// GetByLineNr gets an item - starting with 1
func (Cart Cart) HasDeliveryInfos() bool {
	if len(Cart.DeliveryInfos) > 0 {
		return true
	}
	return false
}

// GetByLineNr gets an item - starting with 1
func (Cart Cart) GetByLineNr(lineNr int) (*Item, error) {
	var item Item
	if len(Cart.Cartitems) >= lineNr && lineNr > 0 {
		return &Cart.Cartitems[lineNr-1], nil
	} else {
		return &item, errors.New("Line in cart not existend")
	}
}

// GetMainShippingEMail
func (Cart Cart) GetMainShippingEMail() string {
	for _, info := range Cart.DeliveryInfos {
		if info.DeliveryLocation.Address != nil {
			if info.DeliveryLocation.Address.Email != "" {
				return info.DeliveryLocation.Address.Email
			}
		}
	}
	return ""
}

// GetByItemId gets an item by its id
func (Cart Cart) GetByItemId(itemId string) (*Item, error) {
	for _, currentItem := range Cart.Cartitems {
		log.Println("Cart GetByItemId:" + currentItem.ID)
		if currentItem.ID == itemId {
			return &currentItem, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("itemId %v in cart not existend", itemId))
}

// HasItem checks if a cartitem for that sku exists and returns lineNr if found
func (cart Cart) HasItem(marketplaceCode string, variantMarketplaceCode string) (bool, int) {
	for lineNr, item := range cart.Cartitems {
		if item.MarketplaceCode == marketplaceCode && item.VariantMarketPlaceCode == variantMarketplaceCode {
			return true, lineNr + 1
		}
	}
	return false, 0
}

func inStruct(value string, list []string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

// ItemCount - returns amount of Cartitems
func (Cart Cart) ItemCount() int {
	count := 0
	for _, item := range Cart.Cartitems {
		count += item.Qty
	}

	return count
}

// GetItemIds - returns amount of Cartitems
func (Cart Cart) GetItemIds() []string {
	ids := []string{}
	for _, item := range Cart.Cartitems {
		ids = append(ids, item.ID)
	}
	return ids
}

// GetCartItemsByOriginalIntend - returns the cart Items grouped by the original DeliveryIntent
func (Cart Cart) GetCartItemsByOriginalDeliveryIntent() map[string][]Item {
	result := make(map[string][]Item)
	for _, item := range Cart.Cartitems {
		result[item.OriginalDeliveryIntent.String()] = append(result[item.OriginalDeliveryIntent.String()], item)
	}
	return result
}

// HasItemWithIntent - returns if the cart has an item with the delivery intent
func (Cart Cart) HasItemWithIntent(intent string) bool {
	for _, item := range Cart.Cartitems {
		if item.OriginalDeliveryIntent.String() == intent {
			return true
		}
	}
	return false
}

func (Cart Cart) HasItemWithDifferentIntent(intent string) bool {
	for _, item := range Cart.Cartitems {
		if item.OriginalDeliveryIntent.String() != intent {
			return true
		}
	}
	return false
}

// check if it is a mixed cart with different delivery intents
func (Cart Cart) HasMixedCart() bool {
	// if there is only one or less items in the cart, it can not be a mixed cart
	if Cart.ItemCount() <= 1 {
		return false
	}

	// get intent from first item
	firstItem := Cart.Cartitems[0]
	firstDeliveryIntent := firstItem.OriginalDeliveryIntent.String()

	return Cart.HasItemWithDifferentIntent(firstDeliveryIntent)
}

func (Cart Cart) GetVoucherSavings() float64 {
	totalSavings := 0.0
	for _, item := range Cart.CartTotals.Totalitems {
		if item.Type == TOTALS_TYPE_VOUCHER {
			totalSavings = totalSavings + math.Abs(item.Price)
		}
	}

	if totalSavings < 0 {
		return 0.0
	}

	return totalSavings
}

func (Cart Cart) GetSavings() float64 {
	totalSavings := 0.0
	for _, item := range Cart.CartTotals.Totalitems {
		if item.Type == TOTALS_TYPE_DISCOUNT {
			totalSavings = totalSavings + math.Abs(item.Price)
		}
	}

	if totalSavings < 0 {
		return 0.0
	}

	return totalSavings
}

func (Cart Cart) HasAppliedCouponCode() bool {
	return len(Cart.AppliedCouponCodes) > 0
}

func (ct CartTotals) GetTotalItemsByType(typeCode string) []Totalitem {
	var totalitems []Totalitem
	for _, item := range ct.Totalitems {
		if item.Type == typeCode {
			totalitems = append(totalitems, item)
		}
	}
	return totalitems
}

func (item Item) GetSavingsByItem() float64 {
	totalSavings := 0.0
	for _, discount := range item.AppliedDiscounts {
		totalSavings = totalSavings + math.Abs(discount.Price)
	}

	if totalSavings < 0 {
		return 0.0
	}

	return totalSavings
}

func (d DeliveryInfo) HasRelatedFlight() bool {
	return d.RelatedFlight != nil
}

func (fd *FlightData) GetScheduledDate() string {
	return fd.ScheduledDateTime.Format("2006-01-02")
}

func (fd *FlightData) GetScheduledDateTime() string {
	return fd.ScheduledDateTime.Format(time.RFC3339)
}

//GetScheduledDateTime string from ScheduledDateTime - used for display
func (fd *FlightData) ParseScheduledDateTime() time.Time {
	//"scheduledDateTime": "2017-11-25T06:30:00Z",
	// @todo this looks really strange
	timeResult, e := time.Parse(FLIGHT_DATE_FORMAT, fd.ScheduledDateTime.Format(FLIGHT_DATE_FORMAT))
	if e != nil {
		return time.Now()
	}
	return timeResult
}

func (di DeliveryInfo) String() string {
	if di.Method == DELIVERY_METHOD_PICKUP {
		return di.Method + "_" + di.DeliveryLocation.Type + "_" + di.DeliveryLocation.Code
	}
	return di.Method
}
