# viam-ros-sensor-bridge

This is a module that is not intended to be published to the public Viam Registry. Instead it is intended to be forked and implemented in private for each user. This is because custome ROS types have to be mapped and this could expose proprietary information.

## How to use this
1. Use the tools from goroslib to convert the IDL files to go structs
2. Add the struct to [msgs.go](messages/msgs.go) (or your own file in that package) 
3. Update the [custom_registry](messages/msgs.go#L40) to add the new type to the registry so that your new types are properly handled. Please see [ThrottlingStates](messages/msgs.go#L87) for an example on doing this. *DO NOT CHANGE `std_msgs_registry`*
4. Add any appropriate [tests](messages/msgs_test.go)
5. Compile the module
   ```
   go build -o viam-ros-sensor-bridge module.go
   ```
   Be sure to use `GOOS` and `GOARCH` environment variables when necessary. For example, to build on x86_64 to run on a Raspberry Pi.
   ```
   GOOS=linux GOARCH=arm64 go build -o viam-ros-sensor-bridge module.go
   ```
6. (Optional) Move the binary to your robot and test
7. (Optional) Publish to your private Viam Repository
   
   ---
   **DO NOT PUBLISH TO THE PUBLIC REPOSITORY**
