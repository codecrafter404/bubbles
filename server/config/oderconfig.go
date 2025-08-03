package config

type OrderConfig struct {
	/// The max. count of orders, before the orderchannel will be blocked
	MaximalConcurrentOrders int
	/// the highest idententifier that can be generated
	MaximalIdentifiers int
}
