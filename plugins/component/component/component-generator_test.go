package component

import (
	"log"
	"strings"
	"testing"

	"github.com/Mattilsynet/map-cli/plugins/component/project"
)

func TestGenerateComponent(t *testing.T) {
	config := NewConfig(
		"/home/solve/git/temp/testComponent",
		"testComponent",
		"github.com/Mattilsynet/test-component",
		[]string{"nats-core", "nats-jetstream", "nats-kv"},
		WithComponentCode(),
		WithWitPackage())
	config.ExportNatsCoreWit = true
	config.ExportNatsCoreRequestReplyWit = true
	config.ImportNatsKvWit = true
	config.ExportNatsJetstreamWit = true
	config.ImportNatsJetstreamWit = true

	err := GenerateApp(config)

	log.Println(err)
}

func TestComponentGeneration(t *testing.T) {
	expected := `
	//go:generate go run github.com/bytecodealliance/wasm-tools-go/cmd/wit-bindgen-go generate -world test-component -out gen ./wit
package main

import (
	"log/slog"
	"go.wasmcloud.dev/component/log/wasilog"	
	"github.com/Mattilsynet/test-component/pkg/nats"
)
var (
	logger *slog.Logger
	conn *nats.Conn
	js *nats.JetstreamContext
	kv *nats.KeyValue
	
)
func init() 
	logger = wasilog.ContextLogger("test-component")
	conn := nats.NewConn()
        conn.RegisterRequestReply(test-componentRequestReplier)
	var err error
	js, err = conn.Jetstream()
	if err != nil {
		logger.Error("error getting jetstreamcontext", "err", err)
		return
	}
	kv, err = js.KeyValue()
	if kv == nil {
		logger.Error("error getting keyvalue", "err", err)
		return
	}
	js.Subscribe(test-componentConsumer, logger)	
	
}

func test-componentConsumer(msg *nats.Msg) {
// this is where you put your application logic when a nats msg arrives at specified stream subject from wadm.yaml
// here you can use nats.Publish or kv.Get, etc depending on what you asked for
}


func test-componentRequestReplier(msg *nats.Msg) msg *nats.Msg {
// here you'll read from incomming msg and take out the msg.Reply and put into your own created msg
// and return it, and the provider will take it from there
}


func test-componentSubscriber(msg *nats.Msg) {
// this is where you put your application logic when a nats msg arrives at specified subject, specified under wadm.yaml nats-core 'subscriptions'
// here you can use nats.Publish or kv.Get, etc depending on what you asked for
}

//main should never be used in a wasm component, everything inside init()
func main() {}
	`
	tmpls := project.Templs
	config := NewConfig(
		"/home/solve/git/temp",
		"test-component",
		"github.com/Mattilsynet/test-component",
		[]string{"nats-core", "nats-jetstream", "nats-kv"},
		WithComponentCode(),
		WithWitPackage())
	config.ExportNatsJetstreamWit = true
	config.ExportNatsCoreWit = true
	config.ExportNatsCoreRequestReplyWit = true
	a, err := ExecuteTmplWithData(config, tmpls["component.go"])
	if err != nil {
		log.Println("err: ", err)
		t.Fail()
	}
	if strings.TrimSpace(a) != strings.TrimSpace(expected) {
		t.Fail()
	}
}

func TestWitGeneration(t *testing.T) {
	expected := `
package Mattilsynet:test-component;

world test-component{
  //wasmcloud-go sdk for ease of logging and other goodies
  include wasmcloud:component-go/imports@0.1.0;
  
  import wasmcloud:messaging/consumer@0.2.0;
  export wasmcloud:messaging/handler@0.2.0;
  
  import mattilsynet:provider-jetstream-nats/jetstream-publish;
  export mattilsynet:provider-jetstream-nats/jetstream-consumer;
  
  //nats-kv watch keys hasn't been implemented yet in github.com/Mattilsynet/map-nats-kv
  import mattilsynet:map-kv/key-value;
  
}
`
	tmpls := project.Templs
	config := NewConfig(
		"/home/solve/git/temp",
		"test-component",
		"github.com/Mattilsynet/map-test",
		[]string{"nats-core", "nats-jetstream", "nats-kv"},
		WithComponentCode(),
		WithWitPackage())
	config.ImportNatsCoreWit = true
	config.ExportNatsCoreWit = true
	config.ImportNatsJetstreamWit = true
	config.ExportNatsJetstreamWit = true
	config.ImportNatsKvWit = true
	a, err := ExecuteTmplWithData(config, tmpls["wit/world.wit"])
	if err != nil {
		log.Println("err: ", err)
		t.Fail()
	}
	if strings.TrimSpace(a) != strings.TrimSpace(expected) {
		t.Fail()
	}
}

