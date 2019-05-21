//
// Copyright 2019 Insolar Technologies GmbH
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
//

package bus

// func RetryIncorrectPulse(ctx context.Context, sender Sender, msg *message.Message) (<-chan *message.Message, func()) {
// 	const retryCount = 3
//
// 	retries := retryCount
// 	for {
// 		reps, done := sender.Send(ctx, msg)
// 		rep, ok := <-reps
// 		if !ok {
// 			return reps, done
// 		}
//
// 		if rep.Metadata.Get(MetaType) != payload.TypeError {
// 			res := make(chan *message.Message, 1)
// 			res <- rep
// 			return res, func() {
// 				done()
// 				close(res)
// 			}
// 		}
//
// 		pl := payload.Error{}
// 		err := pl.Unmarshal(rep.Payload)
// 		if err != nil {
// 			inslogger.FromContext(ctx).Error("Failed to decode reply")
// 		}
// 		if err == nil || !strings.Contains(err.Error(), "Incorrect message pulse") {
// 			return rep, err
// 		}
//
// 		if retries <= 0 {
// 			inslogger.FromContext(ctx).Warn("got incorrect message pulse too many times")
// 			return rep, err
//
// 		}
// 		retries--
//
// 		time.Sleep(100 * time.Millisecond)
// 	}
// }
