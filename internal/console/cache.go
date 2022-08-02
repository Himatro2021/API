package console

import (
	"context"
	"time"

	"github.com/Himatro2021/API/internal/config"
	"github.com/Himatro2021/API/internal/db"
	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:     "cache",
	Short:   "interact with redis server",
	Long:    "Used to interact with redis cache server",
	Example: "cache get cacheKey",
	Run:     cache,
}

func init() {
	RootCmd.AddCommand(cacheCmd)
}

func cache(_ *cobra.Command, args []string) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         config.RedisAddr(),
		Password:     config.RedisPassword(),
		DB:           config.RedisCacheDB(),
		DialTimeout:  config.RedisTimeout(),
		MinIdleConns: config.RedisMinIdleConn(),
		MaxIdleConns: config.RedisMaxIdleConn(),
	})
	client := db.NewCacher(redisClient)

	if args[0] == "get" {
		r, err := client.Get(context.TODO(), "test")
		if err != nil {
			logrus.Error(err)
			return
		}

		logrus.Info(r)
		return
	}

	if args[0] == "set" {
		if err := client.Set(context.TODO(), "test", `{"json":"test"}`, time.Hour*1); err != nil {
			logrus.Error(err)
		}
		return
	}
}
