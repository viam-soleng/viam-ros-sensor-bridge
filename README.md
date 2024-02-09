# viam-ros-sensor-bridge

This is a module that is not intended to be published to the public Viam Registry. Instead it is intended to be forked and implemented in private for each user. This is because custome ROS types have to be mapped and this could expose proprietary information.

## How to use this
1. Use the tools from goroslib to convert the IDL files to go structs
2. Add the struct to [custom_messages.go](messages/custom_messages.go) (or your own file in that package) 
3. Update the [custom_registry](messages/custom_messages.go#L8) to add the new type to the registry so that your new types are properly handled. Please see [ThrottlingStates](messages/msgs.go#L13) for an example on doing this. *DO NOT CHANGE `standard_messages.go` or `message_handler.go`*
4. Add any appropriate [tests](messages/custom_messages_test.go)
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

## FAQ
### Q: What format should sensor values be in?
A sensor needs to return data in the exact format of the message to be sent to ROS. You can find the format of ROS message data by looking at the definitions in [`goroslib`](https://github.com/bluenviron/goroslib/tree/main/pkg/msgs/std_msgs).

For example, if you want to send an [Int32](https://github.com/bluenviron/goroslib/blob/main/pkg/msgs/std_msgs/msg_int32.go) from your sensor to ROS, the readings() method needs to return the integer value where the key is `Data`.
```
{ Data: 42 }
```

For a [string](https://github.com/bluenviron/goroslib/blob/main/pkg/msgs/std_msgs/msg_string.go), the format is:
```
{ Data: "Hello World!" }
```

### Q: I don't want to make my message types public, how do I keep in sync with your repository?
It is not uncommon that you may want to protect your messages by not publishing them publicly. The easiest way to do this is to create a private repository on GitHub and manually syncing this repository with yours.
1. Create an empty private repository under your organization
1. Clone this repository
1. Change the remote name
   ```
   git remote rename origin upstream
   ```
1. Add your new repository as the origin
   ```
   git remote add origin <new_git_uri>
   ```
1. Push the code to your repository
   ```
   git push origin
   ```

Now to sync changes from upstream:
```
git fetch upstream
git push origin
```

Be careful when making changes to limit them just to just the following files: 
* [messages/custom_messages.go](messages/custom_messages.go)
* [messages/custom_messages_test.go](messages/custom_messages_test.go)
* [utils/utils.go](utils/utils.go)

I will try to restrict changes to these files so they won't conflict with your local changes.
