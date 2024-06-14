// Copyright (c) 2023 suwei007@gmail.com
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

// StreamingOutputHandler is a function type that handles streaming output byte slices.
type StreamingOutputHandler func(ctx context.Context, part []byte) error

// ConverseStreamingOutputDeltaHandler is a function type that handles streaming output delta updates.
type ConverseStreamingOutputDeltaHandler func(ctx context.Context, part types.ContentBlockDelta) error

// ConverseStreamingOutputStartHandler is a function type that handles streaming output start events.
type ConverseStreamingOutputStartHandler func(ctx context.Context, part types.ContentBlockStart) error

// Response represents the response structure for a completion or streaming output.
type Response struct {
	Completion string `json:"completion"`
	StopReason string `json:"stopReason"`
}