func TestGeneratePkgNats(t *testing.T) {
	expected := `
	package nats

import (
	"errors"
	"log/slog"

	"github.com/Mattilsynet/map-test/gen/wasmcloud/messaging/consumer"
	"github.com/Mattilsynet/map-test/gen/wasmcloud/messaging/handler"
	"github.com/Mattilsynet/map-test/gen/wasmcloud/messaging/types"
	"github.com/bytecodealliance/wasm-tools-go/cm"
)
type (
	Conn struct {
        js JetStreamContext
	}
	JetStreamContext struct {
		bucket KeyValue
		}

	Msg struct {
		Subject string
		Reply   string
		Data    []byte
		Header  map[string][]string
	}
)

func NewConn() *Conn {
	return &Conn{}
}
type MsgHandler func(msg *Msg)


func toWitExportSubscription(msgHandler MsgHandler) func(msg jetstreamconsumer.Msg) cm.Result[string, struct{}, string] {
	return func(msg jetstreamconsumer.Msg) cm.Result[string, struct{}, string] {
		natsMsg := &Msg{
			Subject: msg.Subject,
			Reply:   msg.Reply,
			Data:    msg.Data.Slice(),
			Header:  toNatsHeaders(msg.Headers),
		}
		msgHandler(natsMsg)
		return cm.OK[cm.Result[string, struct{}, string]](struct{}{})
	}
}
func toNatsHeaders(header cm.List[jetstream_types.KeyValue]) map[string][]string {
	natsHeaders := make(map[string][]string)
	for _, kv := range header.Slice() {
		natsHeaders[kv.Key] = kv.Value.Slice()
	}
	return natsHeaders
}
func FromBrokerMessageToNatsMessage(bm types.BrokerMessage) *Msg {
	if bm.ReplyTo.None() {
		return &Msg{
			Data:    bm.Body.Slice(),
			Subject: bm.Subject,
			Reply:   "",
		}
	} else {
		return &Msg{
			Data:    bm.Body.Slice(),
			Subject: bm.Subject,
			Reply:   *bm.ReplyTo.Some(),
		}
	}
}


func ToBrokenMessageFromNatsMessage(nm *Msg) types.BrokerMessage {
	if nm.Reply == "" {
		return types.BrokerMessage{
			Subject: nm.Subject,
			Body:    cm.ToList(nm.Data),
			ReplyTo: cm.None[string](),
		}
	} else {
		return types.BrokerMessage{
			Subject: nm.Subject,
			Body:    cm.ToList(nm.Data),
			ReplyTo: cm.Some(nm.Subject),
		}
	}
}

func (nc *Conn) Publish(msg *Msg) error {
	bm := ToBrokenMessageFromNatsMessage(msg)
	result := consumer.Publish(bm)
	if result.IsErr() {
		return errors.New(*result.Err())
	}
	return nil
}

func (conn *Conn) RegisterSubscription(fn func(*Msg)) {
	handler.Exports.HandleMessage = func(msg types.BrokerMessage) (result cm.Result[string, struct{}, string]) {
		natsMsg := FromBrokerMessageToNatsMessage(msg)
		fn(natsMsg)
		return cm.OK[cm.Result[string, struct{}, string]](struct{}{})
	}
}

	`
	tmpls := project.Templs
	config := NewConfig(
		"/home/solve/git/temp",
		"test-component",
		"github.com/Mattilsynet/map-test",
		[]string{"nats-core", "nats-jetstream", "nats-kv"},
		WithComponentCode(),
		WithWitPackage())
	config.ImportNatsCoreWit = true
	config.ExportNatsCoreWit = true
	config.ImportNatsJetstreamWit = true
	config.ExportNatsJetstreamWit = true
	config.ImportNatsKvWit = true
	a, err := ExecuteTmplWithData(config, tmpls["pkg/nats/nats.go"])
	if err != nil {
		log.Println("err: ", err)
		t.Fail()
	}
	if strings.TrimSpace(a) != strings.TrimSpace(expected) {
		t.Fail()
	}
}

