package config

type OrderConfig struct {
	/// The max. count of notifications that can be qued up, before the notification will get thrown away
	MaxNotificationQueue int
	/// the highest idententifier that can be generated
	MaximalIdentifiers int
}
