// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package servicecontrol

import (
	"time"

	pb "github.com/googleapis/googleapis/google/api/servicecontrol/v1/"
	"google.golang.org/grpc"
)

type clientState struct {
	client     pb.ServiceControlClient
	connection *grpc.ClientConn
}

func createAPIClient(address string) (*clientState, error) {
	cs := clientState{}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	var err error
	if cs.connection, err = grpc.Dial(address, opts...); err != nil {
		return nil, err
	}

	cs.client = pb.NewServiceControlClient(cs.connection)
	return &cs, nil
}

func deleteAPIClient(cs *clientState) error {
	// TODO: This is to compensate for this bug: https://github.com/grpc/grpc-go/issues/1059
	//       Remove this delay once that bug is fixed.
	time.Sleep(50 * time.Millisecond)

	err := cs.connection.Close()
	cs.client = nil
	cs.connection = nil
	return err
}