func TestGeneratePkgJetstream(t *testing.T) {
	expected := `package nats

func (c *Conn) Jetstream() (*JetStreamContext, error) {
	return &c.js, nil
}
func (js *JetStreamContext) Subscribe(msgHandler MsgHandler, logger *slog.Logger) {
	jetstreamconsumer.Exports.HandleMessage = toWitExportSubscription(msgHandler)
}
func (js *JetStreamContext) Publish(subj string, data []byte) error {
	return js.PublishMsg(&Msg{Subject: subj, Data: data})
}

func (js *JetStreamContext) PublishMsg(msg *Msg) error {
	jpMsg := jetstreampublish.Msg{
		Headers: toWitNatsHeaders(msg.Header),
		Data:    cm.ToList(msg.Data),
		Subject: msg.Subject,
	}
	result := jetstreampublish.Publish(jpMsg)
	if !result.IsOK() {
		return errors.New(*result.Err())
	}
	return nil
}

func toWitNatsHeaders(header map[string][]string) cm.List[jetstream_types.KeyValue] {
	keyValueList := make([]jetstream_types.KeyValue, 0)
	for k, v := range header {
		keyValueList = append(keyValueList, jetstream_types.KeyValue{
			Key:   k,
			Value: cm.ToList(v),
		})
	}
	return cm.ToList(keyValueList)
}`
	tmpls := project.Templs
	config := NewConfig(
		"/home/solve/git/temp",
		"test-component",
		"github.com/Mattilsynet/map-test",
		[]string{"nats-core", "nats-jetstream", "nats-kv"},
		WithComponentCode(),
		WithWitPackage())
	config.ImportNatsCoreWit = true
	config.ExportNatsCoreWit = true
	config.ImportNatsJetstreamWit = true
	config.ExportNatsJetstreamWit = true
	config.ImportNatsKvWit = true
	a, err := ExecuteTmplWithData(config, tmpls["pkg/nats/js.go"])
	if err != nil {
		log.Println("err: ", err)
		t.Fail()
	}
	if strings.TrimSpace(a) != strings.TrimSpace(expected) {
		t.Fail()
	}
}

func TestGeneratePkgKv(t *testing.T) {
	expected := `package nats

import (
	"errors"

	"github.com/Mattilsynet/map-test/gen/mattilsynet/map-kv/key-value"
	"github.com/bytecodealliance/wasm-tools-go/cm"
)

type (
	KeyValue      struct{}
	KeyValueEntry struct {
		key   string
		value []byte
	}
)

func (e *KeyValueEntry) Key() string   { return e.key }
func (e *KeyValueEntry) Value() []byte { return e.value }

func (js *JetStreamContext) KeyValue() (*KeyValue, error) {
	js.bucket = KeyValue{}
	return &js.bucket, nil
}

func (js *KeyValue) Get(key string) (*KeyValueEntry, error) {
	result := keyvalue.Get(key)
	if result.IsOK() {
		resVal := result.OK().Value.Slice()
		resKey := result.OK().Key
		return &KeyValueEntry{resKey, resVal}, nil
	}
	if result.IsErr() {
		return nil, errors.New(*result.Err())
	}
	return nil, errors.New("unknown error when getting keyvalue from map-kv with key: " + key)
}

func (js *KeyValue) GetAll() ([]*KeyValueEntry, error) {
	listKeys := keyvalue.ListKeys()
	if listKeys.IsOK() {
		keys := listKeys.OK().Slice()
		var entries []*KeyValueEntry
		for _, key := range keys {
			result := keyvalue.Get(key)
			if result.IsOK() {
				resVal := result.OK().Value.Slice()
				resKey := result.OK().Key
				entries = append(entries, &KeyValueEntry{resKey, resVal})
			}
			if result.IsErr() {
				return nil, errors.New(*result.Err())
			}
		}
		return entries, nil
	}
	if listKeys.IsErr() {
		return nil, errors.New(*listKeys.Err())
	}
	return nil, errors.New("unknown error when getting all keyvalues from map-kv")
}

func (js *KeyValue) Put(key string, value []byte) error {
	result := keyvalue.Put(key, cm.ToList(value))
	if result.IsOK() {
		return nil
	}
	if result.IsErr() {
		return errors.New(*result.Err())
	}
	return errors.New("unknown error when putting keyvalue in map-kv with key: " + key)
}

func (js *KeyValue) Create(key string, value []byte) error {
	result := keyvalue.Create(key, cm.ToList(value))
	if result.IsOK() {
		return nil
	}
	if result.IsErr() {
		return errors.New(*result.Err())
	}
	return errors.New("unknown error when creating keyvalue in map-kv with key: " + key)
}

func (js *KeyValue) Delete(key string) error {
	result := keyvalue.Delete(key)
	if result.IsOK() {
		return nil
	}
	if result.IsErr() {
		return errors.New(*result.Err())
	}
	return errors.New("unknown error when deleting keyvalue in map-kv with key: " + key)
}

func (js *KeyValue) ListKeys() ([]string, error) {
	result := keyvalue.ListKeys()
	if result.IsOK() {
		return result.OK().Slice(), nil
	}
	if result.IsErr() {
		return nil, errors.New(*result.Err())
	}
	return nil, errors.New("unknown error when listing keys in map-kv")
}
`
	tmpls := project.Templs
	config := NewConfig(
		"/home/solve/git/temp",
		"test-component",
		"github.com/Mattilsynet/map-test",
		[]string{"nats-core", "nats-jetstream", "nats-kv"},
		WithComponentCode(),
		WithWitPackage())
	config.ImportNatsCoreWit = true
	config.ExportNatsCoreWit = true
	config.ImportNatsJetstreamWit = true
	config.ExportNatsJetstreamWit = true
	config.ImportNatsKvWit = true
	a, err := ExecuteTmplWithData(config, tmpls["pkg/nats/kv.go"])
	if err != nil {
		log.Println("err: ", err)
		t.Fail()
	}
	if strings.TrimSpace(a) != strings.TrimSpace(expected) {
		t.Fail()
	}
}

