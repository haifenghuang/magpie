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

addr = "127.0.0.1:8080"
printf("* Running on %s\n", addr)
http.listenAndServe(addr)
