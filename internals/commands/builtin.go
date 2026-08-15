package commands

import (
	"Raven/internals/database"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (r *CommandRegistry) registerDefaults() {
	r.Register("PING", pingHandler)
	r.Register("SET", setHandler)
	r.Register("GET", getHandler)
	r.Register("DEL", delHandler)
	r.Register("EXISTS", existsHandler)
	r.Register("KEYS", keysHandler)
	r.Register("EXPIRE", expireHandler)
	r.Register("TTL", ttlHandler)
	r.Register("INCR", incrHandler)
	r.Register("DECR", decrHandler)
	r.Register("MGET", mgetHandler)
	r.Register("MSET", msetHandler)
}

func pingHandler(args []string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " ") + "\n", nil
	}
	return "PONG\n", nil
}

func setHandler(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'set' command")
	}
	key := args[0]
	val := args[1]
	var ttl uint64 = 0
	if len(args) >= 3 {
		parsedTTL, err := strconv.ParseUint(args[2], 10, 64)
		if err == nil {
			if parsedTTL > 0 {
				ttl = uint64(time.Now().Unix()) + parsedTTL
			}
		}
	}
	database.Data_store.SetValue(key, val, ttl)
	return "OK\n", nil
}

func getHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'get' command")
	}
	key := args[0]
	valObj, exists := database.Data_store.GetValue(key)
	if !exists {
		return "(nil)\n", nil
	}

	// Check TTL expiration
	if valObj.Ttl > 0 && uint64(time.Now().Unix()) > valObj.Ttl {
		database.Data_store.DelValue(key)
		return "(nil)\n", nil
	}

	return fmt.Sprintf("%v\n", valObj.Value), nil
}

func delHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'del' command")
	}
	deletedCount := 0
	for _, key := range args {
		if database.Data_store.DelValue(key) {
			deletedCount++
		}
	}
	return fmt.Sprintf("(integer) %d\n", deletedCount), nil
}

func existsHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'exists' command")
	}
	existsCount := 0
	for _, key := range args {
		valObj, exists := database.Data_store.GetValue(key)
		if exists {
			if valObj.Ttl > 0 && uint64(time.Now().Unix()) > valObj.Ttl {
				database.Data_store.DelValue(key)
			} else {
				existsCount++
			}
		}
	}
	return fmt.Sprintf("(integer) %d\n", existsCount), nil
}

func keysHandler(args []string) (string, error) {
	database.Data_store.GetValue("") // Trigger RLock
	var keysList []string
	now := uint64(time.Now().Unix())

	for k, v := range database.Data_store.Data {
		if v.Ttl > 0 && now > v.Ttl {
			continue
		}
		keysList = append(keysList, k)
	}

	if len(keysList) == 0 {
		return "(empty list or set)\n", nil
	}
	return strings.Join(keysList, "\n") + "\n", nil
}

func expireHandler(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'expire' command")
	}
	key := args[0]
	seconds, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return "", fmt.Errorf("ERR value is not an integer or out of range")
	}

	valObj, exists := database.Data_store.GetValue(key)
	if !exists {
		return "(integer) 0\n", nil
	}

	ttl := uint64(time.Now().Unix()) + seconds
	database.Data_store.SetValue(key, valObj.Value, ttl)
	return "(integer) 1\n", nil
}

func ttlHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'ttl' command")
	}
	key := args[0]
	valObj, exists := database.Data_store.GetValue(key)
	if !exists {
		return "(integer) -2\n", nil // Key does not exist
	}

	if valObj.Ttl == 0 {
		return "(integer) -1\n", nil // Key exists but has no associated expire
	}

	now := uint64(time.Now().Unix())
	if now > valObj.Ttl {
		database.Data_store.DelValue(key)
		return "(integer) -2\n", nil
	}

	remaining := int64(valObj.Ttl - now)
	return fmt.Sprintf("(integer) %d\n", remaining), nil
}

func incrHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'incr' command")
	}
	key := args[0]
	valObj, exists := database.Data_store.GetValue(key)
	var currentVal int64 = 0

	if exists {
		strVal := fmt.Sprintf("%v", valObj.Value)
		parsed, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return "", fmt.Errorf("ERR value is not an integer or out of range")
		}
		currentVal = parsed
	}

	newVal := currentVal + 1
	database.Data_store.SetValue(key, strconv.FormatInt(newVal, 10), 0)
	return fmt.Sprintf("(integer) %d\n", newVal), nil
}

func decrHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'decr' command")
	}
	key := args[0]
	valObj, exists := database.Data_store.GetValue(key)
	var currentVal int64 = 0

	if exists {
		strVal := fmt.Sprintf("%v", valObj.Value)
		parsed, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return "", fmt.Errorf("ERR value is not an integer or out of range")
		}
		currentVal = parsed
	}

	newVal := currentVal - 1
	database.Data_store.SetValue(key, strconv.FormatInt(newVal, 10), 0)
	return fmt.Sprintf("(integer) %d\n", newVal), nil
}

func mgetHandler(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'mget' command")
	}
	var results []string
	now := uint64(time.Now().Unix())

	for _, key := range args {
		valObj, exists := database.Data_store.GetValue(key)
		if !exists || (valObj.Ttl > 0 && now > valObj.Ttl) {
			results = append(results, "(nil)")
		} else {
			results = append(results, fmt.Sprintf("%v", valObj.Value))
		}
	}
	return strings.Join(results, "\n") + "\n", nil
}

func msetHandler(args []string) (string, error) {
	if len(args) < 2 || len(args)%2 != 0 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'mset' command")
	}
	for i := 0; i < len(args); i += 2 {
		key := args[i]
		val := args[i+1]
		database.Data_store.SetValue(key, val, 0)
	}
	return "OK\n", nil
}
