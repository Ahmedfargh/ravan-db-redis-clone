package commands

func (r *CommandRegistry) registerDefaults() {
	r.Register("PING", &PingCommand{})
	r.Register("SET", &SetCommand{})
	r.Register("GET", &GetCommand{})
	r.Register("DEL", &DelCommand{})
	r.Register("EXISTS", &ExistsCommand{})
	r.Register("KEYS", &KeysCommand{})
	r.Register("EXPIRE", &ExpireCommand{})
	r.Register("TTL", &TtlCommand{})
	r.Register("INCR", &IncrCommand{})
	r.Register("DECR", &DecrCommand{})
	r.Register("MGET", &MgetCommand{})
	r.Register("MSET", &MsetCommand{})
	r.Register("ECHO", &EchoCommand{})
	r.Register("VALUEBYINDEX", &ValueByIndex{})
}
