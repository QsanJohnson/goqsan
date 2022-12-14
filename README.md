
# goqsan
Go http client to manage Qsan XEVO models.

## Install
```
go get github.com/QsanJohnson/goqsan
```

## Usage
Here is an sample code.
```
import (
	"github.com/QsanJohnson/goqsan"
	"fmt"
	"context"
)
	
ctx := context.Background()

client := goqsan.NewClient("192.xxx.xxx.xxx")
systemAPI := goqsan.NewSystem(client)
res, err := systemAPI.GetAbout(ctx)
if err == nil {
	fmt.Printf("%+v\n", res);
}

authClient, err := client.GetAuthClient(ctx, "admin", "1234")
volumeAPI := goqsan.NewVolume(authClient)
vols, err := volumeAPI.ListVolumes(ctx, "")
if err == nil {
	fmt.Printf("%+v\n", vols);
}
```

## Debugging
Add flag.Parse() at the begining in main(),
then execute go run with "-v=4 -alsologtostderr" arguments.
```
go run test.go -v=4 -alsologtostderr
```


## Testing

You have to create a test.conf file for integration test. The following is an example,
```
QSAN_IP = 192.xxx.xxx.xxx
QSAN_USERNAME = admin
QSAN_PASSWORD = 1234
POOL_ID = xxxxxx
```
* POOL_ID is Pool ID to be created/deleted volume on.

Then run integration test
```
go test
```

Or run integration test with log level
```
export GOQSAN_LOG_LEVEL=4
go test
```