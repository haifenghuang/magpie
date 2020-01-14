/*
    Rest Service Example
    For more advanced example, please see 'service.mp'
*/

let login = fn(writer, request) {
    if request.method == http.METHOD_POST {
        writer.header.add("Content-Type", "application/json; charset=utf-8")
        writer.write(json.toJson({ "sessionId": "3d5bd2cA15ef0476" }))
    }
}

let logout = fn(writer, request) {
    if request.method == http.METHOD_POST {
        writer.writeHeader(http.STATUS_CREATED) // return http status code 201
    }
}

http.handleFunc("/authentication/login", login)
http.handleFunc("/authentication/logout", logout)

addr = "0.0.0.0:8080"
printf("* Running on %s\n", addr)
http.listenAndServe(addr)
