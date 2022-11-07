package filterstore

import (
	"strings"

	"github.com/Golang-Tools/optparams"
	"github.com/Golang-Tools/redishelper/v2/redisproxy"
	"github.com/go-redis/redis/v8"
)

type RedisFilterStore struct {
	conn *redisproxy.RedisProxy
	Opts Options
}

func NewRedisFilterStore() *RedisFilterStore {
	s := new(RedisFilterStore)
	s.conn = redisproxy.New()
	s.Opts = DefaultOpts
	return s
}
func (s *RedisFilterStore) Init(opts ...optparams.Option[Options]) error {
	optparams.GetOption(&s.Opts, opts...)
	resinitparams := []optparams.Option[redisproxy.Options]{}
	address := strings.Split(s.Opts.RedisURL, ",")
	if len(address) > 1 {
		copt := redis.ClusterOptions{
			Addrs: address,
		}
		switch s.Opts.RedisRouteMod {
		case "bylatency":
			{
				copt.RouteByLatency = true
			}
		case "randomly":
			{
				copt.RouteRandomly = true
			}
		}
		resinitparams = append(resinitparams, redisproxy.WithClusterOptions(&copt))
	} else {
		resinitparams = append(resinitparams, redisproxy.WithURL(s.Opts.RedisURL))
	}
	if s.Opts.ConnTimeout > 0 {
		resinitparams = append(resinitparams, redisproxy.WithQueryTimeoutMS(s.Opts.ConnTimeout))
	}
	err := s.conn.Init(resinitparams...)
	if err != nil {
		return err
	}
	return nil
}
