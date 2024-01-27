# viam-ros-sensor-bridge

This is a module that is not intended to be published to the public Viam Registry. Instead it is intended to be forked and implemented in private for each user. This is because custome ROS types have to be mapped and this could expose proprietary information.

## How to use this
1. Use the tools from goroslib to convert the IDL files to go structs
2. Add the struct to [msgs.go](messages/msgs.go) (or your own file in that package) 
3. Update the [GetMessageType](messages/msgs.go#L49) and [ConvertMessage](messages/msgs.go#L64) methods so that your new types are properly handled. Please see [ThrottlingStates](messages/msgs.go#L119) and [convertThrottlingStates](messages/msgs.go#L132) for demos on doing this.
4. Add any appropriate [tests](messages/msgs_test.go)
5. Compile the module
   ```
   go build -o viam-ros-sensor-bridge module.go
   ```
   Be sure to use `GOOS` and `GOARCH` environment variables when necessary
6. (Optional) Move the binary to your robot and test
7. (Optional) Publish to your private Viam Repository
   
   **DO NOT PUBLISH TO THE PUBLIC REPOSITORY**
