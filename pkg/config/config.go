package config

type Config struct {
	Multipliers *MultiplierConfig
}

type MultiplierConfig struct {
	Failed float64
	Hard   float64
	Normal float64
	Easy   float64
}

var DefaultConfig = &Config{
	Multipliers: &MultiplierConfig{
		Failed: 0.0,
		Hard:   1.0,
		Normal: 1.5,
		Easy:   2.0,
	},
}
