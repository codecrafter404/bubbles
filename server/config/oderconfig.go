package config

type OrderConfig struct {
	/// The max. count of notifications that can be qued up, before the notification will get thrown away
	MaxNotificationQueue int `default:"1000"`
	/// the highest idententifier that can be generated
	MaximalIdentifiers int `default:"100"`

	LoadBalancingStrategy LoadBalancing `default:"round_robin"`

	LoadBalancingRandomTimeout int `default:"100"`
}

type LoadBalancing string

const (
	LoadBalancingRoundRobin = "round_robin"
	LoadBalancingRandom     = "random"
)

func (lb LoadBalancing) String() string {
	return string(lb)
}
