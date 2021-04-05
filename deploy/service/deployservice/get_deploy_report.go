package deployservice

import (
	"context"
	"io/ioutil"
	"strconv"

	"github.com/aws/aws-sdk-go/service/s3"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"

	deploypb "github.com/mjm/pi-tools/deploy/proto/deploy"
)

func (s *Server) GetDeployReport(ctx context.Context, req *deploypb.GetDeployReportRequest) (*deploypb.GetDeployReportResponse, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(label.Int64("deployment.id", req.GetDeployId()))

	key := strconv.FormatInt(req.GetDeployId(), 10)
	span.SetAttributes(
		label.String("report.bucket", s.Config.ReportBucket),
		label.String("report.key", key))

	res, err := s.S3.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: &s.Config.ReportBucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	msgData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(label.Int("report.size", len(msgData)))

	var report deploypb.Report
	if err := proto.Unmarshal(msgData, &report); err != nil {
		return nil, err
	}

	span.SetAttributes(
		label.String("report.commit_sha", report.GetCommitSha()),
		label.String("report.commit_message", report.GetCommitMessage()))

	return &deploypb.GetDeployReportResponse{
		Report: &report,
	}, nil
}
