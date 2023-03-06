package main

import (
	"errors"
	"flag"
	"os"

	"github.com/opensourceways/kafka-lib/kafka"
	"github.com/opensourceways/kafka-lib/mq"
	"github.com/opensourceways/robot-gitee-lib/client"
	"github.com/opensourceways/robot-gitee-lib/framework"
	"github.com/opensourceways/server-common-lib/config"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/opensourceways/server-common-lib/secret"
	"github.com/sirupsen/logrus"
)

type options struct {
	service liboptions.ServiceOptions
	gitee   liboptions.GiteeOptions
}

func (o *options) Validate() error {
	if err := o.service.Validate(); err != nil {
		return err
	}

	return o.gitee.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.gitee.AddFlags(fs)
	o.service.AddFlags(fs)

	fs.Parse(args)
	return o
}

func main() {
	logrusutil.ComponentInit(botName)

	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.Validate(); err != nil {
		logrus.WithError(err).Fatal("Invalid options")
	}

	configAgent := config.NewConfigAgent(func() config.Config {
		return new(configuration)
	})
	if err := configAgent.Start(o.service.ConfigFile); err != nil {
		logrus.WithError(err).Fatal("Error starting config agent.")
	}

	defer configAgent.Stop()

	secretAgent := new(secret.Agent)
	if err := secretAgent.Start([]string{o.gitee.TokenPath}); err != nil {
		logrus.WithError(err).Fatal("Error starting secret agent.")
	}

	defer secretAgent.Stop()

	cfg, err := getConfig(&configAgent)
	if err != nil {
		logrus.WithError(err).Fatal("get config failed")
	}

	if err = connectKafka(cfg.KafkaAddress); err != nil {
		logrus.WithError(err).Fatal("init kafka failed")
	}

	c := client.NewClient(secretAgent.GetTokenGenerator(o.gitee.TokenPath))

	e := newEvent(cfg, c)
	subscribers, err := e.subscribe()
	if err != nil {
		logrus.WithError(err).Fatal("subscribe failed")
	}
	defer func() {
		for k, v := range subscribers {
			if err := v.Unsubscribe(); err != nil {
				logrus.Errorf("failed to unsubscribe for topic:%s, err:%v", k, err)
			}
		}
	}()

	r := newRobot(c)

	framework.Run(r, o.service)
}

func getConfig(agent *config.ConfigAgent) (*configuration, error) {
	_, cfg := agent.GetConfig()
	c, ok := cfg.(*configuration)
	if !ok {
		return nil, errors.New("can't convert to configuration")
	}

	return c, nil
}

func connectKafka(address string) error {
	err := kafka.Init(
		mq.Addresses(address),
		mq.Log(logrus.WithField("module", "kfk")),
	)
	if err != nil {
		return err
	}

	return kafka.Connect()
}
