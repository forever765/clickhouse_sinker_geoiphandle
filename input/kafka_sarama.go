/*Copyright [2019] housepower

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package input

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"strings"
	"time"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"github.com/Shopify/sarama"
	"github.com/housepower/clickhouse_sinker/config"
	"github.com/housepower/clickhouse_sinker/ipHandle/entity"
	"github.com/housepower/clickhouse_sinker/model"
	"github.com/housepower/clickhouse_sinker/statistics"
	"github.com/housepower/clickhouse_sinker/util"
	"github.com/pkg/errors"
	"github.com/xdg-go/scram"
	"go.uber.org/zap"
)

var _ Inputer = (*KafkaSarama)(nil)

// KafkaSarama implements input.Inputer
type KafkaSarama struct {
	cfg       *config.Config
	taskCfg   *config.TaskConfig
	cg        sarama.ConsumerGroup
	sess      sarama.ConsumerGroupSession
	stopped   chan struct{}
	putFn     func(msg model.InputMessage)
	cleanupFn func()
}

// NewKafkaSarama get instance of kafka reader
func NewKafkaSarama() *KafkaSarama {
	return &KafkaSarama{}
}

type MyConsumerGroupHandler struct {
	k *KafkaSarama //point back to which kafka this handler belongs to
}

func (h MyConsumerGroupHandler) Setup(sess sarama.ConsumerGroupSession) error {
	h.k.sess = sess
	return nil
}

func (h MyConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	begin := time.Now()
	h.k.cleanupFn()
	util.Logger.Info("consumer group cleanup",
		zap.String("task", h.k.taskCfg.Name),
		zap.String("consumer group", h.k.taskCfg.ConsumerGroup),
		zap.Int32("generation id", h.k.sess.GenerationID()),
		zap.Duration("cost", time.Since(begin)))
	return nil
}

func SearchIP(json_raw *jsonvalue.V) *jsonvalue.V {
	// handle ip_src and ip_dst
	objs := []string{"src", "dst"}
	for _, obj := range objs {
		ip, _ := json_raw.Get("ip_" + obj)
		naliRsp := entity.ParseIP(ip.String()).String()

		naliRsp = strings.TrimRight(naliRsp, "]")
		PureResult := strings.Split(naliRsp, "[")

		var PureResultList []string
		// PureResult  --->   [220.166.187.228 四川省资阳市简阳市]
		if len(PureResult) > 1 {
			PureResultList = strings.Fields(PureResult[1])
		}
		LPR := len(PureResultList)
		loc := "Unknown"
		isp := "Unknown"

		if LPR == 0 {
			// if nali return null result, default value is "Unknown"
		} else if LPR == 1 {
			// only have location
			loc = PureResultList[0]
		} else if LPR > 1 {
			// 国外的地名和isp可能有空格
			loc = PureResultList[0]
			isp = strings.Join(PureResultList[1:], "")
		} else {
			util.Logger.Warn(fmt.Sprintf("nali return unknown data: %s, 个数：%v", PureResultList, LPR))
		}
		// Replace ... to 局域网
		if strings.Contains(loc, "同一内部网") || strings.Contains(isp, "同一内部网") {
			loc = "局域网"
			isp = "局域网"
		}
		if strings.Contains(isp, "]") {
			isp = strings.TrimRight(isp, "]")
			// fmt.Println("有漏网之鱼, ] replaced!")
		}
		json_raw.SetString(loc).At("loc_" + obj)
		json_raw.SetString(isp).At("isp_" + obj)
	}
	// FinalResult := json_raw.MustMarshalString()
	return json_raw
}

// Unknown/Unknown to Unknown
// Unknown/XXX to  XXX
func ReduceUnknown(json_raw *jsonvalue.V) *jsonvalue.V {
	class_raw, err := json_raw.Get("class")
	if err != nil {
		return json_raw
	}
	class := class_raw.String()
	if class == "Unknown/Unknown" {
		class = strings.Replace(class, "Unknown/Unknown", "Unknown", -1)
	} else if strings.Contains(class, "/") {
		ClassList := strings.Split(class, "/")
		if ClassList[0] != ClassList[1] {
			class = ClassList[1]
		}
	}
	json_raw.SetString(class).At("class")
	return json_raw
}

func HandleMsg(json []byte) []byte {
	// Unmarshal JSON
	json_raw, err := jsonvalue.UnmarshalString(string(json))
	if err != nil {
		errLog := "JSON 解码失败：" + err.Error() + "\n"
		util.Logger.Error(errLog)
	}
	GeoDoneResult := SearchIP(json_raw)
	ClassDone := ReduceUnknown(GeoDoneResult)
	FinalResult := ClassDone.MustMarshalString()
	return []byte(FinalResult)
}

func (h MyConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// util.Logger.Info(fmt.Sprintf("GeoipHandle status: %v", h.k.taskCfg.GeoipHandle))
	for msg := range claim.Messages() {
		// if need handle geoip
		if h.k.taskCfg.GeoipHandle {
			msg.Value = HandleMsg(msg.Value)
		}
		h.k.putFn(model.InputMessage{
			Topic:     msg.Topic,
			Partition: int(msg.Partition),
			Key:       msg.Key,
			Value:     msg.Value,
			Offset:    msg.Offset,
			Timestamp: &msg.Timestamp,
		})
	}
	return nil
}

// Init Initialise the kafka instance with configuration
func (k *KafkaSarama) Init(cfg *config.Config, taskCfg *config.TaskConfig, putFn func(msg model.InputMessage), cleanupFn func()) (err error) {
	k.cfg = cfg
	k.taskCfg = taskCfg
	kfkCfg := &cfg.Kafka
	k.stopped = make(chan struct{})
	k.putFn = putFn
	k.cleanupFn = cleanupFn
	config := sarama.NewConfig()
	if config.Version, err = sarama.ParseKafkaVersion(kfkCfg.Version); err != nil {
		err = errors.Wrapf(err, "")
		return
	}
	if kfkCfg.TLS.CaCertFiles == "" && kfkCfg.TLS.TrustStoreLocation != "" {
		if kfkCfg.TLS.CaCertFiles, _, err = util.JksToPem(kfkCfg.TLS.TrustStoreLocation, kfkCfg.TLS.TrustStorePassword, false); err != nil {
			return
		}
	}
	if kfkCfg.TLS.ClientKeyFile == "" && kfkCfg.TLS.KeystoreLocation != "" {
		if kfkCfg.TLS.ClientCertFile, kfkCfg.TLS.ClientKeyFile, err = util.JksToPem(kfkCfg.TLS.KeystoreLocation, kfkCfg.TLS.KeystorePassword, false); err != nil {
			return
		}
	}
	if kfkCfg.TLS.Enable {
		config.Net.TLS.Enable = true
		if config.Net.TLS.Config, err = util.NewTLSConfig(kfkCfg.TLS.CaCertFiles, kfkCfg.TLS.ClientCertFile, kfkCfg.TLS.ClientKeyFile, kfkCfg.TLS.EndpIdentAlgo == ""); err != nil {
			return
		}
	}
	// check for authentication
	if kfkCfg.Sasl.Enable {
		config.Net.SASL.Enable = true
		if config.Version.IsAtLeast(sarama.V1_0_0_0) {
			config.Net.SASL.Version = sarama.SASLHandshakeV1
		}
		config.Net.SASL.Mechanism = (sarama.SASLMechanism)(kfkCfg.Sasl.Mechanism)
		switch config.Net.SASL.Mechanism {
		case "SCRAM-SHA-256":
			config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
		case "SCRAM-SHA-512":
			config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
		default:
		}
		config.Net.SASL.User = kfkCfg.Sasl.Username
		config.Net.SASL.Password = kfkCfg.Sasl.Password
		config.Net.SASL.GSSAPI = kfkCfg.Sasl.GSSAPI
	}
	if taskCfg.Earliest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	config.ChannelBufferSize = 1024
	cg, err := sarama.NewConsumerGroup(strings.Split(kfkCfg.Brokers, ","), taskCfg.ConsumerGroup, config)
	if err != nil {
		return err
	}
	k.cg = cg
	return nil
}

// kafka main loop
func (k *KafkaSarama) Run(ctx context.Context) {
	taskCfg := k.taskCfg
LOOP_SARAMA:
	for {
		handler := MyConsumerGroupHandler{k}
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := k.cg.Consume(ctx, []string{taskCfg.Topic}, handler); err != nil {
			if errors.Is(err, context.Canceled) {
				util.Logger.Info("KafkaSarama.Run quit due to context has been canceled", zap.String("task", k.taskCfg.Name))
				break LOOP_SARAMA
			} else if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				util.Logger.Info("KafkaSarama.Run quit due to consumer group has been closed", zap.String("task", k.taskCfg.Name))
				break LOOP_SARAMA
			} else {
				statistics.ConsumeMsgsErrorTotal.WithLabelValues(taskCfg.Name).Inc()
				err = errors.Wrap(err, "")
				util.Logger.Error("sarama.ConsumerGroup.Consume failed", zap.String("task", k.taskCfg.Name), zap.Error(err))
				continue
			}
		}
	}
	k.stopped <- struct{}{}
}

func (k *KafkaSarama) CommitMessages(ctx context.Context, msg *model.InputMessage) error {
	k.sess.MarkOffset(msg.Topic, int32(msg.Partition), msg.Offset+1, "")
	return nil
}

// Stop kafka consumer and close all connections
func (k *KafkaSarama) Stop() error {
	k.cg.Close()
	<-k.stopped
	return nil
}

// Description of this kafka consumer, which topic it reads from
func (k *KafkaSarama) Description() string {
	return "kafka consumer of topic " + k.taskCfg.Topic
}

// Predefined SCRAMClientGeneratorFunc, copied from https://github.com/Shopify/sarama/blob/master/examples/sasl_scram_client/scram_client.go

var SHA256 scram.HashGeneratorFcn = func() hash.Hash { return sha256.New() }
var SHA512 scram.HashGeneratorFcn = func() hash.Hash { return sha512.New() }

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}
