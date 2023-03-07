package main

import (
	"flag"
	"os"

	"github.com/opensourceways/kafka-lib/kafka"
	"github.com/opensourceways/kafka-lib/mq"
	"github.com/opensourceways/robot-gitee-lib/client"
	"github.com/opensourceways/robot-gitee-lib/framework"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/opensourceways/server-common-lib/secret"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitee-software-package/config"
	"github.com/opensourceways/robot-gitee-software-package/event"
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

	cfg, err := config.LoadConfig(o.service.ConfigFile)
	if err != nil {
		logrus.WithError(err).Fatal("get config failed")
	}

	secretAgent := new(secret.Agent)
	if err = secretAgent.Start([]string{o.gitee.TokenPath}); err != nil {
		logrus.WithError(err).Fatal("Error starting secret agent.")
	}

	defer secretAgent.Stop()

	if err = connectKafka(cfg.Event.KafkaAddress); err != nil {
		logrus.WithError(err).Fatal("init kafka failed")
	}

	c := client.NewClient(secretAgent.GetTokenGenerator(o.gitee.TokenPath))

	e := event.NewEvent(&cfg.Event, c, botName)
	subscribers, err := e.Subscribe()
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
