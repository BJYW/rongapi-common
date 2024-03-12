package config

var AliMnsConfig struct {
	Key           string `goblet:"key,WC9nreih4P8Mbwwz"`
	Secret        string `goblet:"secret,HVllgK8DF9pvURVnTm0dtbFIqS8e6l"`
	ChargeQueue   string `goblet:"charge_queue,charge"`
	PolicyQueue   string `goblet:"policy_queue,policy"`
	RechargeQueue string `goblet:"recharge_queue,recharge"`
	OwerId        string `goblet:"owerid,50274035"`
	Location      string `goblet:"location,cn-beijing"`
}
