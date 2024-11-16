# Websocket Router

**Package Not Production Ready**

## Step 1 Create Top Level Router

``
tlRtr, err := CreateToplevelRouter()
``

## Sep 2 Create subrouter
``
subRouter, err := router.CreateSubRouter("path2")
``

This creates a router at the pat ```/{router path}/path2```.

## Step 3 Create endpoint
```
router.CreateEndpoint("user", func(payload []byte, companyID int) ([]byte, error) {

fmt.Println("endpoint user")

time.Sleep(2 * time.Second)

return payload, nil
})
```

## Step 4 Accept Request
You can start handling request at any router level but only the top level router has access to all sub routers.
```
response, err := tlRtr.HandleRequest("one/two/three", {Some json payload})
```