func TestGenerateLocalWadm(t *testing.T) {
	expected := `apiVersion: core.oam.dev/v1beta1
kind: Application
metadata:
  name: test-component
  annotations:
    version: v0.0.1
    description: "description"
    authors: "authors"
spec:
  components:
    - name: test-component
      type: component
      properties:
        image: file://./build/test-component_s.wasm
      traits:
        - type: spreadscaler
          properties:
            replicas: 1
- type: link
          properties:
            target:
              name: nats-core
            namespace: wasmcloud
            package: messaging
            interfaces: [consumer]
- type: link
          properties:
            target:
              name: map-nats-kv
              config:
                - name: map-nats-kv-config
                  properties:
                    bucket: "my-bucket"
                    url: "nats://localhost:4222"
            namespace: mattilsynet
            package: map-kv
            interfaces: [key-value]
- name: nats-core
      type: capability
      properties:
        image: ghcr.io/wasmcloud/messaging-nats:canary
        config:
          - name: nats-core-config
            properties:
              cluster_uris: "nats://localhost:4222"
      traits:
        - type: spreadscalar
          properties:
            replicas: 1
- name: nats-jetstream
      type: capability
      properties:
        image: ghcr.io/Mattilsynet/map-nats-jetstream:v0.0.1-pre-17
      traits:
        - type: link
          properties:
            target:
              name: test-component
            source:
              config:
                - name: nats-jetstream-nats-url
                  properties:
                    url: "nats://localhost:4222"
                - name: nats-jetstream-consumer-config
                  properties:
                    stream-name: "stream-name"
                    stream-retention-policy: "workqueue" # oneof "interest, workqueue, limits"
                    subject: "special.subject.>"
                    durable-consumer-name: "test-component-consumer"
            namespace: mattilsynet
            package: provider-jetstream-nats
            interfaces: [jetstream-consumer]
- name: map-nats-kv
      type: capability
      properties:
        image: ghcr.io/mattilsynet/map-nats-kv:latest`
	_ = expected
	tmpls := project.Templs
	config := NewConfig(
		"/home/solve/git/temp",
		"test-component",
		"github.com/Mattilsynet/map-test",
		[]string{"nats-core", "nats-jetstream", "nats-kv"},
		WithComponentCode(),
		WithWitPackage())
	config.ImportNatsCoreWit = true
	config.ExportNatsCoreWit = true
	config.ImportNatsJetstreamWit = true
	config.ExportNatsJetstreamWit = true
	config.ImportNatsKvWit = true
	a, err := ExecuteTmplWithData(config, tmpls["local.wadm.yaml"])
	if err != nil {
		log.Println("err: ", err)
		t.Fail()
	}
	if strings.TrimSpace(a) != strings.TrimSpace(expected) {
		t.Fail()
	}
}
