package messenger

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"gitlab.x.lan/yunshan/droplet-libs/app"
	dt "gitlab.x.lan/yunshan/droplet-libs/zerodoc"
	pb "gitlab.x.lan/yunshan/message/zero"
)

func Marshal(doc *app.Document) ([]byte, error) {
	if doc.Tag == nil || doc.Meter == nil {
		return nil, errors.New("No tag or meter in document")
	}

	var msgType MessageType
	switch v := doc.Meter.(type) {
	case *dt.UsageMeter:
		msgType = MSG_USAGE
	case *dt.PerfMeter:
		msgType = MSG_PERF
	case *dt.GeoMeter:
		msgType = MSG_GEO
	case *dt.FlowMeter:
		msgType = MSG_FLOW
	case *dt.PlatformMeter:
		msgType = MSG_PLATFORM
	case *dt.ConsoleLogMeter:
		msgType = MSG_CONSOLE_LOG
	case *dt.TypeMeter:
		msgType = MSG_TYPE
	default:
		return nil, fmt.Errorf("Unknown supported type %T", v)
	}

	msg := &pb.ZeroDocument{}
	msg.Timestamp = proto.Uint32(doc.Timestamp)

	var tag *dt.Tag
	tag, ok := doc.Tag.(*dt.Tag)
	if !ok {
		return nil, fmt.Errorf("Unknown supported tag type %T", doc.Tag)
	}
	msg.Tag = dt.TagToPB(tag)

	msg.Meter = &pb.Meter{}
	switch msgType {
	case MSG_USAGE:
		meter := doc.Meter.(*dt.UsageMeter)
		msg.Meter.Usage = dt.UsageMeterToPB(meter)
	case MSG_PERF:
		meter := doc.Meter.(*dt.PerfMeter)
		msg.Meter.Perf = dt.PerfMeterToPB(meter)
	case MSG_GEO:
		meter := doc.Meter.(*dt.GeoMeter)
		msg.Meter.Geo = dt.GeoMeterToPB(meter)
	case MSG_FLOW:
		meter := doc.Meter.(*dt.FlowMeter)
		msg.Meter.Flow = dt.FlowMeterToPB(meter)
	case MSG_PLATFORM:
		meter := doc.Meter.(*dt.PlatformMeter)
		msg.Meter.Platform = dt.PlatformMeterToPB(meter)
	case MSG_CONSOLE_LOG:
		meter := doc.Meter.(*dt.ConsoleLogMeter)
		msg.Meter.ConsoleLog = dt.ConsoleLogMeterToPB(meter)
	case MSG_TYPE:
		meter := doc.Meter.(*dt.TypeMeter)
		msg.Meter.Type = dt.TypeMeterToPB(meter)
	}
	msg.Actions = proto.Uint32(doc.Actions)

	// TODO: 传入buffer
	b := make([]byte, msg.Size())
	_, err := msg.MarshalTo(b)
	if err != nil {
		return nil, fmt.Errorf("Marshaling protobuf failed: %s", err)
	}

	return b, nil
}

func Unmarshal(b []byte) (*app.Document, error) {
	if b == nil {
		return nil, errors.New("No input byte")
	}

	msg := &pb.ZeroDocument{}
	if err := msg.Unmarshal(b); err != nil {
		return nil, fmt.Errorf("Unmarshaling protobuf failed: %s", err)
	}

	doc := &app.Document{}
	doc.Timestamp = msg.GetTimestamp()
	doc.Tag = dt.PBToTag(msg.GetTag())
	meter := msg.GetMeter()
	switch {
	case meter.GetUsage() != nil:
		doc.Meter = dt.PBToUsageMeter(meter.GetUsage())
	case meter.GetPerf() != nil:
		doc.Meter = dt.PBToPerfMeter(meter.GetPerf())
	case meter.GetGeo() != nil:
		doc.Meter = dt.PBToGeoMeter(meter.GetGeo())
	case meter.GetFlow() != nil:
		doc.Meter = dt.PBToFlowMeter(meter.GetFlow())
	case meter.GetPlatform() != nil:
		doc.Meter = dt.PBToPlatformMeter(meter.GetPlatform())
	case meter.GetConsoleLog() != nil:
		doc.Meter = dt.PBToConsoleLogMeter(meter.GetConsoleLog())
	case meter.GetType() != nil:
		doc.Meter = dt.PBToTypeMeter(meter.GetType())
	}
	doc.Actions = msg.GetActions()

	return doc, nil
}
