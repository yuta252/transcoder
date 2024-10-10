package main

import (
	"context"
	"fmt"
	"os"

	transcoder "cloud.google.com/go/video/transcoder/apiv1"
	"cloud.google.com/go/video/transcoder/apiv1/transcoderpb"
	"google.golang.org/protobuf/types/known/durationpb"
)

type reqType string

const (
	reqTypeGet    reqType = "GET"
	reqTypeCreate reqType = "CREATE"
)

func getJob(projectID string, location string, jobID string) (*transcoderpb.Job, error) {
	ctx := context.Background()
	client, err := transcoder.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &transcoderpb.GetJobRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectID, location, jobID),
	}

	response, err := client.GetJob(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetJob: %w", err)
	}

	return response, nil
}

// createJobWithSetNumberImagesSpritesheet creates a job from an ad-hoc configuration and generates
// two spritesheets from the input video. Each spritesheet contains a set number of images.
func createJobWithSetNumberImagesSpritesheet(projectID string, location string, inputURI string, outputURI string) (*transcoderpb.Job, error) {
	ctx := context.Background()
	client, err := transcoder.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	req := &transcoderpb.CreateJobRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectID, location),
		Job: &transcoderpb.Job{
			InputUri:  inputURI,
			OutputUri: outputURI,
			JobConfig: &transcoderpb.Job_Config{
				Config: &transcoderpb.JobConfig{
					ElementaryStreams: []*transcoderpb.ElementaryStream{
						{
							Key: "video_stream0",
							ElementaryStream: &transcoderpb.ElementaryStream_VideoStream{
								VideoStream: &transcoderpb.VideoStream{
									CodecSettings: &transcoderpb.VideoStream_H264{
										H264: &transcoderpb.VideoStream_H264CodecSettings{
											BitrateBps:   550000,
											FrameRate:    60,
											HeightPixels: 360,
											WidthPixels:  640,
										},
									},
								},
							},
						},
					},
					MuxStreams: []*transcoderpb.MuxStream{
						{
							Key:               "output_video0",
							Container:         "mp4",
							ElementaryStreams: []string{"video_stream0"},
						},
					},
					SpriteSheets: []*transcoderpb.SpriteSheet{
						{
							FilePrefix:         "large-sprite-sheet",
							SpriteWidthPixels:  128,
							SpriteHeightPixels: 72,
							ColumnCount:        1,
							RowCount:           1,
							// ExtractionStrategy: &transcoderpb.SpriteSheet_Interval{Interval: durationpb.New(1)},
							StartTimeOffset: durationpb.New(0),
							Quality:         100,
						},
					},
				},
			},
		},
	}
	// Creates the job. Jobs take a variable amount of time to run. You can query for the job state.
	// See https://cloud.google.com/transcoder/docs/how-to/jobs#check_job_status for more info.
	response, err := client.CreateJob(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("createJobWithSetNumberImagesSpritesheet: %w", err)
	}

	return response, nil
}

func main() {
	// Set your projectID
	projectID := "ProjectID"
	location := "asia-northeast1"
	inputURI := "gs://test-nakano-thumbnail-extraction/input/sample.mp4"
	outputURI := "gs://test-nakano-thumbnail-extraction/output/"
	reqType := reqTypeCreate

	if reqType == reqTypeCreate {
		response, err := createJobWithSetNumberImagesSpritesheet(projectID, location, inputURI, outputURI)
		if err != nil {
			fmt.Println("failed to create job:", err)
			os.Exit(1)
		}
		fmt.Printf("get job: %+v", response)
	} else if reqType == reqTypeGet {
		response, err := getJob(projectID, location, "a942a66f-8b9a-4e00-8d05-ad6479f3e63b")
		if err != nil {
			fmt.Println("failed to get job:", err)
			os.Exit(1)
		}
		fmt.Printf("get job: %+v", response)
	}
}
